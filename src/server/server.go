// Package server provides the main server implementation.
package server

import (
	"BloTils/src/db"
	"BloTils/src/models"
	"BloTils/src/sentry"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ContextKey string

var ProtectedRoutes = []string{}
var AdminOnly = []string{}
var DynamicRoutes = []string{}

const (
	ServerConfigContext ContextKey = "serverconfiguration"
	AppContext          ContextKey = "appcontext"
	SessionContext      ContextKey = "sessiondata"
	InternalServerError string     = "Internal Server Error"
)

func Start(server *models.Server) {
	srv := &http.Server{
		Handler:      cors.Default().Handler(server.Router),
		Addr:         addr(server.Config),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	sig := <-quit
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	cleanup()

	log.Println("Server stopped")
}

func cleanup() {
	log.Println("Running cleanup tasks...")
	sentry.Flush(5 * time.Second)
	db.CloseDB()
	log.Println("Cleanup complete")
}

func addr(c *models.ServerConfig) string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func InitServer(config *models.ServerConfig) *models.Server {
	router := mux.NewRouter()
	server := &models.Server{
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

func SetupMiddlewares(router *mux.Router, app models.App) {
	router.Use(SentryMiddleware())
	router.Use(ContextUpdateMiddleware(AppContext, app))
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(CORSPolicySettingMiddleware)
	router.Use(CSRFMiddleware)
	router.Use(AuthMiddleware)
	router.Use(ContentTypeSettingMiddleware)

}
