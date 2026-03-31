# MongoDB

## Installation

### Operator

```bash
# CRD
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-kubernetes/1.7.0/public/crds.yaml

# Create ns
kubectl create ns mongodb

# Operator(s)
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-kubernetes/1.7.0/public/mongodb-kubernetes.yaml

kubectl get all -n mongodb
```

### Quick local-storage dynamic provider

```bash
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml

# mark sc as default
kubectl patch storageclass local-path \
  -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```


### Create an MongoDB Database

```yaml
cat <<EOF | kubectl apply -f -
apiVersion: mongodbcommunity.mongodb.com/v1
kind: MongoDBCommunity
metadata:
  name: example-mongodb
  namespace: mongodb
spec:
  members: 1  # up to 3
  type: ReplicaSet
  version: "6.0.5"
  security:
    authentication:
      modes: ["SCRAM"]
  users:
    - name: my-user
      db: admin
      passwordSecretRef: # a reference to the secret that will be used to generate the user's password
        name: my-user-password
      roles:
        - name: clusterAdmin
          db: admin
        - name: userAdminAnyDatabase
          db: admin
      scramCredentialsSecretName: my-scram
  additionalMongodConfig:
    storage.wiredTiger.engineConfig.journalCompressor: zlib
---
apiVersion: v1
kind: Secret
metadata:
  name: my-user-password
  namespace: mongodb
type: Opaque
stringData:
  password: my-password
EOF
```

used images

```bash
kubectl get pods --all-namespaces -o jsonpath="{.items[*].spec['initContainers', 'containers'][*].image}" |\
tr -s '[[:space:]]' '\n' |\
sort |\
uniq -c | grep mongo
   1 quay.io/mongodb/mongodb-agent:108.0.2.8729-1
   1 quay.io/mongodb/mongodb-community-server:6.0.5-ubi8
   1 quay.io/mongodb/mongodb-kubernetes-operator-version-upgrade-post-start-hook:1.0.10
   1 quay.io/mongodb/mongodb-kubernetes-readinessprobe:1.0.23
   1 quay.io/mongodb/mongodb-kubernetes:1.7.0
```