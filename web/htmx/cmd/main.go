package main

import (
	"go-web-app/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", handlers.HomeHandler)
	http.HandleFunc("GET /lissajous", handlers.LissajousHandler)
	http.HandleFunc("GET /dynamic-content", handlers.DynamicContentHandler)

	log.Printf("starting to listen at :8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
