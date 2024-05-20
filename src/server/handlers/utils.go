package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return errors.ErrUnsupported
		}
	}
	return nil
}

// decode_request is a utility function that decodes the request body into the provided destination interface.
// If there is an error decoding the request body, it will log the error and write an appropriate HTTP error response.
// The function returns the error, if any, from the decoding process.
func decode_request(w http.ResponseWriter, r *http.Request, dst interface{}) error {
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
			msg = fmt.Sprintf("Request body contains badly-formed JSON")

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
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return err
	}
	return nil

}
