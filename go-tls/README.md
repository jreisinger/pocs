HTTP over TLS demo in Go. Stolen from [Liz Rice](https://github.com/lizrice/secure-connections) :-).

certs:

```
# Create CA private key and certificate. Plus generate private key and a signed
# certificate for localhost.
minica -ca-key cakey.key -ca-cert cacert.crt -domains localhost

# You can use github.com/pete911/certinfo to check the generated certificates.
certinfo cacert.crt
certinfo localhost/cert.pem
```

server:

```
go run server/main.go
```

client:

```
curl https://localhost:8080 # observe the client and server error
curl https://localhost:8080 --cacert cacert.crt
```