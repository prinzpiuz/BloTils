// Package server provides the main server implementation.
package server

import (
	"BloTils/src/models"
	"fmt"
	"log"
	"net/http"
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
	}
	log.Fatal(srv.ListenAndServe())
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
	router.Use(ContextUpdateMiddleware(AppContext, app))
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(CORSPolicySettingMiddleware)
	router.Use(CSRFMiddleware)
	router.Use(AuthMiddleware)
	router.Use(ContentTypeSettingMiddleware)

}
