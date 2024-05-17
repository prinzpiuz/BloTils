package handlers

import (
	"BloTils/src/server"
	"encoding/json"
	"log"
	"net/http"
)

// Ping handles ping requests and responds with a JSON "pong" message.
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
