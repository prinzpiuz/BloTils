// Package routes provides the HTTP routes for the server.
package routes

import (
	"BloTils/src/server"
	local_handlers "BloTils/src/server/handlers"
	"net/http"
	"os"
	"time"

	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
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
	lmt := init_rate_limiter()
	setRoutes(server, "/", local_handlers.IndexPage, http.MethodGet)
	setRoutes(server, "/clap_counter", local_handlers.ClapCounterPage, http.MethodGet)
	// API routes
	setRoutes(server, "/api/v1/ping", local_handlers.Ping, http.MethodGet)
	setRoutes(server, "/api/v1/count_like", tollbooth.LimitFuncHandler(lmt, local_handlers.GetClaps).ServeHTTP, http.MethodGet, http.MethodPost)
}

// init_rate_limiter creates a new rate limiter with a limit of 1 request per hour.
// It sets the IP lookup method to use the RemoteAddr field, and configures the
// error message and content type to be returned when the limit is reached.
func init_rate_limiter() *limiter.Limiter {
	var lmt = tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookup(limiter.IPLookup{
		Name:           "RemoteAddr",
		IndexFromRight: 0,
	})
	lmt.SetMessage("You have reached maximum request limit.")
	lmt.SetMessageContentType("application/json")
	return lmt
}

// ServeStaticFiles registers a file server handler on the provided server to serve
// static files from the configured static files directory. The path prefix "/static/"
// is used to match requests for static files.
func ServeStaticFiles(server *server.Server) {
	server.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(server.Config.StaticFiles))))
}
