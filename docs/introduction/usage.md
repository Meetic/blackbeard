---
title: "Usage"
date: 2018-07-18T12:42:17+01:00
anchor: "usage"
weight: 15
---

Blackbeard provide a CLI interface in addition to a REST API. This way, you can use it either to run automated tests in a CI pipeline or plug your own UI for manuel deployment purpose.

{{% block note %}}
Blackbeard requires `kubectl` to be installed and configured to work.
{{% /block %}}

#### Creating a new isolated env

```sh
blackbeard create -n my-feature
```

This command actually *create a namespace* and generate a JSON configuration file containing default values. This file, called an `inventory`, is where you may update the values to apply specifically to your new namespace (such as microservice version)

#### Applying changes

```sh
blackbeard apply -n my-feature
```

This command apply your kubernetes manifest (modified by the values you have put in the generated `inventory` previously) to your newly created namespace

#### Getting services endpoint / ports

```sh
blackbeard get services -n my-feature
```

This step will prompt a list of exposed services in the namespace. If you need to connect to a database in order to test data insertion, it is where you will find useful info.

#### Getting back to previous state

```sh
blackbeard delete -n my-feature
```

Delete all generated files and delete the namespace.