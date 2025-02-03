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
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	UnAuthorized     = errors.New("UnAuthorized")
	BadRequest       = errors.New("Bad Request")
	NotFound         = errors.New("Not Found")
	Unsupported      = errors.New("Unsupported")
	Internal         = errors.New("Internal Server Error")
	Timeout          = errors.New("Timeout")
	BadGateway       = errors.New("Bad Gateway")
	Service          = errors.New("Service Unavailable")
	MethodNotAllowed = errors.New("Method Not Allowed")
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

func get_domain(url_string string) string {
	url, err := url.Parse(url_string)
	if err != nil {
		msg := fmt.Sprintf("Error parsing URL: %s", url_string)
		log.Print(msg)
		return ""
	}
	return url.Host
}

// check_for_domain checks if the given domain name is configured for the Blotils application.
// It retrieves the domain from the database and performs the following checks:
//  1. If the domain is not found, it returns an error with a log message.
//  2. If the domain is found and the context is a ClapCounter, it checks if likes are enabled for the domain.
//     If likes are not enabled, it returns an error with a log message.
//
// If all checks pass, it returns nil.
func check_for_domain(r *http.Request, domain_name string, context any) (db.Domain, error) {
	db_connection := get_db_connection(r)
	domain := db.GetDomain(db_connection, domain_name)
	if domain.IsEmpty() {
		msg := "This Domain(%s) Is Not Configured For Blotils"
		log.Printf(msg, domain_name)
		return db.Domain{}, fmt.Errorf(msg, domain_name)
	}
	_, ok := context.(ClapCounter)
	if ok {
		if !domain.LikesEnabled() {
			msg := "Like Counting Is Not Enabled For This Domain (%s)"
			log.Printf(msg, domain_name)
			return db.Domain{}, fmt.Errorf(msg, domain_name)
		}
	}

	return domain, nil
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

// path_to_cookie_str returns a string representation of a cookie name based on the provided path.
func path_to_cookie_str(cookie_name string, path string) string {
	return fmt.Sprintf("%s%s", cookie_name, strings.Join(strings.Split(path, "/"), "_"))
}

// setCookie sets an HTTP cookie with the provided name, value, and path. The cookie is set with
// HttpOnly, Secure, and SameSite=None attributes to ensure it is only accessible by the server
// and is transmitted securely over HTTPS.
func setCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   "127.0.0.1",
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, &cookie)

}

// getCookie retrieves the value of the cookie with the given name from the provided HTTP request.
// If the cookie is not found, it returns an error. If there is any other error retrieving the
// cookie, it wraps the error and returns it.
func getCookie(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			err = fmt.Errorf("Cookie Not Found")
		default:
			err = fmt.Errorf(err.Error())
		}
		return nil, err
	}
	return cookie, nil
}

// add_common_files appends the "favicon" template file to the given list of HTML template files.
// This function is used to ensure that the "favicon" template is always included when rendering
// HTML templates.
func add_common_files(files []string) []string {
	const html_location = "templates/%s.html"
	return append(files, fmt.Sprintf(html_location, "favicon"))
}

type TemplateData struct {
	Title    string
	MetaDesc string
	Static   string
	Data     interface{}
}

// set_template_data sets the default title and meta description for the template data if they are not already set.
// It also sets the static file path from the server configuration.
// The updated template data is returned.
func set_template_data(data TemplateData, r *http.Request) TemplateData {
	serverConfig, ok := r.Context().Value(server.ContextServerConfig).(*server.ServerConfig)
	if ok && serverConfig != nil {
		if data.Title == "" {
			data.Title = "BloTils - aka Blog uTils"
		}
		if data.MetaDesc == "" {
			data.MetaDesc = "A Blog Utils and Analytics Platform"
		}
		data.Static = serverConfig.StaticFiles
	}
	return data
}

// generateHTML renders the specified HTML templates with the provided data and writes the
// result to the given http.ResponseWriter.
//
// The filenames parameter specifies the names of the HTML template files to be rendered,
// without the ".html" extension. The templates are loaded from the "public/html/"
// directory.
//
// The data parameter provides the data to be used in rendering the templates.
//
// If there is an error parsing or executing the templates, the error is returned.
func generateHTML(w http.ResponseWriter, data TemplateData, filenames ...string) error {
	files := make([]string, 0, len(filenames))
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	files = add_common_files(files)
	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(w, "layout", data)
	return err
}
