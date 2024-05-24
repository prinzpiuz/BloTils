// Package handlers provides HTTP request handlers for the application.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// ClapCounter represents a counter for tracking the number of "claps" or interactions
// with a specific URL, message, and page.
type ClapCounter struct {
	URL         string `json:"url"`
	Message     string `json:"message"`
	Count       int    `json:"count"`
	Success     bool   `json:"success"`
	Page        string `json:"page"`
	remote_addr string
}

// SetClapCounter sets the URL, Message, Count, and Success fields of the ClapCounter struct.
// This function is used to update the state of the ClapCounter.
func (clapCounter *ClapCounter) SetClapCounter(url string, message string, count int, success bool) {
	clapCounter.URL = url
	clapCounter.Message = message
	clapCounter.Count = count
	clapCounter.Success = success
}

func GetClaps(w http.ResponseWriter, r *http.Request) {
	// check for domain
	// decode body
	// get count and return if request is GET
	// check if already counted for given ip,
	// if counted increase the ip count
	// if not counted add the ip to list and icrease count to one there
	// add url to and increase like count and return current count

	var clapCounter ClapCounter
	clapCounter.remote_addr = r.RemoteAddr
	clapCounter.URL = r.Referer()
	switch r.Method {
	case "GET":
		clapCounter.Page = get_path(r, w)
		domain_check_err := check_for_domain(r, clapCounter.URL, clapCounter)
		if domain_check_err != nil {
			clapCounter.SetClapCounter(clapCounter.URL, domain_check_err.Error(), 0, false)
		}
		// clapCounter.count = get_like_count()
		clapCounter.SetClapCounter(clapCounter.URL, "Clap Count", 0, true)
	case "POST":
		content_type_err := check_for_request_content_type(w, r)
		if content_type_err != nil {
			return
		}
		err := decode_request(r, w, &clapCounter)
		if err != nil {
			return
		}
	default:
		msg := "Method Not Allowed"
		log.Println(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}
	jsonData, err := json.Marshal(clapCounter)
	if err != nil {
		log.Printf("Error Encoding JSON: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
