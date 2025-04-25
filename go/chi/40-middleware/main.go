package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(MyMiddlerware)
	r.Get("/", MyHandler)
	http.ListenAndServe(":3000", r)
}

// HTTP middleware setting a value on the request context
func MyMiddlerware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from request context and assign key user to value 123
		ctx := context.WithValue(r.Context(), "user", "123")

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new "user" key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	w.Write([]byte(fmt.Sprintf("hi %s", user)))
}
