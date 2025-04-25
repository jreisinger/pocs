package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.With(httprate.LimitByIP(5, 1*time.Second)).Get("/limit", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok\n"))
		})

		r.Get("/nolimit", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok\n"))
		})
	})

	http.ListenAndServe(":3000", r)
}
