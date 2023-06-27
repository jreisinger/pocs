/*
Ratelims is an HTTP server that rate limits the requests. It's meant for testing
clients that want to interact with a rate limiting server.

It uses the [token bucket] algorithm. It works like a bucket that holds tokens.
Each token represents the permission to serve an HTTP request. When a request is
served a token is taken from the bucket. If the bucket is not full tokens are
added at the fixed rate defined as requests per second. Bucket's size is the
maximum number of requests that can be served immediately, i.e. in a burst.

The server sets RateLimit-Policy and RateLimit headers as defined in [RateLimit
header fields for HTTP] draft.

[token bucket]: https://www.krakend.io/docs/throttling/token-bucket/
[RateLimit header fields for HTTP]: https://ietf-wg-httpapi.github.io/ratelimit-headers/draft-ietf-httpapi-ratelimit-headers.html
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

var (
	r = flag.Float64("r", 10, "requests per second")
	b = flag.Int("b", 10, "burst size")
)

func main() {
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", ok)
	addr := "localhost:8000"
	log.Printf("listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, ratelimit(mux, *r, *b)))
}

func ok(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func ratelimit(next http.Handler, rps float64, burst int) http.Handler {
	var limiter = rate.NewLimiter(rate.Limit(rps), burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rlm := fmt.Sprintf("%.0f;w=1;burst=%d", rps, burst)
		w.Header().Add("RateLimit-Policy", rlm)

		rm := fmt.Sprintf("limit=%.0f, remaining=%.0f, reset=%.0f",
			rps, limiter.Tokens(), limiter.Reserve().Delay().Seconds())
		w.Header().Add("RateLimit", rm)

		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
