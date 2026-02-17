// Package server provides HTTP request handlers for the application.
package server

import (
	"BloTils/src/db"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
)

// Handler name constant for logging and error reporting
const clapCounterHandler = "ClapCounter"

// ClapCounter response messages
const (
	msgClapCount          = "Clap Count"
	msgClapCountedFailed  = "Clap Counted Failed"
	msgClapCountedSuccess = "Clap Counted Successfully"
	msgClapAlreadyCounted = "Clap Already Counted"
	msgDomainNotFound     = "Domain not configured"
	msgInvalidRequest     = "Invalid request"
)

// ClapResponse represents the JSON response for clap counter operations.
// Separating request and response structs follows the single responsibility principle.
type ClapResponse struct {
	URL     string `json:"url"`
	Message string `json:"message"`
	Count   int    `json:"count"`
	Success bool   `json:"success"`
	Page    string `json:"page"`
}

// ClapRequest represents the incoming request body for POST requests.
type ClapRequest struct {
	Page string `json:"page"`
}

// clapContext holds all the contextual data needed for processing a clap request.
// This reduces parameter passing and makes the code more maintainable.
type clapContext struct {
	url        string
	page       string
	remoteAddr string
	dbConn     *sql.DB
}

// newClapContext creates a new clapContext from an HTTP request.
// Returns an error if database connection is unavailable.
func newClapContext(r *http.Request, w http.ResponseWriter) (*clapContext, error) {
	dbConn := GetDbConnection(r)
	if dbConn == nil {
		return nil, errors.New("database connection unavailable")
	}

	ctx := &clapContext{
		remoteAddr: extractClientIP(r),
		url:        getDomain(r.Referer()),
		dbConn:     dbConn,
	}

	// Try to get page from query params first
	ctx.page = getPath(r, w)

	return ctx, nil
}

// extractClientIP extracts the real client IP, considering proxy headers.
// This is important for production as requests often come through load balancers.
func extractClientIP(r *http.Request) string {
	var ip string

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := findFirstComma(xff); idx != -1 {
			ip = xff[:idx]
		} else {
			ip = xff
		}
	} else if xri := r.Header.Get("X-Real-IP"); xri != "" {
		ip = xri
	} else {
		ip = r.RemoteAddr
	}
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return ip
	}
	return host
}

// findFirstComma returns the index of the first comma or -1 if not found.
// Avoiding strings.Index to minimize allocations.
func findFirstComma(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}

// GetClaps handles HTTP requests for the clap counter endpoint.
// Supports GET (retrieve count) and POST (add clap) methods.
func GetClaps(w http.ResponseWriter, r *http.Request) {
	// Validate content type for non-GET requests
	if r.Method != http.MethodGet {
		if err := checkForRequestContentType(w, r); err != nil {
			return
		}
	}

	ctx, err := newClapContext(r, w)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, msgInvalidRequest, "", "")
		return
	}

	// Parse request body for POST if page not in query params
	if ctx.page == "" && r.Method == http.MethodPost {
		var req ClapRequest
		if err := decodeRequest(r, w, &req); err != nil {
			log.Printf("[%s] Error decoding request: %v", clapCounterHandler, err)
			return
		}
		ctx.page = req.Page
	}

	// Validate page is present
	if ctx.page == "" {
		writeErrorResponse(w, http.StatusBadRequest, msgInvalidRequest, ctx.url, "")
		return
	}

	// Verify domain is configured
	domain, err := checkForDomain(r, ctx.url, clapContext{})
	if err != nil {
		log.Printf("[%s] Domain check failed for %s: %v", clapCounterHandler, ctx.url, err)
		writeErrorResponse(w, http.StatusForbidden, msgDomainNotFound, ctx.url, ctx.page)
		return
	}

	// Route to appropriate handler based on method
	switch r.Method {
	case http.MethodGet:
		handleGetClaps(w, ctx)
	case http.MethodPost:
		handlePostClap(w, ctx, domain.Id)
	default:
		setHTTPError(w, ErrMethodNotAllowed, clapCounterHandler, http.StatusMethodNotAllowed)
	}
}

// handleGetClaps processes GET requests to retrieve the current clap count.
func handleGetClaps(w http.ResponseWriter, ctx *clapContext) {
	likes := db.GetLikes(ctx.dbConn, ctx.url, ctx.page)

	writeJSONResponse(w, http.StatusOK, ClapResponse{
		URL:     ctx.url,
		Page:    ctx.page,
		Message: msgClapCount,
		Count:   likes.Count,
		Success: true,
	})
}

// handlePostClap processes POST requests to add a new clap.
func handlePostClap(w http.ResponseWriter, ctx *clapContext, domainID int) {
	// Check if already liked
	if hasAlreadyLikedIP(ctx) {
		likes := db.GetLikes(ctx.dbConn, ctx.url, ctx.page)
		writeJSONResponse(w, http.StatusOK, ClapResponse{
			URL:     ctx.url,
			Page:    ctx.page,
			Message: msgClapAlreadyCounted,
			Count:   likes.Count,
			Success: false,
		})
		return
	}

	// Add the like
	if err := db.UpdateLikeCount(ctx.dbConn, ctx.page, domainID); err != nil {
		log.Printf("[%s] Error updating like count: %v", clapCounterHandler, err)
		likes := db.GetLikes(ctx.dbConn, ctx.url, ctx.page)
		writeJSONResponse(w, http.StatusInternalServerError, ClapResponse{
			URL:     ctx.url,
			Page:    ctx.page,
			Message: msgClapCountedFailed,
			Count:   likes.Count,
			Success: false,
		})
		return
	}

	// Record the IP to prevent duplicate likes
	db.UpdateIPLikeCount(ctx.dbConn, ctx.url, ctx.page, ctx.remoteAddr)

	// Get updated count
	likes := db.GetLikes(ctx.dbConn, ctx.url, ctx.page)

	writeJSONResponse(w, http.StatusOK, ClapResponse{
		URL:     ctx.url,
		Page:    ctx.page,
		Message: msgClapCountedSuccess,
		Count:   likes.Count,
		Success: true,
	})
}

// hasAlreadyLikedIP checks if the IP has already liked this page.
func hasAlreadyLikedIP(ctx *clapContext) bool {
	likedIP := db.GetLikedIP(ctx.dbConn, ctx.url, ctx.page, ctx.remoteAddr)
	return !likedIP.IsEmpty()
}

// writeJSONResponse writes a JSON response with the given status code.
func writeJSONResponse(w http.ResponseWriter, statusCode int, response ClapResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[%s] Error encoding JSON response: %v", clapCounterHandler, err)
	}
}

// writeErrorResponse is a helper for writing error responses.
func writeErrorResponse(w http.ResponseWriter, statusCode int, message, url, page string) {
	writeJSONResponse(w, statusCode, ClapResponse{
		URL:     url,
		Page:    page,
		Message: message,
		Count:   0,
		Success: false,
	})
}

// ClapCounterPage renders the HTML template for the clap counter page.
func ClapCounterPage(w http.ResponseWriter, r *http.Request) {
	templateData := setCommonTemplateData(r, w)
	templateData.Data = map[string]interface{}{"samplePage": true}
	generateHTML(w, templateData, "layout", "clap_counter")
}
