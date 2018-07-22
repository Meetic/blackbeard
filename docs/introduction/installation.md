---
title: "Installation"
date: 2018-07-18T12:42:17+01:00
anchor: "installation"
weight: 16
---
### Recommanded

The simplest way of installing Blackbeard is to use the installation script :

```sh
curl -sf https://raw.githubusercontent.com/Meetic/blackbeard/master/install.sh | sh
```

### Manually

Download your preferred flavor from the releases page and install manually.

### Using Go Get

Note : this method requires Go 1.9+ and dep.

```sh
go get github.com/Meetic/blackbeard
cd $GOPATH/src/github.com/Meetic/blackbeard
make build
```
