# Metricsserver

https://kubernetes-sigs.github.io/metrics-server/

## Disable CA verification

--kubelet-insecure-tls - Do not verify the CA of serving certificates presented by Kubelets. For testing purposes only.

```
kubectl edit deployment metrics-server -n kube-system#


args:
          - --kubelet-insecure-tls
```

