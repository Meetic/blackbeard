package main

import (
	"github.com/Meetic/blackbeard/cmd"
)

// @title Blackbeard API
// @version 1.0
// @description This is a REST API to manage Blackbeard. See https://github.com/Meetic/blackbeard

// @contact.name Sébastien Le gall
// @contact.url http://le-gall.bzh
// @contact.email seb@le-gall.net

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	cmd.Execute()
}
