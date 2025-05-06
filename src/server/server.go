// Package server provides the main server implementation.
package server

import (
	"BloTils/src/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ContextKey string

var ProtectedRoutes = []string{}
var AdminOnly = []string{}
var DynamicRoutes = []string{}

const (
	ServerConfigContext ContextKey = "serverconfiguration"
	SessionContext      ContextKey = "sessiondata"
	InternalServerError string     = "Internal Server Error"
)

type Server struct {
	Router *mux.Router
	Config ServerConfig
}

type ServerConfig struct {
	Host        string
	Port        int
	StaticFiles string
}

func (server *Server) Start() {
	srv := &http.Server{
		Handler:      cors.Default().Handler(server.Router),
		Addr:         server.Config.Addr(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func InitServer(config ServerConfig) *Server {
	router := mux.NewRouter()
	// order matters
	router.Use(config.ContextUpdateMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(corsPolicySettingMiddleware)
	router.Use(csrfMiddleware)
	router.Use(authMiddleware)
	router.Use(contentTypeSettingMiddleware)
	server := &Server{
		Router: router,
		Config: config,
	}
	return server
}

type RouteDetails struct {
	Path          string
	Handler       http.HandlerFunc
	LoginRequired bool
	AdminRequired bool
	DynamicRoute  bool `default:"false"`
	Methods       []string
}

// SetRoute configures and registers a route with the server's router, handling route protection and logging.
// It supports setting routes as login-required (with dynamic or static paths), admin-only routes,
// and applies combined logging to the route handler. The route is registered with specified HTTP methods.
func (server *Server) SetRoute(routeDetails RouteDetails) {
	if routeDetails.LoginRequired {
		if routeDetails.DynamicRoute {
			firstPart := strings.Split(routeDetails.Path, "/")[1]
			DynamicRoutes = append(DynamicRoutes, firstPart)
		} else {
			ProtectedRoutes = append(ProtectedRoutes, routeDetails.Path)
		}
	}
	if routeDetails.AdminRequired {
		AdminOnly = append(AdminOnly, routeDetails.Path)
	}

	server.Router.Handle(routeDetails.Path, handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(routeDetails.Handler))).Methods(routeDetails.Methods...)
}

// contextUpdateMiddleware is a middleware that injects the ServerConfig into the request context.
func (c *ServerConfig) ContextUpdateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ServerConfigContext, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// contentTypeSettingMiddleware sets the Content-Type header to
// application/json for all requests that have api in url otherwise text/html
func contentTypeSettingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "api"):
			w.Header().Add("Content-Type", "application/json")
		case strings.Contains(r.URL.Path, "css"):
			w.Header().Add("Content-Type", "text/css")
		case strings.Contains(r.URL.Path, "js"):
			w.Header().Add("Content-Type", "text/javascript")
		case strings.Contains(r.URL.Path, "favicon"):
			w.Header().Add("Content-Type", "image/png")
		default:
			w.Header().Add("Content-Type", "text/html")
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

// corsPolicySettingMiddleware is a middleware function that sets the appropriate CORS headers
// for requests that contain "api" in the URL path. This allows cross-origin requests
// to the API endpoints.
//
// The middleware wraps the next http.Handler and sets the "Access-Control-Allow-Origin"
// and "Access-Control-Allow-Headers" headers before passing the request to the
// next handler.
func corsPolicySettingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "api") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		next.ServeHTTP(w, r)
	})
}

// csrfMiddleware is a middleware function that provides Cross-Site Request Forgery (CSRF) protection
// for non-API POST requests. It validates the CSRF token from the request cookie against
// the token submitted in the form, rejecting requests with missing or invalid tokens.
func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowedMethods = []string{http.MethodPost, http.MethodPut, http.MethodDelete}
		if slices.Contains(allowedMethods, r.Method) && !strings.Contains(r.URL.Path, "api") {
			var expectedToken string

			cookie, err := r.Cookie("csrfToken")
			if err != nil {
				msg := "Error: CSRF token missing"
				log.Println(msg)
				http.Error(w, msg, http.StatusForbidden)
				return
			}
			expectedToken = cookie.Value

			actualToken := r.FormValue("csrfToken")
			if actualToken != expectedToken {
				msg := "Error: CSRF token mismatch"
				log.Println(msg)
				http.Error(w, msg, http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// authMiddleware is a middleware function that handles session authentication for HTTP requests.
// It checks for a valid session cookie and redirects to the login page if no valid session exists.
// If a valid session is found, it adds the session to the request context for downstream handlers.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextpage := fmt.Sprintf("/login?next=%s", r.URL.Path)
		firstPart := strings.Split(r.URL.Path, "/")[1]
		if !slices.Contains(ProtectedRoutes, r.URL.Path) && !slices.Contains(DynamicRoutes, firstPart) {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			http.Redirect(w, r, nextpage, http.StatusFound)
			return
		}
		serverConfig, ok := r.Context().Value(ServerConfigContext).(*ServerConfig)
		if ok && serverConfig != nil {
			db_connection := serverConfig.DB.Connection
			session := db.GetSession(db_connection, cookie.Value)
			if session.IsEmpty() {
				http.Redirect(w, r, nextpage, http.StatusFound)
			}
			if session.IsExpired() {
				err := db.DeleteSessionWithSessionId(db_connection, session.SessionId)
				if err != nil {
					log.Println(err)
				}
				http.Redirect(w, r, nextpage, http.StatusFound)
			}
			if session.User.IsNormalUser() && slices.Contains(AdminOnly, r.URL.Path) {
				http.Redirect(w, r, nextpage, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), SessionContext, &session)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		http.Redirect(w, r, nextpage, http.StatusFound)
	})
}
