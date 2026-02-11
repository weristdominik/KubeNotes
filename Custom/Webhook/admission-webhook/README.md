# admission webhook

```bash
go mod init admission-webhook
mkdir webhook

cd ..
go get k8s.io/api/core/v1
go get k8s.io/apimachinery/pkg/apis/meta/v1

mkdir certs
openssl req -newkey rsa:2048 -nodes -keyout certs/tls.key \
    -x509 -days 365 -out certs/tls.crt -subj "/CN=localhost"
```

Create Certificate

Since [Kubernetes v1.23 we need SANs](https://docs.cloud.google.com/kubernetes-engine/docs/deprecations/webhookcompatibility#:~:text=Instead%2C%20Kubernetes%20will%20only%20rely,resource%20(i.e.%20Pod)%20creation.)

```bash
mkdir certs
cd certs

openssl req -x509 -nodes -days 365 \
  -newkey rsa:2048 \
  -keyout tls.key \
  -out tls.crt \
  -config openssl.cnf \
  -extensions req_ext
```

build docker 

> if using minikube: 'eval $(minikube docker-env)' BEFORE running docker build command

```bash
# Same dir as Dockerfil
docker build -t admission-webhook:v0.9 .
```

deploy to cluster

```bash
 kubectl apply -f kubernetes/deployment.yaml
 kubectl apply -f kubernetes/admissionWebhookConfiguration.yaml
```

mock to service

```bash
curl -k -X POST https://localhost:8443/validate \
  -H "Content-Type: application/json" \
  -d @AdmissionReview-example.json
```

If a new Pod gets deployed an ["AdmissionReview" config](../AdmissionReview-example.json) is sent

If someone creates or updates an Pod without `label: zone`:

```bash
kubectl run whoami --image=traefik/whoami:v1.11 -n default --labels=run=whoami,env=test 
Error from server: admission webhook "pod-validator.default.svc.cluster.local" denied the request: Pod is missing 'zone' label

kubectl edit pod whoami
error: pods "whoami" could not be patched: admission webhook "pod-validator.default.svc.cluster.local" denied the request: Pod is missing 'zone' label
You can run `kubectl replace -f /var/folders/qb/xqnq2s7d76zf6sz1925558p80000gn/T/kubectl-edit-621644048.yaml` to try this update again.

# correct usage
kubectl run whoami --image=traefik/whoami:v1.11 -n default --labels=zone=local,run=whoami,env=test 
```
