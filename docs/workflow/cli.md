---
title: "CLI"
date: 2018-07-18T09:14:17+01:00
anchor: "cli"
weight: 21
---
Blackbeard provide a set of CLI command.

CLI usage may be use for different purpose.

Locally, for developers, it is useful for managing multiple namespace either to develop a microservice calling other applications running a specified version, or to test the entire stack using specified version of multiple services in the stack.

In a CI/CD pipeline, for automated end-to-end testing.

### Create a new env

```sh
cd {your playbook}
blackbeard create -n {namespace name}
```

* create a Kubernetes namespace;
* generate a `inventory` file for the newly created namespace;
* generate a set of yml `manifest` based on the playbook `templates`.

### Update values & apply changes

```sh
cd {your playbook}/inventories
## edit the inventory file you want to update
cd ..
blackbeard apply -n {namespace name}
```

* apply values defined in the `inventory` file to the playbook `templates`;
* update the yml `manifest` using the newly updated values from the `inventory` file;
* run a `kubectl apply` command and apply changes in the manifest to the namespace.

### List namespaces

```sh
blackbeard get namespaces
```

* prompt a list of available Kubernetes namespace;

for each namespace :

* indicate if the namespace is managed by a local `inventory` or not.
* indicate the status of the namespace (aka : percentage of pods in a "running" state)

Exemple :

```sh
Namespace	Phase	Status	Managed
backend		Active	100%	false
john    	Active	73% 	true
default		Active	0%    false
kevin   	Active	73%   false
team1	   	Active	73%   true
```

### Get useful informations about services

```sh
blackbeard get services -n my-feature
```

* prompt a list of exposed services

{{% block info %}}
Exposed services are Kubernetes services exposed using `NodePort` or http services exposed via `Ingress`
{{% /block %}}


### Get Help

```sh
Usage:
  blackbeard [command]

Available Commands:
  apply       Apply a given inventory to the associated namespace
  create      Create a namespace and generated a dedicated inventory.
  delete      Delete a namespace
  get         Show informations about a given namespace.
  help        Help about any command
  reset       Reset a namespace based on the template files and the default inventory.
  serve       Launch the blackbeard server

Flags:
      --config string   config file (default is $HOME/.blackbeard.yaml)
      --dir string      Use the specified dir as root path to execute commands. Default is the current dir.
  -h, --help            help for blackbeard

Use "blackbeard [command] --help" for more information about a command.
```