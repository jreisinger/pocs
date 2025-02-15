package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	addr        = "localhost:44307"
	upstreamURL = "https://localhost:44300"
)

func main() {
	// Define the backend server URL
	backendURL, err := url.Parse(upstreamURL)
	if err != nil {
		log.Fatalf("Failed to parse backend URL: %v", err)
	}

	// Load the backend server's certificate
	caCert, err := os.ReadFile("localhost.pem")
	if err != nil {
		log.Fatalf("Failed to read backend server certificate: %v", err)
	}

	// Create a certificate pool with the backend server's certificate
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a custom transport with the certificate pool
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}

	// Create a reverse proxy
	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Update the request URL to the backend server
		r.URL.Scheme = backendURL.Scheme
		r.URL.Host = backendURL.Host
		// r.RequestURI = ""

		// Forward the request to the backend server
		resp, err := transport.RoundTrip(r)
		if err != nil {
			http.Error(w, "Failed to forward request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy the response headers and body to the client
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// Start the proxy server with HTTPS support
	log.Println("Listening on", addr)
	if err := http.ListenAndServeTLS(addr, "localhost.pem", "localhost-key.pem", proxy); err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}
