apiVersion: v1
kind: Service
metadata:
  name: blackbeard-example-app
  labels:
    app: blackbeard-example-app
spec:
  clusterIP: None
  ports:
    - port: 50051
      name: blackbeard-example-app
  selector:
    app: blackbeard-example-app
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blackbeard-example-app
  labels:
    app: blackbeard-example-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: blackbeard-example-app
  template:
    metadata:
      labels:
        app: blackbeard-example-app
    spec:
      containers:
      - name: blackbeard-example-app
        image: seblegall/blackbeard-example-api:{{.Values.api.version}}
        ports:
        - containerPort: 50051
