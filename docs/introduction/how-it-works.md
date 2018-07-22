---
title: "How it works?"
anchor: "how-it-works"
weight: 12
---
![how it works?](/img/blackbeard_mechanism.png)

## Playbooks

Blackbeard use *playbooks* to manage namespaces. A playbook is a collection of kubernetes manifest describing your stack, written as templates. A playbook also require a `default.json`, providing the default values to apply to the templates.

Playbooks are created as files laid out in a particular directory tree.