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
)

type ContextKey string

const ContextServerConfig ContextKey = "serverconfig"
const InternalServerError string = "Internal Server Error"

type Server struct {
	Router *mux.Router
	Config ServerConfig
}

type ServerConfig struct {
	Host string
	Port int
	DB   db.DB
}

func (server *Server) Start() {

	srv := &http.Server{
		Handler:      server.Router,
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
	router.Use(corsSettingMiddleware)
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
// application/json for all requests passed to the next handler.
func contentTypeSettingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func corsSettingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "api") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		next.ServeHTTP(w, r)
	})
}
