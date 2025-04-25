```sh
go run main.go
```

```sh
curl localhost:3000 -v                              # 200
curl localhost:3000/blah -v                         # 404
curl localhost:3000/articles/20250425-demo -v       # 404
curl localhost:3000/api/articles/20250425-demo -v   # 200
```