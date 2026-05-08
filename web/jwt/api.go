package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	router := gin.Default()
	router.GET("/protected", protectedMessage)

	log.Fatal(router.Run("localhost:8080"))
}

type responseJSON struct {
	Message string `json:"message"`
}

var (
	// I took this values from the sample token at https://jwt.io/
	Key             = []byte("your-256-bit-secret")
	ValidAlgMethods = []string{"HS256"}
)

func protectedMessage(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")

	token, err := jwt.Parse(
		getToken(authorizationHeader),
		func(token *jwt.Token) (any, error) {
			return Key, nil
		},
		jwt.WithValidMethods(ValidAlgMethods),
		// jwt.WithExpirationRequired(), // should be used in prod
	)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responseJSON{Message: err.Error()})
		return
	}
	if !token.Valid {
		c.IndentedJSON(http.StatusUnauthorized, responseJSON{Message: "token not valid"})
		return
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responseJSON{Message: err.Error()})
		return
	}
	if sub != "philosopher" {
		c.IndentedJSON(http.StatusUnauthorized, responseJSON{Message: "sub claim in payload is not philosopher"})
		return
	}

	message := "the man's ultimate happiness consists in the contemplation of truth"
	c.IndentedJSON(http.StatusOK, responseJSON{Message: message})
}

func getToken(authorizationHeader string) string {
	fields := strings.Fields(authorizationHeader)

	// Bearer <token>
	if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
		return fields[1]
	}

	// <token>
	if len(fields) == 1 {
		return fields[0]
	}

	return ""
}
