package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte(os.Getenv("JWT_KEY"))

func main() {
	if len(secret) == 0 {
		log.Fatal("JWT_KEY environment variable must be set")
	}

	http.HandleFunc("POST /authenticate", authenticateHandler)
	http.HandleFunc("GET /verify", verifyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func keyFunc(token *jwt.Token) (any, error) {
	return secret, nil
}

func createAndSignToken(username string) (string, error) {
	issuedAt := time.Now()
	expirationTime := issuedAt.Add(time.Hour)

	claims := jwt.MapClaims{
		"iat":  issuedAt.Unix(),
		"exp":  expirationTime.Unix(),
		"name": username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Username and password are required\n")
		return
	}

	if err := verifyCredentials(username, password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "authenticating: %s\n", err)
		return
	}

	tokenString, err := createAndSignToken(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "TOKEN: %s\n", tokenString)
}

func verifyCredentials(username, password string) error {
	filename := "passwd"

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close() // Add this line

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return fmt.Errorf("%s must be lines of username:bcrypt-password-hash", filename)
		}
		fileUsername := parts[0]
		filePasswordHash := parts[1]
		if username == fileUsername {
			return bcrypt.CompareHashAndPassword([]byte(filePasswordHash), []byte(password))
		}
	}
	if err := s.Err(); err != nil {
		return err
	}

	return fmt.Errorf("username %s not found in %s", username, filename)
}

type CustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func verifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("unknown claims type")
	}

	return claims, nil
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	token := strings.TrimPrefix(header, "Bearer ")

	if header == "" || token == header {
		log.Println("missing authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := verifyToken(token)
	if err != nil {
		log.Println("invalid authorization token:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome back, %s!\n", claims.Name)
}
