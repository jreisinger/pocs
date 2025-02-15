```
mkcert localhost

go run ./webserver/main.go
curl https://localhost:44300 --cacert localhost.pem

go run ./l4/main.go
curl https://localhost:44300 --cacert localhost.pem

go run ./l7/main.go
curl https://localhost:44300 --cacert localhost.pem
```