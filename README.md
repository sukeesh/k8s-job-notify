# Kubernetes Job/CronJob Notifier 

This tool sends an alert to slack whenever there is a [Kubernetes](https://github.com/kubernetes/kubernetes) cronJob/Job failure/success.

**No extra setup** required to deploy this tool on to your cluster, just apply below K8s deploy manifest 🎉  

This uses `InClusterConfig` for accessing Kubernetes API.

Limitations
-----
- Namespace scoped i.e., each namespace should have this deploy *separately*
- **All the jobs** in the namespace are fetched and verified for failures
  - Will add support for selectors in future 📋 
   
Development
----
If you wish to run this locally, clone this repository, set `webhook` and `namespace` env variables  
```$xslt
$ export webhook="slack_webhook_url" && export namespace="<namespace_name>" && go build &&  ./k8s-job-notify
```

Docker 🐳
--- 
Docker images are hosted at [hub.docker/k8s-job-notify](https://hub.docker.com/r/sukeesh/k8s-job-notify)

To start using this
---

Create and apply below kubernetes deployment in your cluster
```$xslt
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: kjn
  name: k8s-job-notify
  namespace: <namespace_name>
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kjn
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "false"
      labels:
        app: kjn
    spec:
      restartPolicy: Never  
      containers:
      - env:
        - name: webhook
          value: <slack_webhook_url> # creating a secret for this var is recommended
        - name: namespace
          valueFrom:
            fieldRef:
                fieldPath: metadata.namespace
        - name: incluster
          value: "1"
        image: sukeesh/k8s-job-notify:latest
        name: k8s-job-notify
        resources:
          limits:
            cpu: 500m
            memory: 256Mi
          requests:
            cpu: 500m
            memory: 128Mi
```

 