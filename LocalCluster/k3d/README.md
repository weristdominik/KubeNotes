# k3d

# Download and run installer script
```bash
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

k3d version
```

# Download and install kubectl

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"

chmod +x kubectl
sudo mv -v kubectl /usr/local/bin/
```

# Create a Cluster

Start docker service !

```bash
k3d cluster create traefik \
  --port 80:80@loadbalancer \
  --port 443:443@loadbalancer \
  --port 8000:8000@loadbalancer \
  --k3s-arg "--disable=traefik@server:0"
```

1. [Offical Kubernetes Gateway api crds resources](https://gateway-api.sigs.k8s.io/guides/#install-standard-channel)
2. [Traefik Gateway Crds](https://doc.traefik.io/traefik/reference/install-configuration/providers/kubernetes/kubernetes-gateway/#traefik-kubernetes-with-gateway-api)
3. Traefik helm repo + apply

```bash
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.1/standard-install.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v3.6/docs/content/reference/dynamic-configuration/kubernetes-gateway-rbac.yml
helm repo add traefik https://traefik.github.io/charts
helm repo update

helm install traefik traefik/traefik -f traefik-values.yaml --wait
```

`kubectl get all -n traefik`
