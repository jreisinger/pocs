package main

import (
	"chi-demo/internal/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Setup routes
	routes.SetupRoutes(r)

	// Start the server
	http.ListenAndServe(":8080", r)
}
