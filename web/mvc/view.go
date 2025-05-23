package main

import (
	"fmt"
	"net/http"
)

func RenderUser(w http.ResponseWriter, user User) {
	fmt.Fprintf(w, "<h1>User Profile</h1><p>ID: %d</p><p>Name: %s</p>", user.ID, user.Name)
}

func RenderNotFound(w http.ResponseWriter) {
	http.Error(w, "User not found", http.StatusNotFound)
}
