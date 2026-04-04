# cert-manager

## Installation

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.20.0/cert-manager.yaml
```

## CA

```yaml
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: kubenet-local-ca
  namespace: cert-manager
spec:
  isCA: true
  duration: 43800h # 5 years
  commonName: kubenet.local
  secretName: kubenet-local-key-pair
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: selfsigned
    kind: ClusterIssuer
    group: cert-manager.io
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: kubenet-local-ca
spec:
  ca:
    secretName: kubenet-local-key-pair
EOF
```
