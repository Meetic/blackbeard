# Blackbeard
[![Build Status](https://travis-ci.org/Meetic/blackbeard.svg?branch=master)](https://travis-ci.org/Meetic/blackbeard) [![Go Report Card](https://goreportcard.com/badge/github.com/Meetic/blackbeard)](https://goreportcard.com/report/github.com/Meetic/blackbeard) [![GitHub license](https://img.shields.io/github/license/Meetic/blackbeard.svg)](https://github.com/Meetic/blackbeard/blob/master/LICENSE) 
[![GitHub release](https://img.shields.io/github/release/Meetic/blackbeard.svg)](https://github.com/Meetic/blackbeard) [![Twitter](https://img.shields.io/twitter/url/https/github.com/Meetic/blackbeard.svg?style=social)](https://twitter.com/intent/tweet?text=Wow:&url=https%3A%2F%2Fgithub.com%2FMeetic%2Fblackbeard) 

## Introduction
**Blackbeard is a namespace manager for Kubernetes.** It helps you to develop and test with Kubernetes using namespaces.

Blackbeard helps you to deploy your Kubernetes manifests on multiple namespaces, making each of them running a different version of your microservices. You may use it to manage development environment (one namespace per developer) or for testing purpose (one namespace for each feature to test before deploying in production).

Blackbeard use *playbooks* to manage namespaces. A playbook is a collection of kubernetes manifest describing your stack, written as templates. A playbook also require a `default.json`, providing the default values to apply to the templates.

Playbooks are created as files laid out in a particular directory tree.

## Requirements

You must have `kubectl` installed and configured to use Blackbeard

## Installation

```sh
curl -sf https://raw.githubusercontent.com/Meetic/blackbeard/master/install.sh | sh
```

## Documentation

You can find Blackbeard documentation on the [Blackbeard website](https://blackbeard.netlify.com)

