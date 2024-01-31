package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"tls-http/utils"
)

func main() {
	server := getServer()
	http.HandleFunc("/", myHandler)
	must(server.ListenAndServeTLS("", ""))
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request from %s", r.RemoteAddr)
	w.Write([]byte("Hello from server!\n"))
}

func getServer() *http.Server {
	tls := &tls.Config{
		GetCertificate: utils.CertReqFunc("localhost/cert.pem", "localhost/key.pem"),
	}
	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tls,
	}
	return server
}

// --- helper functions ---

func must(err error) {
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
