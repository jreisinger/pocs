package routes

import (
	"chi-demo/internal/handlers"
	"chi-demo/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) *chi.Mux {
	// Public route
	r.Get("/", handlers.HomeHandler)

	// Protected route group
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.Auth)
		protected.Get("/protected", handlers.ProtectedHandler)
	})

	return r
}
