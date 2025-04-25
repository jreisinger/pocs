package main

import (
	"fmt"
	"net/http"

	"subrouters/db"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Creating a new router
	apiRouter := chi.NewRouter()
	apiRouter.Get("/articles/{date}-{slug}", getArticle)

	// Mounting the new sub router on the main router
	r.Mount("/api", apiRouter)

	http.ListenAndServe(":3000", r)
}
func getArticle(w http.ResponseWriter, r *http.Request) {
	dateParam := chi.URLParam(r, "date")
	slugParam := chi.URLParam(r, "slug")
	article, err := db.GetArticle(dateParam, slugParam)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(fmt.Appendf(nil, "fetching article %s-%s: %v", dateParam, slugParam, err))
		return
	}
	if article == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("article not found"))
		return
	}
	w.Write(article)
}
