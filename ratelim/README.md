Ratelim contains HTTP server and client for testing rate limiting.

```
$ go run server/ratelims.go # -h for help
2023/06/27 22:35:06 listening at localhost:8000

$ go run client/ratelimc.go -r 100 -n 11
1   22:36:03.583733 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=9
2   22:36:03.588827 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=8
3   22:36:03.599635 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=7
4   22:36:03.609698 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=6
5   22:36:03.620019 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=5
6   22:36:03.629796 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=4
7   22:36:03.639778 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=3
8   22:36:03.650108 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=2
9   22:36:03.659562 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=1
10  22:36:03.669362 200 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=0
11  22:36:03.679855 429 RateLimit-Policy: 10;w=1;burst=10;policy='token bucket' RateLimit: limit=10, remaining=0
112.602625ms
10  200 OK
1   429 Too Many Requests
```

More

- https://gobyexample.com/rate-limiting
- https://www.alexedwards.net/blog/how-to-rate-limit-http-requests
- [counting semaphore](https://github.com/adonovan/gopl.io/blob/master/ch8/crawl2/findlinks.go)
