```sh
go run main.go
```

```sh
curl localhost:3000/articles/20250425-demo -v # 200
curl localhost:3000/articles/20250425-blah -v # 404
curl localhost:3000/articles/25-04-25-demo -v # 422
```