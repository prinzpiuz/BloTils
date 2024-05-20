package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type GetClapRequest struct {
	URL string `json:"url"`
}

func GetClaps(w http.ResponseWriter, r *http.Request) {
	// check for domain
	// decode body
	// get count and return if request is GET
	// check if already counted for given ip,
	// if counted increase the ip count
	// if not counted add the ip to list and icrease count to one there
	// add url to and increase like count and return current count
	println(r.RemoteAddr)

	switch r.Method {
	case "GET":
		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			panic(err)
		}

		println(q.Get("url"))
	case "POST":
		content_type_err := check_for_request_content_type(w, r)
		if content_type_err != nil {
			return
		}
		var request_body GetClapRequest
		err := decode_request(w, r, &request_body)
		if err != nil {
			return
		}

		print(request_body.URL)
	default:
		msg := "Method Not Allowed"
		log.Println(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}

	p := map[string]string{"success": "true"}
	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Printf("Error Encoding JSON: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
