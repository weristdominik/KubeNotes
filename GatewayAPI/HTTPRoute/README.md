# HTTPRoute

Curl commands:

(make sure [minikube forwarding & tunnel](../README.md) is active)

```bash
curl -i -H "Host: cafe.example.com" http://localhost:80

curl -i -H "Host: cafe.example.com" -H "version: v2" http://localhost:80

curl -i -H "Host: tea.example.com" http://localhost:80
```