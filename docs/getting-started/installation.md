---
title: "Installation"
anchor: "installation"
weight: 21
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
