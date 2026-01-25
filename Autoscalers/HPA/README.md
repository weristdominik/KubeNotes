# HPA

## Load Generator for HPA

`kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -n hpa -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://hpa-svc.hpa.svc.cluster.local; done"`