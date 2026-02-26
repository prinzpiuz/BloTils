// Package handlers provides HTTP request handlers for the application.
package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// Ping is an HTTP handler that responds with a "pong" message.
// It also checks the connection to the database and returns an error if the connection fails.
func Ping(w http.ResponseWriter, r *http.Request) {
	msg := "pong"
	db := GetDbConnection(r)
	if db != nil {
		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg = "Database Ping Error"
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		msg = "Database Not Found"
	}
	appData := GetAppDataFromContext(r.Context())
	responseBody := map[string]string{"message": msg, "version": appData.AppConfig.Version}
	jsonData, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error Encoding JSON: %v", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
