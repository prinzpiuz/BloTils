// Package routes provides the HTTP routes for the server.
package server

import (
	"BloTils/src/models"
	"net/http"
	"os"
	"strings"
	"time"

	tollbooth "github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
	"github.com/gorilla/handlers"
)

// SetRoute configures and registers a route with the server's router, handling route protection and logging.
// It supports setting routes as login-required (with dynamic or static paths), admin-only routes,
// and applies combined logging to the route handler. The route is registered with specified HTTP methods.
func setRoute(server *models.Server, routeDetails RouteDetails) {
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

// RegisterRoutes registers the API routes for the server.
// It calls setRoutes to add each route, specifying the path, handler function,
// and HTTP methods allowed.
func RegisterRoutes(s *models.Server) {
	lmt := init_rate_limiter()
	setRoute(s, RouteDetails{
		Path:          "/",
		Handler:       IndexPage,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: false,
	})
	setRoute(s, RouteDetails{
		Path:          "/clap_counter",
		Handler:       ClapCounterPage,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: false,
	})
	setRoute(s, RouteDetails{
		Path:          "/create_account",
		Handler:       CreateAccountPage,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
	})
	setRoute(s, RouteDetails{
		Path:          "/forgot_password",
		Handler:       ForgotPassword,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
	})
	setRoute(s, RouteDetails{
		Path:          "/reset_password/{token}",
		Handler:       ResetPassword,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
		DynamicRoute:  true,
	})
	setRoute(s, RouteDetails{
		Path:          "/login",
		Handler:       Login,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: false,
	})
	setRoute(s, RouteDetails{
		Path:          "/logout",
		Handler:       Logout,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: true,
	})
	setRoute(s, RouteDetails{
		Path:          "/admin",
		Handler:       AdminPage,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: true,
		LoginRequired: true,
	})
	setRoute(s, RouteDetails{
		Path:          "/approve_account/{user_id}",
		Handler:       ApproveUserRequest,
		Methods:       []string{http.MethodGet},
		AdminRequired: true,
		LoginRequired: true,
		DynamicRoute:  true,
	})
	setRoute(s, RouteDetails{
		Path:          "/delete_user/{user_id}",
		Handler:       DeleteUser,
		Methods:       []string{http.MethodGet},
		AdminRequired: true,
		LoginRequired: true,
		DynamicRoute:  true,
	})
	setRoute(s, RouteDetails{
		Path:          "/domains",
		Handler:       GetUserDomains,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: true,
	})
	setRoute(s, RouteDetails{
		Path:          "/add_domain",
		Handler:       AddDomain,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: true,
	})
	setRoute(s, RouteDetails{
		Path:          "/domain/{domain_id}",
		Handler:       EditDomainSettings,
		Methods:       []string{http.MethodGet, http.MethodPost},
		AdminRequired: false,
		LoginRequired: true,
		DynamicRoute:  true,
	})
	setRoute(s, RouteDetails{
		Path:          "/delete_domain/{domain_id}",
		Handler:       DeleteDomain,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: true,
		DynamicRoute:  true,
	})
	// API routes
	setRoute(s, RouteDetails{
		Path:          "/api/v1/ping",
		Handler:       Ping,
		Methods:       []string{http.MethodGet},
		AdminRequired: false,
		LoginRequired: false,
	})
	setRoute(s, RouteDetails{
		Path:          "/api/v1/count_like",
		Handler:       tollbooth.LimitFuncHandler(lmt, GetClaps).ServeHTTP,
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
func ServeStaticFiles(server *models.Server) {
	server.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(server.Config.StaticFiles))))
}
