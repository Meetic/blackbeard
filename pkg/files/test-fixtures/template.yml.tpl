{{range .Values.microservices}}
---
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: {{.name}}
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: fpm-{{.name}}
    spec:
      containers:
        - name: {{.name}}
          image:  docker.io/{{.name}}:{{.version}}
{{end}}