# Cilium CNI

https://docs.cilium.io/en/stable/installation/k8s-install-helm/

```bash
helm repo add cilium https://helm.cilium.io/

helm upgrade --install cilium cilium/cilium \
  --version 1.18.5 \
  --namespace kube-system \
  --set ipv6.enabled=true \
  --set kubeProxyReplacement=true \
  --set ipam.mode=cluster-pool \
  --set ipam.operator.clusterPoolIPv6PodCIDRList="{fd00:10:244::/56}" \
  --set ipam.operator.clusterPoolIPv6MaskSize=64
```

*Do this if clusterPoolIPv4/6CIDRList is not working*

delete _kind: ciliumnode_ and recreate

```bash
kubectl delete ciliumnode --all

kubectl delete pod -n kube-system -l k8s-app=cilium

kubectl get ciliumnode -o yaml
```
