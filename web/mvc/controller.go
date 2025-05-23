package main

import (
	"log"
	"net/http"
	"strconv"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, found := GetUserByID(id)
	if !found {
		RenderNotFound(w)
		return
	}
	RenderUser(w, user)
}

func main() {
	http.HandleFunc("/user", userHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
