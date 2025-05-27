package main

import (
	"log"
	"net/http"
	"todo/handlers"
)

func main() {
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/toggle-done", handlers.ToggleDone)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
