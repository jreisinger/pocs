Ratelim contains HTTP server and client for testing rate limiting.

```
$ go install server/ratelims.go 
$ go install client/ratelimc.go
```

```
$ ratelims &
[1] 53298
2023/06/27 10:32:10 listening at localhost:8000
$ ratelimc
1   10:32:13.620281 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=10, reset=0
2   10:32:13.713355 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=9, reset=0
3   10:32:13.813883 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=8, reset=0
4   10:32:13.912703 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=7, reset=0
5   10:32:14.013696 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=6, reset=0
6   10:32:14.112829 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=5, reset=0
7   10:32:14.214347 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=4, reset=0
8   10:32:14.314252 200 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=3, reset=0
9   10:32:14.412203 429 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=2, reset=0
10  10:32:14.512998 429 RateLimit-Policy: 10;w=1;burst=10 RateLimit: limit=10, remaining=2, reset=0
1.00410925s
8   200 OK
2   429 Too Many Requests
```