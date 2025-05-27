package middleware

import "net/http"

// Auth checks for a valid authentication token in the request
// headers. NOTE: hardcoded-credentials credentials in source code create the
// risk of unauthorized access!
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer my-secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
