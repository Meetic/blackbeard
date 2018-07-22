package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	version = 2
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello API version %d\n", version)
}

func main() {
	log.Printf("API version %d is running", version)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":50051", nil)
}
