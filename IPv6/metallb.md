# MetalLB

https://metallb.universe.tf/installation/

https://metallb.universe.tf/configuration/

```bash
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.3/config/manifests/metallb-native.yaml

kubectl apply -f values.yaml
```

## values.yaml

```bash
# kubernetes/metallb/values.yaml
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: ipv6-pool
  namespace: metallb-system
spec:
  addresses:
  - "2a01:xxxx:xxxx:xxxx:0001:0000:0000:0100/128"
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: default-ipv6
  namespace: metallb-system
spec:
  ipAddressPools:
  - ipv6-pool
  interfaces:
  - enp1s0
```