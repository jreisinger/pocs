[Distroless images](https://github.com/GoogleContainerTools/distroless) are:

* more secure (smaller attack surface, better signal/noise ratio of scanners)
* smaller

Usage:

```sh
docker build -t distroless .
docker run distroless
```

Debug:

```sh
# Add :debug tag in the final FROM in Dockerfile.
docker build -t distroless .
docker run -it --entrypoint=sh distroless
```