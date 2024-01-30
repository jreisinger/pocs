[Distroless images](https://github.com/GoogleContainerTools/distroless) are smaller and thus:

* faster to start
* faster to upload and download
* more secure (smaller attack surface, better signal/noise ratio of scanners)

```sh
docker build -t distroless .
docker run --name distroless -p 8080:8080 -d distroless
curl localhost:8080
docker rm -f distroless

# To debug add :debug tag in the final FROM in Dockerfile.
docker build -t distroless-debug .
docker run -it --entrypoint=sh distroless-debug
/ # /app
2024/01/30 12:24:00 starting a web server on port 8080
```
