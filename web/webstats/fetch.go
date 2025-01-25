package webstats

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type fetch struct {
	time time.Duration
	size int64
}

func doFetch(url string) (fetch, error) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return fetch{}, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fetch{}, fmt.Errorf("reading body: %v", err)
	}
	return fetch{size: n, time: time.Since(start)}, nil
}
