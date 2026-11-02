An example of gRPC from the network automation domain.

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative networkstuff/*.proto

go run ./server
go run ./client # in 2nd terminal
```

For more
- https://grpc.io/docs/languages/go/quickstart/
- https://github.com/oreilly-npa-book/examples/tree/main/ch10-apis/grpc
