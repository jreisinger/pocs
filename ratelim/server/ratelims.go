/*
Ratelims is an HTTP server that rate limits the requests. It's meant for testing
clients that want to interact with a rate limiting server.

It uses the [token bucket] algorithm. It works like a bucket that holds tokens.
Each token represents the permission to serve an HTTP request. When a request is
served a token is taken from the bucket. If the bucket is not full tokens are
added at the fixed rate defined as requests per time interval. Bucket's size is
the maximum number of requests that can be served immediately, i.e. in a burst.

The server sets RateLimit-Policy and (partially) RateLimit headers as defined in
[RateLimit header fields for HTTP] draft.

[token bucket]: https://www.krakend.io/docs/throttling/token-bucket/
[RateLimit header fields for HTTP]: https://ietf-wg-httpapi.github.io/ratelimit-headers/draft-ietf-httpapi-ratelimit-headers.html
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var (
	r = flag.Float64("r", 10, "requests rate")
	t = flag.Duration("t", time.Second, "requests rate time interval")
	b = flag.Int("b", 10, "burst size")
)

func main() {
	flag.Parse()

	limit := rate.Every(time.Duration(float64(*t) / float64(*r)))
	limiter := rate.NewLimiter(limit, *b)

	rlpHeader := fmt.Sprintf("%.0f;w=%.0f;burst=%d;policy='token bucket'", *r, t.Seconds(), *b)

	mux := http.NewServeMux()
	mux.HandleFunc("/", ok)
	addr := "localhost:8000"
	log.Printf("listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, ratelimit(mux, limiter, rlpHeader)))
}

func ok(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func ratelimit(next http.Handler, limiter *rate.Limiter, rlpHeader string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow := limiter.Allow()

		w.Header().Add("RateLimit-Policy", rlpHeader)
		// NOTE: reset value is missing; it is required by the draft.
		rlHeader := fmt.Sprintf("limit=%d, remaining=%v",
			nonNegativeInteger(float64(limiter.Limit())), nonNegativeInteger(limiter.Tokens()))
		w.Header().Add("RateLimit", rlHeader)

		if allow {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})
}

func nonNegativeInteger(f float64) int {
	if f < 0 {
		return 0
	}
	return int(f)
}
