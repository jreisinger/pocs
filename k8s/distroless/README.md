[Distroless images](https://github.com/GoogleContainerTools/distroless) are smaller and thus:

* faster to start
* faster to upload and download
* more secure (smaller attack surface, better signal/noise ratio of scanners)

```sh
docker build -t wizards .
docker run --name wizards -p 8080:8080 -d wizards
curl localhost:8080
docker rm -f wizards

# To interactively debug.
docker build -t wizards-debug -f Dockerfile-debug .
docker run -it --entrypoint=sh wizards-debug
/ # /wizards
2024/01/30 12:24:00 starting a web server on port 8080
```
