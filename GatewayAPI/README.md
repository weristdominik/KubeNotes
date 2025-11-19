# Gateway API

## References

## Installation

```bash
minikube start
```

```bash
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.0/experimental-install.yaml
```

```bash
helm pull oci://ghcr.io/nginxinc/charts/nginx-gateway-fabric --untar
cd nginx-gateway-fabric

kubectl create namespace nginx-gateway

helm install ngf . -n nginx-gateway
```

```bash
kubectl kustomize HTTPRoute/ | kubectl apply -f -
```

```bash
kubectl get svc ngf-nginx-gateway-fabric 
# New Terminal window
kubectl port-forward deployment/ngf-nginx-gateway-fabric 8080:80 8443:443 -n nginx-gateway
# Antoher new Terminal window
minikube tunnel
```

Curl:
```bash
curl -i -H "Host: domain.example.com" http://localhost:80
```
