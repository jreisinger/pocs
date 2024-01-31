Tls-tcp demonstrates how TLS works at TCP (or socket) layer.

```
mkcert localhost # https://github.com/FiloSottile/mkcert
go run server/echo.go -cert localhost.pem -key localhost-key.pem

# in second terminal
go run client/main.go -cert localhost.pem
go run certinfo/main.go -insecure
go run certinfo/main.go -addr example.com:443
```

Source: https://eli.thegreenplace.net/2021/go-socket-servers-with-tls/
