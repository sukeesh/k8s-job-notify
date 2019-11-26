# Kubernetes Job/CronJob Notifier

Sends an alert when a Job/CronJob succeeds/fails to Slack. This runs every 10 minutes to check for any success/failures of Kubernetes Jobs/CronJobs.

Docker images are hosted at [hub.docker.com/r/sukeesh/k8s-job-notify](https://hub.docker.com/r/sukeesh/k8s-job-notify)

Docker pull command 
```$xslt
$ docker pull sukeesh/k8s-job-notify
```
Usage
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
          value: <slack_webhook_url>
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
 
 TODO
 ---
 - Add support for labels
 - Create a CRD
 