# Export Your Keycloak config json

http://keycloak.com/admin

Realm-settings > Actions (top right) > Partial export > Include groups and roles ON + Include clients ON > Export

Store file at `cilium-l7-oidc/realm-export.json`

create configmap for [keycloak.yaml](../keycloak.yaml)

```bash
kubectl create configmap keycloak-realm --from-file=realm.json=realm-export.json -n keycloak-ns
```