package main

import (
	"log"

	"github.com/Meetic/blackbeard/cmd"
)

func main() {
	if err := cmd.NewBlackbeardCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
