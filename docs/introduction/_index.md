---
title: "What is Blackbeard?"
anchor: "introduction"
weight: 10
---

**Blackbeard is a namespace manager for Kubernetes.** :thumbsup: It helps you to develop and test with Kubernetes using namespaces.

{{% block tip %}}
Kubernetes namespaces provide an easy way to isolate your components, your development environment, or your staging environment.
{{% /block %}}

Blackbeard helps you to deploy your Kubernetes manifests on multiple namespaces, making each of them running a different version of your microservices. You may use it to manage development environment (one namespace per developer) or for testing purpose (one namespace for each feature to test before deploying in production).

## Purpose

*When working in a quite large team or, in a quite large project, have you ever experienced difficulties to test multiple features at the same time?*

Usually, teams have 2 alternatives :

* Stack "features to test" in a queue and wait for the staging environment to be available;
* Try more or less successfully to create and maintain an "on demand" staging environment system, where each environment is dedicated to test a specified feature.

Blackbeard helps you to create ephemeral environments where you can test a set of features before pushing in production.

It also provide a very simple workflow, making things easier if you plan to run automated end-to-end testing.

## How it works ?

![how it works?](/img/blackbeard_mechanism.png)

## Playbook

Blackbeard use *playbooks* to manage namespaces. A playbook is a collection of kubernetes manifest describing your stack, written as templates. A playbook also require a `default.json`, providing the default values to apply to the templates.

Playbooks are created as files laid out in a particular directory tree.