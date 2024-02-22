The scratch base image used in the second stage is *completely empty* so it contains only what we copy into it.

```
docker build -t fibspin .     # image has less than 2MB (as of 2024-02)
docker run fibspin {0..50}    # this will put some load on CPUs
```

```
docker login # just once

docker image tag fibspin reisinge/fibspin
docker image push reisinge/fibspin

k create job fibspin --image=reisinge/fibspin -- /fibspin {0..50}
stern fibspin
k delete job fibspin
```
