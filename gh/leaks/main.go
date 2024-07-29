package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
)

type Event struct {
	Type string
	Repo struct {
		Name string // user/repo
	}
}

type Repo struct {
	Name  string
	Path  string // local clone
	Leaks []Leak
	Err   error
}

type Leak struct {
	Description string
	Secret      string
	File        string
	Entropy     float64
	Email       string
	Date        time.Time
	Message     string
	RuleID      string
}

type Count struct {
	mu            sync.Mutex
	publicEvents  int
	reposCloned   int
	reposSearched int
}

const progname = "ghleaks"

func main() {
	log.SetFlags(0)
	log.SetPrefix(progname + ": ")

	c := flag.Int("c", 100, "repos to clone in parallel")
	r := flag.String("r", "/tmp", "where to clone the repos")
	s := flag.Int("s", 100, "repos to search for leaks in parallel")
	// v := flag.Bool("v", false, "print details about found leaks")
	flag.Parse()

	madePublic := make(chan Repo)
	cloned := make(chan Repo)
	searched := make(chan Repo)

	var count Count

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		dec := json.NewDecoder(os.Stdin)
		for {
			var event Event
			if err := dec.Decode(&event); err == io.EOF {
				break
			} else if err != nil {
				log.Printf("decode json: %v", err)
			}
			if event.Type == "PublicEvent" {
				madePublic <- Repo{Name: event.Repo.Name}
				count.mu.Lock()
				count.publicEvents++
				count.mu.Unlock()
			}
		}
		close(madePublic)
		wg.Done()
	}()

	// https://github.com/orgs/community/discussions/44515
	for i := 0; i < *c; i++ {
		wg.Add(1)
		go func() {
			for repo := range madePublic {
				if repo.Err == nil {
					repo.Path, repo.Err = cloneRepo(repo, *r)
				}
				cloned <- repo
				count.mu.Lock()
				count.reposCloned++
				count.mu.Unlock()
			}
			close(cloned)
			wg.Done()
		}()
	}

	for i := 0; i < *s; i++ {
		wg.Add(1)
		go func() {
			for repo := range cloned {
				if repo.Err == nil {
					repo.Leaks, repo.Err = searchLeaks(repo, *r)
				}
				searched <- repo
				count.mu.Lock()
				count.reposSearched++
				count.mu.Unlock()
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(searched)
	}()

	for repo := range searched {
		if repo.Err != nil {
			log.Print(repo.Err)
		} else if len(repo.Leaks) > 0 {
			data, err := json.MarshalIndent(repo, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(data))
		}
		log.Printf("events: %d, cloned: %d, searched: %d\n", count.publicEvents, count.reposCloned, count.reposSearched)
	}
}

// cloneRepo clones the repo into reposDir and returns the full path of the
// cloned the repo.  If the repo already exists it does nothing.
func cloneRepo(repo Repo, reposDir string) (string, error) {
	dir := filepath.Join(reposDir, repo.Name)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://github.com/%s", repo.Name)
	_, err := git.PlainClone(dir, false, &git.CloneOptions{URL: url})
	if err != nil && !errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return dir, fmt.Errorf("clone %s: %v", url, err)
	}

	return dir, nil
}

// searchLeaks searches for leaks in the given repository. If leaks are found it
// returns the number and details about the leaks.
func searchLeaks(repo Repo, reposDir string) ([]Leak, error) {
	reportPath := filepath.Join(reposDir, repo.Name, "ghleaks.json")
	out, err := exec.Command("gitleaks", "detect", "--no-color", "--no-banner", "--verbose", "--source", repo.Path, "--report-path", reportPath).CombinedOutput()
	if err != nil { // gitleaks exits with 1 leaks are found or when there's an error

		data, err := os.ReadFile(reportPath)
		if err != nil {
			return nil, err
		}
		var leaks []Leak
		if err := json.Unmarshal(data, &leaks); err != nil {
			return nil, err
		}

		if len(leaks) > 0 {
			return leaks, nil
		}

		// leaksfound := regexp.MustCompile(`WRN leaks found: (\d+)`) // 1:32PM WRN leaks found: 1
		// for _, line := range strings.Split(string(out), "\n") {
		// 	match := leaksfound.FindStringSubmatch(line)
		// 	if len(match) > 0 {
		// 		n, err := strconv.Atoi(match[1])
		// 		if err != nil {
		// 			return 0, "", err
		// 		}
		// 		return n, string(out), nil
		// 	}
		// }

		return nil, fmt.Errorf("run gitleaks: %s: %v", out, err)
	}

	return nil, nil
}
