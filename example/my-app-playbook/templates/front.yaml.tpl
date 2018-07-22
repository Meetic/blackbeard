apiVersion: v1
kind: Service
metadata:
  name: blackbeard-example-web
  labels:
    app: blackbeard-example-web
spec:
  type: NodePort
  ports:
    - port: 8080
      targetPort: 8080
      protocol: TCP
      name: blackbeard-example-web
  selector:
    app: blackbeard-example-web
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blackbeard-example-web
  labels:
    app: blackbeard-example-web
spec:
  replicas: 1
  selector:
    matchLabels:
      app: blackbeard-example-web
  template:
    metadata:
      labels:
        app: blackbeard-example-web
    spec:
      containers:
      - name: blackbeard-example-web
        image: seblegall/blackbeard-example-front:{{.Values.front.version}}
        ports:
        - containerPort: 8080
