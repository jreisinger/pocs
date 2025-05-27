package handlers

import (
	"net/http"
	"strconv"
	"todo"
	"todo/components"

	"github.com/a-h/templ"
)

func Index(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.List(todo.Todos)).ServeHTTP(w, r)
}

func ToggleDone(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invald id", http.StatusBadRequest)
		return
	}

	for i := range todo.Todos {
		if todo.Todos[i].ID == id {
			todo.Todos[i].Done = !todo.Todos[i].Done
		}
	}
}
