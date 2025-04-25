```sh
go run main.go
```

```sh
while true; do echo -n "$(date +%H:%M:%S) "; curl --silent localhost:3000/api/limit; sleep 0.1; done    # sometimes Too Many Requests (429)
while true; do echo -n "$(date +%H:%M:%S) "; curl --silent localhost:3000/api/nolimit; sleep 0.1; done  # all ok (200)
```