# Contributing

## Installation

### Requirements
* go >= 1.8

### Go installation

On Linux, follow the white rabbit : [https://golang.org/doc/install](https://golang.org/doc/install)

Then, you need to configure what is called a "workspace".

By default, the workspace is `$HOME/go`.

If you want to use a different one, you have to set up you GOPATH env var. : [https://github.com/golang/go/wiki/Setting-GOPATH](https://github.com/golang/go/wiki/Setting-GOPATH)

Last thing : if you want your binary to be executed from anywhere you should also add `$GOPATH/bin` to the PATH env var.

### Dependencies installation

```sh
make dep
```

This simple make target will `go get` all the tool you need to work on Blackbeard

### Build

```sh
make
```

### Tests
```sh
make test
```

## Contributing

All pull request and issue are welcomed. Note that this repo come with Travis-ci integration.
So any PR will have to pass at least the travis job before being merged.

### Convention

The Blackbeard code source tends to follow the Go code convention describe in the [Go code review comments](https://github.com/golang/go/wiki/CodeReviewComments).
