---
title: "defaults.json"
date: 2018-07-19T12:42:17+01:00
anchor: "defaults.json"
weight: 33
---

The `defaults.json` file is the default inventory. It contains defaults values to apply on the templates. Blackbeard uses it to generated per namespace inventory.

The only constraint on the `default.json` file is :

* Must contains a `namespace` key, containing the value "default".

**Example :** `defaults.json`

```json
{
    "namespace": "default",
    "values": {
        "api": {
            "version": "1.0.0",
            "memoryLimit": "128m"
        },
        "front": {
            "version": "1.0.0"
        }
    }
}
```