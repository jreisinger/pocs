# mutating-wh

Kubernetes [admission controllers](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) intercept an API request (to `kube-apiserver`) before it is persisted into `etcd`. The controllers can mutate and/or validate the request. 

![](apirequest.png)

Among the compiled-in admission controllers there are three special - MutatingAdmissionWebhook, ValidatingAdmissionWebhook, and ValidatingAdmissionPolicy - that you can use to customize cluster behavior at admission time.

This repo is based on [pete's template](https://github.com/pete911/template-wh) and contains MutatingAdmissionWebhook. Admission webhooks are HTTP callbacks (Kubernetes service in this case) that receive admission requests, do something with them and respond. They are configured (registered) in the API at run time using MutatingWebhookConfiguration. See the repo code for more. To test it end to end:

```
make e2e-test
```
