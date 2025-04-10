Simpleweb demonstrates a simple web site that is distributed as a [distroless](https://github.com/GoogleContainerTools/distroless) container image and deployed to Kubernetes via Helm.

Build and run locally:

```
docker build -t wizards .
docker run --name wizards -p 1212:8080 -d wizards

curl localhost:1212
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
helm template ./helm | kubectl apply -f -

kubectl run tmp --image=busybox --rm -it --restart=Never -- wget wizards:8080 -qO-

helm template ./helm | kubectl delete -f -
```
