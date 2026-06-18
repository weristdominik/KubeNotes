# Kind IPv6 Only

Found at: https://oneuptime.com/blog/post/2026-03-20-kind-clusters-ipv6-configuration/view

## Configure docker.service

https://docs.docker.com/engine/daemon/ipv6/

```bash
sudo vi /etc/docker/daemon.json
```

```json
{
  "ipv6": true,
  "fixed-cidr-v6": "fd00:1::/64"
}
```

```bash
sudo systemctl restart docker
```


## Create kind cluster config

```yaml
# kind-ipv6-only.yaml - kind cluster configuration for IPv6-only
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  # IPv6-only cluster
  ipFamily: ipv6
  apiServerAddress: 127.0.0.1
  podSubnet: "fd00:10:244::/56"
  serviceSubnet: "fd00:10:96::/112"
nodes:
  - role: control-plane
  - role: worker
```

## Create the dual-stack cluster

```bash
kind create cluster --config kind-ipv6-only.yaml --name ipv6-only
```

## Quick test

```bash
# Node IPs
kubectl get nodes -o wide

# Start Pod
kubectl run test --image=busybox:1.36 --restart=Never -- sleep 3600
kubectl get pods -o wide

# Get the IPv6 address of the test pod
POD_IPV6=$(kubectl get pod test -o go-template='{{range .status.podIPs}}{{printf "%s\n" .ip}}{{end}}' | grep ':')

# From another pod, ping the IPv6 address
kubectl run pinger --image=busybox:1.36 --restart=Never -- ping -6 -c 3 "$POD_IPV6"

# Check the result
kubectl logs pinger
```
