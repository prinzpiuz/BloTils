// Package server provides the main server implementation.
package server

import (
	"BloTils/src/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ContextKey string

const ContextServerConfig ContextKey = "serverconfig"
const InternalServerError string = "Internal Server Error"

type Server struct {
	Router *mux.Router
	Config ServerConfig
}

type ServerConfig struct {
	Host        string
	Port        int
	DB          db.DB
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

func New(config ServerConfig) *Server {
	router := mux.NewRouter()
	err := config.DB.Initialize()
	if err != nil {
		log.Fatalf("Error Initializing DB: %v", err)
	}
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(corsPolicySettingMiddleware)
	router.Use(contentTypeSettingMiddleware)
	router.Use(config.contextUpdateMiddleware)
	server := &Server{
		Router: router,
		Config: config,
	}
	return server
}

// contextUpdateMiddleware is a middleware that injects the ServerConfig into the request context.
func (c *ServerConfig) contextUpdateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextServerConfig, c)
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
