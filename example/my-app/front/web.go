package main

import (
	"io"
	"net/http"

	"log"
)

const (
	apiURL = "http://blackbeard-example-app:50051"
)

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		panic(err)
	}

}

func main() {
	log.Print("front web running")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
