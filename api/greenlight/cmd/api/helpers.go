package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// NOTE: readIDParam doesn't use any deps from our application struct so it
// could just be a regular function.  But in general it's good to setup all
// app-specific handler and helpers as methods on application:
//   - it helps maintain consistency in the code structure
//   - it future-proofs code for when handler and helpers in the future *do* need
//     access to the struct
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// When httprouter is parsing a request, any interpolated URL parameters
	// will be stored in the request context.
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id paramater")
	}

	return id, nil
}

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	js = append(js, '\n') // small nicety for terminal users

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
