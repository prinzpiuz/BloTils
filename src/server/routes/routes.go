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
	server.Router.Handle(path, handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler))).Methods(methods...)
}

// RegisterRoutes registers the API routes for the server.
// It calls setRoutes to add each route, specifying the path, handler function,
// and HTTP methods allowed.
func RegisterRoutes(server *server.Server) {
	setRoutes(server, "/v1/ping", local_handlers.Ping, "GET")
	setRoutes(server, "/v1/count/{url}", local_handlers.GetClaps, "GET")
}
