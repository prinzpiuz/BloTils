// Package routes provides the HTTP routes for the server.
package routes

import (
	"BloTils/src/server"
	local_handlers "BloTils/src/server/handlers"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// setRoutes registers the given handler function for the specified path and HTTP
// methods on the provided server. It uses gorilla/mux for routing and
// gorilla/handlers for logging.
func setRoutes(server *server.Server, path string, handler http.HandlerFunc, methods ...string) {
	server.Router.Handle(path, handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(handler))).Methods(methods...)
}

// RegisterRoutes registers the API routes for the server.
// It calls setRoutes to add each route, specifying the path, handler function,
// and HTTP methods allowed.
func RegisterRoutes(server *server.Server) {
	setRoutes(server, "/api/v1/ping", local_handlers.Ping, http.MethodGet)
	setRoutes(server, "/api/v1/count_like", local_handlers.GetClaps, http.MethodGet, http.MethodPost)
}
