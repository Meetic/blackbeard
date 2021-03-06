---
title: "Templates"
anchor: "templates"
weight: 41
---

Playbooks expect `templates` and a *default* inventory.

Templates are very simple. The only constraints are :

* Templates must contains a valid Kubernetes manifest (yaml)
* Template files must have a `.tpl` extension

**Example** : `api.yml.tpl`

```yaml
---
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: api
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
        - name: api
          image:  myCompany/MyApp:{{.Values.api.version}}
          args: ["-Xms{{.Values.api.memoryLimit}}", "-Xmx{{.Values.api.memoryLimit}}", "-Dconfig.resource=config.conf"]
          imagePullPolicy: Always
---
kind: Service
apiVersion: v1
metadata:
  name: api
spec:
  selector:
    app: api
  ports:
    - protocol: TCP
      port: 8080
```

{{% block tip %}}
Under the hood, Blackbeard use [Go templating system](https://golang.org/pkg/text/template/).

Blackbeard compile templates using the content of the inventory file. Thus, two variables are available inside the template :

* `.Values` : contains a json object
* `.Namespace` : contains a string

You can also use this custom functions inside template :

* `getFile "somefile.yml"` : return file content in string
* `sha256sum` : return sha256 hash of a string

{{% /block %}}


