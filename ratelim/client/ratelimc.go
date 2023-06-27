/*
Ratelimc is an HTTP client for testing a rate limiting server. For each request
it reports the time, status code, RateLimit-Policy and RateLimit headers.
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"time"
)

var (
	n = flag.Int("n", 10, "number of requests")
	r = flag.Float64("r", 10, "requests rate")
	t = flag.Duration("t", time.Second, "time interval")
	u = flag.String("u", "http://localhost:8000", "URL")
)

func main() {
	flag.Parse()

	statuses := make(map[string]int)

	limiter := time.Tick(time.Duration(float64(*t) / float64(*r)))
	ch := make(chan string)

	start := time.Now()
	go func() {
		for i := 0; i < *n; i++ {
			<-limiter
			get(*u, i, statuses, ch)
		}
	}()
	for i := 0; i < *n; i++ {
		fmt.Println(<-ch)
	}
	fmt.Println(time.Since(start))
	printStats(statuses)
}

func printStats(statuses map[string]int) {
	keys := make([]string, 0, len(statuses))
	for k := range statuses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%-3d %s\n", statuses[k], k)
	}
}

func get(url string, i int, statuses map[string]int, c chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprintf("%-3d %s %v", i+1, time.Now().Format("15:04:05.000000"), err)
		return
	}
	defer resp.Body.Close()
	if _, ok := statuses[resp.Status]; !ok {
		statuses[resp.Status] = 1
	} else {
		statuses[resp.Status]++
	}
	rl := resp.Header.Get("RateLimit")
	rlp := resp.Header.Get("RateLimit-Policy")
	c <- fmt.Sprintf("%-3d %s %d RateLimit-Policy: %s RateLimit: %s",
		i+1, time.Now().Format("15:04:05.000000"), resp.StatusCode, rlp, rl)
}
