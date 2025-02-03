// Package handlers provides HTTP request handlers for the application.
package handlers

import (
	"BloTils/src/db"
	"encoding/json"
	"log"
	"net/http"
)

const handlerReference = "ClapCounter"

// Message constants for ClapCounter
const (
	ClapCount          = "Clap Count"
	ClapCountedFailed  = "Clap Counted Failed"
	ClapCountedSuccess = "Clap Counted Successfully"
	ClapAlreadyCounted = "Clap Already Counted"
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

// update_ip_like_count updates the like count for the IP address associated with the given ClapCounter.
// It retrieves the database connection from the HTTP request and calls the UpdateIPLikeCount function
// in the db package to update the like count for the IP address.
func update_ip_like_count(r *http.Request, clapCounter ClapCounter) {
	db_connection := get_db_connection(r)
	db.UpdateIPLikeCount(db_connection, clapCounter.URL, clapCounter.Page, clapCounter.remote_addr)
}

// get_likes retrieves the like count for the given URL and page from the database.
// It takes an http.Request and a ClapCounter as input, and returns a db.Likes struct
// containing the like count. If the like count is empty, it returns an empty db.Likes struct.
func get_likes(r *http.Request, clapCounter ClapCounter) db.Likes {
	db_connection := get_db_connection(r)
	likes := db.GetLikes(db_connection, clapCounter.URL, clapCounter.Page)
	if likes.IsEmpty() {
		return db.Likes{}
	}
	return likes
}

// already_liked_IP checks if the remote address associated with the given clapCounter has already
// liked the content identified by the clapCounter.URL and clapCounter.Page. It retrieves the list of
// liked IP addresses from the database and returns true if the remote address is found in the list,
// indicating the user has already liked the content.
func already_liked_IP(r *http.Request, clapCounter ClapCounter) bool {
	db_connection := get_db_connection(r)
	liked_ips := db.GetLikedIP(db_connection, clapCounter.URL, clapCounter.Page, clapCounter.remote_addr)
	if liked_ips.IsEmpty() {
		return false
	}
	return true
}

// add_like_to_page updates the like count for the given page and domain ID in the database.
// It takes an http.Request, a ClapCounter struct, and a domain ID as input, and calls the
// UpdateLikeCount function in the db package to update the like count for the given page.
// If the update is successful, it returns nil, otherwise it returns an error.
func add_like_to_page(r *http.Request, clapCounter ClapCounter, doamin_id int) error {
	db_connection := get_db_connection(r)
	return db.UpdateLikeCount(db_connection, clapCounter.Page, doamin_id)
}

// GetClaps is an HTTP handler function that handles requests to the clap counter endpoint.
// It processes the incoming request, updates the clap counter state, and writes the response
// as a JSON-encoded ClapCounter struct.
//
// The function first checks if the request is authorized by calling the check_for_domain function.
// If the request is not authorized, it sets the ClapCounter state accordingly and writes an error
// response.
//
// If the request is authorized, the function retrieves the current like count for the requested
// URL and page using the get_likes function. It then processes the request based on the HTTP method:
//
//   - GET requests: The function sets the ClapCounter state with the current like count and a success
//     message.
//   - POST requests: The function first checks if the request content type is valid, and then decodes
//     the request body into the ClapCounter struct. It then checks if the remote IP address has
//     already liked the content. If not, it updates the like count by calling the add_like_to_page
//     function, and then updates the IP like count by calling the update_ip_like_count function.
//     Finally, it sets the ClapCounter state with the updated like count and a success or failure
//     message.
//
// If any errors occur during the processing, the function logs the errors and writes an appropriate
// error response.
func GetClaps(w http.ResponseWriter, r *http.Request) {
	var clapCounter ClapCounter
	var continue_counting bool = true
	clapCounter.remote_addr = r.RemoteAddr
	clapCounter.URL = get_domain(r.Referer())
	clapCounter.Page = get_path(r, w)
	if clapCounter.Page == "" {
		err := decode_request(r, w, &clapCounter)
		if err != nil {
			log.Printf("Error decoding request: %v", err)
			return
		}
	}
	domain, domain_check_err := check_for_domain(r, clapCounter.URL, clapCounter)
	if domain_check_err != nil {
		continue_counting = false
		clapCounter.SetClapCounter(clapCounter.URL, domain_check_err.Error(), 0, false)
		setHTTPError(w, UnAuthorized, handlerReference, http.StatusForbidden)
	}
	if continue_counting {
		likes := get_likes(r, clapCounter)
		switch r.Method {
		case http.MethodGet:
			clapCounter.SetClapCounter(clapCounter.URL, ClapCount, likes.Count, true)
		case http.MethodPost:
			content_type_err := check_for_request_content_type(w, r)
			if content_type_err != nil {
				return
			}
			if already_liked_IP(r, clapCounter) {
				clapCounter.SetClapCounter(clapCounter.URL, ClapAlreadyCounted, likes.Count, false)
			} else {
				err := add_like_to_page(r, clapCounter, domain.ID)
				if err != nil {
					log.Printf("Error Updating Like Count: %v", err)
					clapCounter.SetClapCounter(clapCounter.URL, ClapCountedFailed, likes.Count, false)
					setHTTPError(w, Internal, handlerReference, http.StatusInternalServerError)
				}
				clapCounter.SetClapCounter(clapCounter.URL, ClapCountedSuccess, likes.Count+1, true)
			}
			update_ip_like_count(r, clapCounter)
		default:
			setHTTPError(w, MethodNotAllowed, handlerReference, http.StatusMethodNotAllowed)
		}
	}
	jsonData, err := json.Marshal(clapCounter)
	if err != nil {
		log.Printf("Error Encoding JSON: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// ClapCounterPage renders the HTML template for the clap counter page.
func ClapCounterPage(w http.ResponseWriter, r *http.Request) {
	templateData := set_template_data(TemplateData{}, r)
	err := generateHTML(w, templateData, "layout", "clap_counter")
	if err != nil {
		log.Printf("Error Generating HTML: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
