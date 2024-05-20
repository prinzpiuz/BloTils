package handlers

import (
	"BloTils/src/server"
	"encoding/json"
	"log"
	"net/http"
)

// Ping is an HTTP handler that responds with a "pong" message.
// It also checks the connection to the database and returns an error if the connection fails.
func Ping(w http.ResponseWriter, r *http.Request) {
	serverConfig, ok := r.Context().Value(server.ContextServerConfig).(*server.ServerConfig)
	if ok && serverConfig != nil {
		err := serverConfig.DB.Connection.Ping()
		if err != nil {
			log.Printf("Database Error %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	p := map[string]string{"message": "pong"}
	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Printf("Error Encoding JSON: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
