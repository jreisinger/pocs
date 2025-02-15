package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	addr = "localhost:44300"
)

func main() {
	http.HandleFunc("/", helloHandler)
	log.Println("Listening on", addr)
	if err := http.ListenAndServeTLS(addr, "localhost.pem", "localhost-key.pem", nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
