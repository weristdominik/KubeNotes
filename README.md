# KubeNotes

---

![Kubernetes Architecture](https://upload.wikimedia.org/wikipedia/commons/6/67/Kubernetes_logo.svg)


---

## Every great project starts with "I`ll clean this up later."

### Debug Container

(Gets deleted once exited)

`kubectl run busybox --image=busybox --rm -it --restart=Never -- /bin/sh`

`kubectl run curl --image=curlimages/curl --rm -it --restart=Never -- -i -k -H "Authorization: Bearer <ACCESS_TOKEN>" app1.app1-ns.svc.cluster.local/app1/`