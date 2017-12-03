# Blackbeard

Blackbeard is a tool that let you manage multiple version of the same stack using kubernetes and namespaces.

If you need to deploy and use mutiple version of a stack (a group of applications), Blackbeard is the tool you need.

Blackbeard is made to be executed using a directory containing configuration files and directories called a *Playbook*.

A *Playbook* is a directory that contains :

* A `defaults.json` file that defines the default values to apply.
* A `templates` directory that contains the configuration templates. Those are typically kubernetes configuration files (yaml).
* An `inventories` directory that will contains the future inventories (One per namespace). The content of this directory should not be versioned.
* A `configs` directory that will contains the future configuration files (one sub-dir per namespace). The content of this directory should not be versioned as well.

By default, Blackbeard will try to use the current directory as a Playbook. You can also specify a default playbook using a configuration file (see configuration).

## Usage

### Requirement

* A working and configure kubectl

### Cli usage

Blackbeard can be use as a cli tool.

You can find examples bellow.

#### Creatin a new env :

```sh
cd {your playbook}
kubgen create -n {namespace name}
```

#### Apply a change :

```sh
cd {your playbook}/inventories
## edit the inventory file you want to update
cd ..
kubgen apply -n {namespace name}
```

#### Getting Help

```sh
blackbeard -h
blackbeard create -h
blackbeard apply -h
```

### REST API / websocket server

blackbeard also provide a webserver able to handle REST queries and a websocket server.

The REST api can be used to do the same thing that you can do using blackbeard as cli tool.

Here are example  : 

#### Create a new namespace

`POST /inventories/`

```json
{
    "namespace": "foo"
}
```

#### Get defaults

`GET /defaults/`

#### Get a namespace

`GET /inventories/foo`

#### Get the list of avaible namespace

`GET /inventories`

#### Update a namespace

`PUT /inventories/foo`

```json
{
    "namespace" : "bar",
    "values": {
        "microservice" : {
            "name": "baz",
            "version" : "1.2.3"
        }
    }
}
```

*Note : the payload here must match the required defaults.json (see creating Paybook bellow)*

## Playbooks

A `playbook` is a bunch of kubernetes configuration files using some variables that may be different depending on what we want to do.

Under the hood, a `playbook` is a suit of templates containing references to variables. The value of those variables are defined in an `inventory` file.

An `inventory` is a file containing a JSON object with at leat to fields :

* `namespace`
* `values`

Inventories are not part of a `playbook`. They are generated from a `defaults.json` file that is part of the playbook and define :

* The structure of the values to be used in the templates
* The default values.

Since templates are Go template, you just have to follow [the rules of the go template engine](https://golang.org/pkg/text/template/). It is very similar to any well known template engine (Jinja,twig, etc.)

There are some rules you need to follow.

### defaults rules

`defaults.json` must contains a JSON object with at leat to fields :

* `namespace`
* `values`

The `namespace` must be a string. It's a common practice to set the value to "default". This value will be replaced for each new generated inventory.

The `values` can be what ever you want : a JSON object, an array of JSON object, etc.

### templates rules

* Templates files must be located in the `templates` directory.
* All templates files must end with `.tpl` extension.

### values rules

You can put whatever you want in the default `values` field. The objects you put in this field will be used to execute the template.

Example :
If you choose to defines values like that : 

```json
{
    "namespace": "test",
    "values": {
        "apis": [
            {
                "name": "test",
                "url": "http://test.kube",
                "version": "1.2.3"
            }
        ]
    }
}
```

You will be able to use those values in the template following the template engine rules : 

```
{{range .Values.apis}}
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
          image:  artifact-docker-rd.meetic.ilius.net/meetic/api-private/{{.name}}:{{.version}}
          imagePullPolicy: Always
{{end}}

```

**Caution :** First letter of `namespace` and `value` must be capitalized when called in the template.

## Configuration

Using a configuration file is not mandatory. But if you use the webserver or always work with the same playbook, it may be easier.

Create a new file `~/.blackbeard.yml` following this guide line : 

```yaml
working-dir: /path/to/your/playbook
```

## Installation

### Requirements
* go > 1.8

### Go installation

On Linux, follow the white rabbit : [https://golang.org/doc/install](https://golang.org/doc/install)

Then, you need to configure what is called a "workspace".

By default, the workspace is `$HOME/go`.

If you want to use a different one, you have to set up you GOPATH env var. : [https://github.com/golang/go/wiki/Setting-GOPATH](https://github.com/golang/go/wiki/Setting-GOPATH)

Last thing : if you want your binary to be executed from anywhere you should also add `$GOPATH/bin` to the PATH env var.

### Dep installation

This tool is used to get all dependencies needed to build the project.

```sh
go get -u github.com/golang/dep/cmd/dep
```

Or, you can follow the installation instruction at [https://github.com/golang/dep](https://github.com/golang/dep)

### Build

Getting dependencies :

```sh
dep ensure
```

Generating binary :

```sh
go build
```

(Here, you will be able to execute the binary using `./blackbeard`)

Installing binary:

```sh
go install
```