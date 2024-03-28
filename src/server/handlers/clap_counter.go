package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetClaps(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	url := params["url"]
	print(url)
	p := map[string]string{"sucesss": "true"}
	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Printf("Error Encoding JSON: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
