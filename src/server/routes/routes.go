// Package routes provides the HTTP routes for the server.
package routes

import (
	"BloTils/src/server"
	localHandlers "BloTils/src/server/handlers"
	"net/http"
	"time"

	tollbooth "github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
)

// RegisterRoutes registers the API routes for the server.
// It calls setRoutes to add each route, specifying the path, handler function,
// and HTTP methods allowed.
func RegisterRoutes(s *server.Server) {
	lmt := init_rate_limiter()
	s.SetRoute(server.RouteDetails{
		Path:          "/",
		Handler:       localHandlers.IndexPage,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: false,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/clap_counter",
		Handler:       localHandlers.ClapCounterPage,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: false,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/create_account",
		Handler:       localHandlers.CreateAccountPage,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/forgot_password",
		Handler:       localHandlers.ForgotPassword,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/reset_password/{token}",
		Handler:       localHandlers.ResetPassword,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
		DynamicRoute:  true,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/login",
		Handler:       localHandlers.Login,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/logout",
		Handler:       localHandlers.Logout,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: true,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/admin",
		Handler:       localHandlers.AdminPage,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: true,
		LoginRequired: true,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/domains",
		Handler:       localHandlers.GetUserDomains,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: true,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/add_domain",
		Handler:       localHandlers.AddDomain,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: true,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/domain/{domain_id}",
		Handler:       localHandlers.EditDomainSettings,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: true,
		DynamicRoute:  true,
	})
	// API routes
	s.SetRoute(server.RouteDetails{
		Path:          "/api/v1/ping",
		Handler:       localHandlers.Ping,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: false,
	})
	s.SetRoute(server.RouteDetails{
		Path:          "/api/v1/count_like",
		Handler:       tollbooth.LimitFuncHandler(lmt, localHandlers.GetClaps).ServeHTTP,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
	})
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
