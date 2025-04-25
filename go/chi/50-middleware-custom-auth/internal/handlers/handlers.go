package handlers

import (
	"net/http"
)

// HomeHandler responds with a welcome message.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Chi demo application!"))
}

// ProtectedHandler responds with a message for authenticated users.
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have accessed a protected route!"))
}
