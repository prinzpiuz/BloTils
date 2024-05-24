// Package handlers provides HTTP request handlers for the application.
// this page contains the utility functions for handlers
package handlers

import (
	"BloTils/src/db"
	"BloTils/src/server"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// check_for_request_content_type checks the Content-Type header of the incoming HTTP request
// and returns an error if it is not "application/json". This ensures that the server only
// accepts JSON-encoded request bodies.
func check_for_request_content_type(w http.ResponseWriter, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			setHTTPError(w, errors.New(msg), "Content-Type", http.StatusUnsupportedMediaType)
			return errors.ErrUnsupported
		}
	}
	return nil
}

// decode_request is a utility function that decodes the request body into the provided destination interface.
// If there is an error decoding the request body, it will log the error and write an appropriate HTTP error response.
// The function returns the error, if any, from the decoding process.
func decode_request(r *http.Request, w http.ResponseWriter, dst interface{}) error {
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&dst)
	var msg string
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg = "Request body contains badly-formed JSON"

		case errors.As(err, &unmarshalTypeError):
			msg = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)

		case errors.Is(err, io.EOF):
			msg = "Request body must not be empty"

		case err.Error() == "http: request body too large":
			msg = "Request body must not be larger than 1MB"

		default:
			msg = err.Error()
		}
		setHTTPError(w, err, msg, http.StatusBadRequest)
		return err
	}
	return nil

}

// get_db_connection returns the database connection stored in the server configuration
// context. If the configuration is not available, it returns nil.
func get_db_connection(r *http.Request) *sql.DB {
	serverConfig, ok := r.Context().Value(server.ContextServerConfig).(*server.ServerConfig)
	if ok && serverConfig != nil {
		return serverConfig.DB.Connection
	}
	return nil
}

// check_for_domain checks if the given domain name is configured for the Blotils application.
// It retrieves the domain from the database and performs the following checks:
//  1. If the domain is not found, it returns an error with a log message.
//  2. If the domain is found and the context is a ClapCounter, it checks if likes are enabled for the domain.
//     If likes are not enabled, it returns an error with a log message.
//
// If all checks pass, it returns nil.
func check_for_domain(r *http.Request, domain_name string, context any) error {
	db_connection := get_db_connection(r)
	domain := db.GetDomain(db_connection, domain_name)
	if domain.IsEmpty() {
		msg := "This Domain(%s) Is Not Configured For Blotils"
		log.Printf(msg, domain_name)
		return fmt.Errorf(msg, domain_name)
	}
	_, ok := context.(ClapCounter)
	if ok {
		if !domain.LikesEnabled() {
			msg := "Like Counting Is Not Enabled For This Domain (%s)"
			log.Printf(msg, domain_name)
			return fmt.Errorf(msg, domain_name)
		}
	}

	return nil
}

// setHTTPError writes an HTTP error response with the provided error message and HTTP status code.
// The reference parameter is used to provide additional context about the error.
func setHTTPError(w http.ResponseWriter, err error, reference string, httpStatus int) {
	msg := fmt.Sprintf("%s:  %v", reference, err.Error())
	log.Println(msg)
	http.Error(w, msg, httpStatus)
}

// get_path extracts the path from the query string of the given HTTP request.
// If there is an error parsing the query string, it sets an HTTP error response
// and returns an empty string.
func get_path(r *http.Request, w http.ResponseWriter) string {

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		setHTTPError(w, err, "Error Parsing Query", http.StatusBadRequest)
		return ""
	}
	s, err := url.Parse(q.Get("page"))
	if err != nil {
		setHTTPError(w, err, "Error Parsing Query", http.StatusBadRequest)
		return ""
	}
	return s.Path
}

// func get_like_count(r *http.Request, clapCounter ClapCounter) db.Likes {
// 	db_connection := get_db_connection(r)
// 	likes := db.GetLikes(db_connection, clapCounter.URL, clapCounter.Page)
// 	if likes.IsEmpty() {
// 		return db.Likes{count=0}
// 	}
// 	return likes
// }
