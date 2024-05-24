// Package handlers provides HTTP request handlers for the application.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Ping is an HTTP handler that responds with a "pong" message.
// It also checks the connection to the database and returns an error if the connection fails.
func Ping(w http.ResponseWriter, r *http.Request) {
	db := get_db_connection(r)
	if db != nil {
		err := db.Ping()
		if err != nil {
			setHTTPError(w, err, "Database Error", http.StatusInternalServerError)
		}

	} else {
		msg := "Database Not Found"
		setHTTPError(w, errors.New(msg), "Database Error", http.StatusInternalServerError)
	}
	p := map[string]string{"message": "pong"}
	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		setHTTPError(w, err, "JSON Error", http.StatusInternalServerError)
	}
}
