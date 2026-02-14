// Package handlers provides HTTP request handlers for the application.
// this page contains the utility functions for handlers
package server

import (
	"BloTils/src/db"
	"BloTils/src/models"
	"context"
	"database/sql"

	"crypto/rand"
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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnAuthorized     = errors.New("unauthorized")
	ErrBadRequest       = errors.New("bad request")
	ErrNotFound         = errors.New("not found")
	ErrUnsupported      = errors.New("unsupported")
	ErrInternal         = errors.New("internal server error")
	ErrTimeout          = errors.New("timeout")
	ErrBadGateway       = errors.New("bad gateway")
	ErrService          = errors.New("service unavailable")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrInvalidToken     = errors.New("invalid token")
)

const (
	SessionExpiry = 24 * time.Hour
	CSRFExpiry    = 5 * time.Minute
	SessionCookie = "sessionID"
	CSRFCookie    = "csrfToken"
)

// GetAppDataFromContext retrieves the application data from the given context.
// It returns the App data if found, or an empty App if no data is present.
// Panics if the context does not contain an App value with the AppContext key.
func GetAppDataFromContext(context context.Context) models.App {
	appData, ok := context.Value(AppContext).(models.App)
	if !ok {
		return models.App{}
	}
	return appData
}

// GetDBConfigFromContext retrieves the database configuration from the application context.
// It returns the DBConfig from the app data, or an empty DBConfig if no configuration is found.
func GetDBConfigFromContext(context context.Context) models.DBConfig {
	appData := GetAppDataFromContext(context)
	if appData.DBConfig.IsEmpty() {
		return models.DBConfig{}
	}
	return appData.DBConfig
}

// GetDbConnection retrieves the database connection from the HTTP request context.
// It returns the database connection if available, or nil if no connection is found.
func GetDbConnection(r *http.Request) *sql.DB {
	db := GetDBConfigFromContext(r.Context())
	if db.IsEmpty() {
		return nil
	}
	return db.Connection
}

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
//  2. If the domain is found and the context is a clapContext, it checks if likes are enabled for the domain.
//     If likes are not enabled, it returns an error with a log message.
//
// If all checks pass, it returns nil.
func checkForDomain(r *http.Request, domain_name string, context any) (db.Domain, error) {
	db_connection := GetDbConnection(r)
	domain := db.GetDomain(db_connection, domain_name)
	if domain.IsEmpty() {
		msg := "this domain(%s) is not configured for blotils"
		log.Printf(msg, domain_name)
		return db.Domain{}, fmt.Errorf(msg, domain_name)
	}
	_, ok := context.(clapContext)
	if ok {
		if !domain.LikesEnabled() {
			msg := "like counting is not enabled for this domain (%s)"
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
func setCookie(r *http.Request, w http.ResponseWriter, name string, value string, path string, httpOnly bool, expires time.Time, maxAge int, SameSitePolicy http.SameSite) {
	var secure = false
	appData := GetAppDataFromContext(r.Context())
	if !appData.AppConfig.IsEmpty() {
		secure = appData.AppConfig.IsProduction()
	}
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: SameSitePolicy,
		MaxAge:   maxAge,
	}
	if !expires.IsZero() {
		cookie.Expires = expires
	}
	if path != "" {
		cookie.Path = path
	}
	http.SetCookie(w, &cookie)

}

// addCommonFiles walks through the partials templates directory and appends all template files
// to the provided slice of files. It is used to dynamically include common template files
// during HTML template rendering.
func addCommonFiles(files []string) []string {
	commonTemplates := "templates/partials"
	err := filepath.Walk(commonTemplates, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			fileName := fmt.Sprintf("%s/%s", commonTemplates, info.Name())
			files = append(files, fileName)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", commonTemplates, err)
	}
	return files
}

type TemplateData struct {
	Title     string
	MetaDesc  string
	Static    string
	CSRFToken string
	Flash     *FlashMessage
	Session   *db.Session
	AppConfig models.AppConfig
	Data      map[string]any
}

// set_template_data sets the default title and meta description for the template data if they are not already set.
// It also sets the static file path from the server configuration.
// The updated template data is returned.
func setCommonTemplateData(r *http.Request, w http.ResponseWriter) TemplateData {
	appData := GetAppDataFromContext(r.Context())
	data := TemplateData{}
	if !appData.IsEmpty() {
		if data.Title == "" {
			data.Title = "BloTils - aka Blog uTils"
		}
		if data.MetaDesc == "" {
			data.MetaDesc = "A Blog Utils and Analytics Platform"
		}
		data.Static = appData.ServerConfig.StaticFiles
	}
	data.CSRFToken = generateCSRFToken(r, w)
	data.Session = GetSessionData(r)
	data.AppConfig = appData.AppConfig
	data.Flash = GetFlash(w, r)
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
	funcs := template.FuncMap{
		"notCurrentAdmin": notCurrentAdmin,
	}
	tmpl := template.New("").Funcs(funcs)
	templates := template.Must(tmpl.ParseFiles(files...))
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
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error Generating Secure Token %v", err)
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

// New CSRF token generation
func generateCSRFToken(r *http.Request, w http.ResponseWriter) string {
	csrfToken := generateSecureToken()
	setCookie(r, w, CSRFCookie, csrfToken, "", true, time.Now().Add(CSRFExpiry), 0, http.SameSiteStrictMode)
	return csrfToken
}

func generateSessionToken(r *http.Request, w http.ResponseWriter) string {
	sessionToken := generateSecureToken()
	setCookie(r, w, SessionCookie, sessionToken, "", true, time.Now().Add(SessionExpiry), 0, http.SameSiteStrictMode)
	return sessionToken
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

// GetSessionData retrieves the session data from the HTTP request context.
// If session data is found, it returns the session; otherwise, it returns an empty session.
func GetSessionData(r *http.Request) *db.Session {
	sessionData, ok := r.Context().Value(SessionContext).(*db.Session)
	if ok {
		return sessionData
	}
	return &db.Session{}
}

// getToggleValues retrieves the value of a form checkbox field and converts it to an integer.
// It returns 1 if the field value is "on", otherwise returns 0.
func getToggleValues(r *http.Request, fieldName string) int {
	toggleValue := r.FormValue(fieldName)
	if toggleValue == "on" {
		return 1
	}
	return 0

}

// commonIntParsingError handles errors that occur during integer parsing by logging the error
// and redirecting the user to the domains page if an error is encountered.
func commonIntParsingError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		log.Printf("Error parsing domain_id: %v", err)
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
	}
}
