The scratch base image is empty so it contains only what we copy there.

```
$ docker build -t fibspin .     # image has less than 2MB (as of 2024-02)
$ docker run fibspin {0..50}    # this will put some load on CPUs
```