package server

import (
	"BloTils/src/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
)

// ContextUpdateMiddleware is a higher-order middleware function that adds custom context data to the request context.
// It takes a context key and associated data, and returns a middleware function that injects this data
// into the request's context before passing it to the next handler.
//
// The returned middleware allows dynamic injection of context-specific information that can be
// retrieved by subsequent handlers using the provided context key.
func ContextUpdateMiddleware(contextDataKey ContextKey, contextData any) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextDataKey, contextData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// contentTypeSettingMiddleware sets the Content-Type header to
// application/json for all requests that have api in url otherwise text/html
func ContentTypeSettingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "api"):
			w.Header().Add("Content-Type", "application/json")
		case strings.Contains(r.URL.Path, "css"):
			w.Header().Add("Content-Type", "text/css")
		case strings.Contains(r.URL.Path, "js"):
			w.Header().Add("Content-Type", "text/javascript")
		case strings.Contains(r.URL.Path, ".svg"):
			w.Header().Add("Content-Type", "image/svg+xml")
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
func CORSPolicySettingMiddleware(next http.Handler) http.Handler {
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
func CSRFMiddleware(next http.Handler) http.Handler {
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

// AuthMiddleware is a middleware function that handles authentication for protected routes.
// It checks for a valid session cookie and ensures the user has appropriate access rights.
// If no valid session is found, the user is redirected to the login page.
// For admin-only routes, it prevents normal users from accessing those endpoints.
// The middleware adds the session information to the request context for downstream handlers.
func AuthMiddleware(next http.Handler) http.Handler {
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
		db_connection := GetDbConnection(r)
		if db_connection != nil {
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
