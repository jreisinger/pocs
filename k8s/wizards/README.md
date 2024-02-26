Wizards demonstrates a simple web site that is:

- distributed as [distroless](https://github.com/GoogleContainerTools/distroless) container
- deployable to Kubernetes via Helm

Build and run locally:

```
docker build -t wizards .
docker run --name wizards -p 8080:8080 -d wizards

curl localhost:8080
docker rm -f wizards
```

Interactively debug:

```
docker build -t wizards-debug -f Dockerfile-debug .
docker run -it --entrypoint=sh wizards-debug
/ # /wizards
```

Push to image registry:

```
docker image tag wizards reisinge/wizards
docker image push reisinge/wizards
```

Deploy to Kubernetes cluster:

```
helm template . | kubectl apply -f -

kubectl run tmp --image=busybox --rm -it --restart=Never -- wget wizards:8080 -qO-
helm template . | kubectl delete -f -
```