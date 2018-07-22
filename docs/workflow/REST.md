---
title: "HTTP"
date: 2018-07-18T09:14:17+01:00
anchor: "http"
weight: 21
---

Blackbeard also provide a web server and a websocket server exposing a REST api.

You can launch the Blackbeard server using the command :

```sh
blackbeard serve --help

Usage:
  blackbeard serve [flags]

Flags:
      --cors   Enable cors
  -h, --help   help for serve

Global Flags:
      --config string   config file (default is $HOME/.blackbeard.yaml)
      --dir string      Use the specified dir as root path to execute commands. Default is the current dir.
```

The REST api documentation is written following the [OpenAPI specifications](https://github.com/OAI/OpenAPI-Specification).

This documentation is available in an HTML format, using Swagger UI.

{{< oai-spec url="https://raw.githubusercontent.com/Meetic/blackbeard/master/swagger.json">}}
