---
title: "Inventories"
date: 2018-07-19T12:42:17+01:00
anchor: "inventories"
weight: 32
---

An inventory is a file containing a json object. This object is made available inside the templates by Blackbeard.

The only constraints on an inventory json object are :

* An inventory must contains a `namespace` key, containing a string. The string value must be the name of the namespace the inventory is associated to;
* An inventory must contains a `values` key. This key may contains whatever you want.
* An inventory file must be located in the `inventories` directory
* An inventory file must be named after the template : `{{namespace}}_inventory.json`

**Example** : `john_inventory.json`

```json
{
    "namespace": "john",
    "values": {
        "api": {
            "version": "1.2.0",
            "memoryLimit": "128m"
        },
        "front": {
            "version": "1.2.1"
        }
    }
}
```

Inventories are generated from the `defaults.json` file. Blackeard copy the `defaults.json` file content, create a inventory for the given namespace (located in the `inventories` directory), past the content default values and change the `namespace` key value with the corresponding namespace