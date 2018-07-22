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