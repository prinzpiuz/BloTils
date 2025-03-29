// Package handlers provides HTTP request handlers for the application.
// this page contains the utility functions for handlers
package handlers

import (
	"BloTils/src/db"
	"BloTils/src/server"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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
	InvalidToken     = errors.New("Invalid Token")
)

const (
	SessionExpiry = 24 * time.Hour
	CSRFExpiry    = 5 * time.Minute
	SessionCookie = "sessionID"
	CSRFCookie    = "csrfToken"
)

// check_for_request_content_type checks the Content-Type header of the incoming HTTP request
// and returns an error if it is not "application/json". This ensures that the server only
// accepts JSON-encoded request bodies.
func checkForRequestContentType(w http.ResponseWriter, r *http.Request) error {
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
func decodeRequest(r *http.Request, w http.ResponseWriter, dst interface{}) error {
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
func GetDbConnection(r *http.Request) *sql.DB {
	serverConfig, ok := r.Context().Value(server.ServerConfigContext).(*server.ServerConfig)
	if ok && serverConfig != nil {
		return serverConfig.DB.Connection
	}
	return nil
}

func getDomain(url_string string) string {
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
func checkForDomain(r *http.Request, domain_name string, context any) (db.Domain, error) {
	db_connection := GetDbConnection(r)
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
func getPath(r *http.Request, w http.ResponseWriter) string {
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
// func path_to_cookie_str(cookie_name string, path string) string {
// 	return fmt.Sprintf("%s%s", cookie_name, strings.Join(strings.Split(path, "/"), "_"))
// }

// setCookie sets an HTTP cookie with the provided name, value, and path. The cookie is set with
// HttpOnly, Secure, and SameSite=None attributes to ensure it is only accessible by the server
// and is transmitted securely over HTTPS.
func setCookie(r *http.Request, w http.ResponseWriter, name string, value string, expires time.Time) {
	var secure bool = false
	serverConfig, ok := r.Context().Value(server.ServerConfigContext).(*server.ServerConfig)
	if ok && serverConfig != nil {
		secure = serverConfig.IsProduction
	}
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

}

// add_common_files appends the "favicon" template file to the given list of HTML template files.
// This function is used to ensure that the "favicon" template is always included when rendering
// HTML templates.
func addCommonFiles(files []string) []string {
	const htmlLocation = "templates/%s.html"
	commonFiles := []string{"favicon", "error_layout", "message_layout", "password", "sidebar", "footer", "topbar", "logo"}
	for _, file := range commonFiles {
		files = append(files, fmt.Sprintf(htmlLocation, file))
	}
	return files
}

type TemplateData struct {
	Title     string
	MetaDesc  string
	Static    string
	CSRFToken string
	Messages  []string
	Errors    []string
	Session   *db.Session
	Data      map[string]interface{}
}

// set_template_data sets the default title and meta description for the template data if they are not already set.
// It also sets the static file path from the server configuration.
// The updated template data is returned.
func setCommonTemplateTata(data TemplateData, r *http.Request, w http.ResponseWriter) TemplateData {
	serverConfig, ok := r.Context().Value(server.ServerConfigContext).(*server.ServerConfig)
	if ok && serverConfig != nil {
		if data.Title == "" {
			data.Title = "BloTils - aka Blog uTils"
		}
		if data.MetaDesc == "" {
			data.MetaDesc = "A Blog Utils and Analytics Platform"
		}
		data.Static = serverConfig.StaticFiles
	}
	data.CSRFToken = generateCSRFToken(r, w)
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
func generateHTML(w http.ResponseWriter, data TemplateData, filenames ...string) {
	files := make([]string, 0, len(filenames))
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	files = addCommonFiles(files)
	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Error Generating HTML: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// passwordHash generates a bcrypt hash for the given password using the default cost.
// It logs an error and panics if password hashing fails.
// Returns the hashed password as a string.
func passwordHash(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error Hashing Password %v", err)
		panic(err)
	}
	return string(hash)
}

// isValidEmail checks if the provided email address is valid by performing two validations:
// 1. Using mail.ParseAddress to check basic email format
// 2. Using a regular expression to validate email structure
// Returns true if the email is valid, false otherwise
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		log.Printf("Not A Valid Email %v", err)
		return false
	}
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// generateSecureToken creates a cryptographically secure random token
// by generating 32 random bytes and encoding them as a URL-safe base64 string.
// Returns a unique, random token suitable for use in security-sensitive contexts.
func generateSecureToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// New CSRF token generation
func generateCSRFToken(r *http.Request, w http.ResponseWriter) string {
	csrfToken := generateSecureToken()
	setCookie(r, w, CSRFCookie, csrfToken, time.Now().Add(CSRFExpiry))
	return csrfToken
}

func generateSessionToken(r *http.Request, w http.ResponseWriter) string {
	sessionToken := generateSecureToken()
	setCookie(r, w, SessionCookie, sessionToken, time.Now().Add(SessionExpiry))
	return sessionToken
}

// sendMessagePage renders a message page with a single message using the provided HTTP request and response writer.
// It sets common template data, adds the specified message, and generates an HTML page using the "layout" and "message_page" templates.
func sendMessagePage(r *http.Request, w http.ResponseWriter, msg string) {
	templateData := setCommonTemplateTata(TemplateData{}, r, w)
	templateData.Messages = []string{msg}
	generateHTML(w, templateData, "layout", "message_page")
}

// getUrlVars retrieves a specific URL variable from an HTTP request using Gorilla Mux.
// It takes an HTTP request and the name of the variable to extract.
// Returns the value of the specified URL variable as a string.
func getUrlVars(r *http.Request, variable string) string {
	vars := mux.Vars(r)
	return vars[variable]
}

// parseForm attempts to parse the HTTP request form data and handles any parsing errors
// by logging the error and sending an HTTP 500 Internal Server Error response.
func parseForm(r *http.Request, w http.ResponseWriter) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error Parsing Form: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// validLogin checks if a user's login credentials are valid by verifying:
// 1. The user exists
// 2. The user account is active
// 3. The user status is approved
// 4. The provided password matches the stored password hash
// Returns true if all conditions are met, false otherwise.
func validLogin(user db.User, password string) bool {
	passwordCheckErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return user.UserExist() && user.IsActive && user.UserStatusApproved() && passwordCheckErr == nil
}

func GetSessionData(r *http.Request) *db.Session {
	sessionData, ok := r.Context().Value(server.SessionContext).(*db.Session)
	if ok {
		return sessionData
	}
	return &db.Session{}
}
