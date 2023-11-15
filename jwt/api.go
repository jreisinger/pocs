package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	router := gin.Default()
	router.GET("/protected", protectedMessage)

	router.Run("localhost:8080")
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
		func(token *jwt.Token) (interface{}, error) {
			return Key, nil
		},
		jwt.WithValidMethods(ValidAlgMethods),
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
	if sub != "Aquinas" {
		c.IndentedJSON(http.StatusUnauthorized, responseJSON{Message: "sub claim in payload is not Aquinas"})
		return
	}

	message := "the man's ultimate happiness consists in the contemplation of truth"
	c.IndentedJSON(http.StatusOK, responseJSON{Message: message})
}

func getToken(authorizationHeader string) string {
	if authorizationHeader == "" {
		return ""
	}

	var token string
	fields := strings.Fields(authorizationHeader)
	if len(fields) == 1 {
		token = fields[0] // <token>
	} else {
		token = fields[1] // Bearer <token> [...]
	}

	return token
}
