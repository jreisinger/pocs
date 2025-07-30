package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret is used to sign and verify JWT tokens.
var jwtSecret []byte

func init() {
	// Generate random 256 bits (32 bytes) for jwtSecret
	jwtSecret = make([]byte, 32)
	if _, err := rand.Read(jwtSecret); err != nil {
		log.Fatalf("failed to generate random jwtSecret: %v", err)
	}
}

func main() {
	r := chi.NewRouter()

	r.Post("/login", loginHandler)
	r.With(jwtMiddleware).Get("/protected", protectedHandler)

	log.Println("server started on localhost:3000")
	log.Fatal(http.ListenAndServe("localhost:3000", r))
}

// loginHandler is a dummy implementation that returns a JWT token
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Normally you'd validate username and password here
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenStr))
}

// jwtMiddleware protects routes via authn
func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		t, err := token.Claims.GetExpirationTime()
		if err != nil {
			log.Printf("getting token expiration time: %v", err)
		} else {
			log.Printf("token is valid till %s", t)
		}

		// Add user context here if needed (e.g. claims["sub"])

		next.ServeHTTP(w, r)
	})
}

// protectedHandler handles the route protected by authn
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You accessed a protected route! 🎉"))
}
