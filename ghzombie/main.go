package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Event struct {
	Type string
	Repo struct {
		Name string // user/repo
	}
	Payload struct {
		Forkee struct {
			FullName string `json:"full_name"` // user/repo
			HtmlUrl  string `json:"html_url"`
		}
		Commits []Commit
	}
}

type Commit struct {
	Sha     string
	Message string
}

type PushEvent struct {
	Repo    string
	Commits []Commit
}

type ForkEvent struct {
	OrigRepo string
	ForkRepo string
	ForkUrl  string
}

type output struct {
	Deleted string
	Commits []outputCommit
}

type outputCommit struct {
	Message string
	Url     string
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("ghzombie: ")

	if len(os.Args) < 2 {
		log.Fatal("supply json file(s) with GitHub events (see gharchive.org)")
	}
	files := os.Args[1:]

	log.Printf("going to parse %d files", len(files))

	events := make(chan Event)
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func() {
			defer log.Printf("finished parsing %s for fork and push events\n", file)
			defer wg.Done()
			es, err := getForkandPushEvents(file)
			if err != nil {
				log.Print(err)
				return
			}
			for _, e := range es {
				events <- e
			}
		}()
	}

	go func() {
		wg.Wait()
		close(events)
	}()

	var pushEvents []PushEvent
	var forkEvents []ForkEvent

	for e := range events {
		switch e.Type {
		case "PushEvent":
			pushEvents = append(pushEvents, PushEvent{
				Repo:    e.Repo.Name,
				Commits: e.Payload.Commits,
			})
		case "ForkEvent":
			forkEvents = append(forkEvents, ForkEvent{
				OrigRepo: e.Repo.Name,
				ForkRepo: e.Payload.Forkee.FullName,
				ForkUrl:  e.Payload.Forkee.HtmlUrl,
			})
		}
	}

	log.Printf("found %d fork events (and %d push events)", len(forkEvents), len(pushEvents))

	if len(forkEvents) <= 0 {
		os.Exit(0)
	}

	all := make(chan ForkEvent)
	deleted := make(chan ForkEvent)

	var mu sync.Mutex
	var count int

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fe := range all {
				mu.Lock()
				count++
				log.Printf("checking fork %s (%d/%d)", fe.ForkUrl, count, len(forkEvents))
				mu.Unlock()
				if isNotFound(fe.ForkUrl) {
					deleted <- fe
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(deleted)
	}()

	go func() {
		for _, fe := range forkEvents {
			all <- fe
		}
	}()

	for fe := range deleted {
		var out output
		out.Deleted = fe.ForkUrl
		for _, pe := range pushEvents {
			if pe.Repo == fe.ForkRepo {
				if len(pe.Commits) <= 0 {
					continue
				}
				for _, c := range pe.Commits {
					out.Commits = append(out.Commits, outputCommit{
						Message: c.Message,
						Url:     fmt.Sprintf("https://github.com/%s/commit/%s", fe.OrigRepo, c.Sha),
					})
				}
			}
		}
		if len(out.Commits) > 0 {
			data, err := json.MarshalIndent(&out, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(data))
		}

	}
}

var client = http.Client{
	Timeout: 3 * time.Second,
}

func isNotFound(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		log.Print(err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusNotFound
}

func getForkandPushEvents(filename string) ([]Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []Event

	dec := json.NewDecoder(file)
	for {
		var event Event
		if err := dec.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			return events, err
		}
		if event.Type == "PushEvent" || event.Type == "ForkEvent" {
			events = append(events, event)
		}
	}

	return events, nil
}
