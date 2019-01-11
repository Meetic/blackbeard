package main

import (
	"github.com/Meetic/blackbeard/cmd"
)

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
