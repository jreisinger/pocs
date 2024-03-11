package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port = "8080"
	msg  = `Do not meddle in the affairs of wizards,
for you are crunchy and good with ketchup.`
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, msg)
	})
	log.Println("starting a web server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
