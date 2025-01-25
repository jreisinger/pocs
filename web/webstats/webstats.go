// Package web fetches web resources to get some statistics about them.
package webstats

import (
	"fmt"
	"time"
)

type Stats struct {
	URL       string
	Err       error
	FetchSize int64
	FetchTime time.Duration
}

func (s Stats) String() string {
	return fmt.Sprintf("%.3fs %9s %s",
		s.FetchTime.Seconds(),
		HumanReadableSize(s.FetchSize),
		s.URL,
	)
}

func Get(urls []string) []Stats {
	ch := make(chan Stats)

	for _, u := range urls {
		go func(url string) {
			fetch, err := doFetch(url)
			ch <- Stats{URL: url, Err: err, FetchSize: fetch.size, FetchTime: fetch.time}
		}(u)
	}

	var stats []Stats
	for range urls {
		resource := <-ch
		stats = append(stats, resource)
	}
	return stats
}

func HumanReadableSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
