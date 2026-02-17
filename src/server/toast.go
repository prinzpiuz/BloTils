package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Flash message types
const (
	FlashSuccess = "success"
	FlashError   = "error"
	FlashInfo    = "info"
	FlashWarning = "warning"
)

const flashCookieName = "blotils_flash"

// FlashMessage represents a toast notification
type FlashMessage struct {
	Type    string `json:"type"` // success, error, info, warning
	Message string `json:"message"`
}

// SetFlash stores a flash message in a cookie
func SetFlash(w http.ResponseWriter, r *http.Request, msgType, message string) {
	flash := FlashMessage{
		Type:    msgType,
		Message: message,
	}

	jsonData, err := json.Marshal(flash)
	if err != nil {
		log.Printf("Error marshaling flash message: %v", err)
		return
	}

	// Base64 encode to handle special characters
	encoded := base64.URLEncoding.EncodeToString(jsonData)
	setCookie(r, w, flashCookieName, encoded, "/", false, time.Time{}, 60, http.SameSiteLaxMode)
}

// GetFlash retrieves and clears the flash message
func GetFlash(w http.ResponseWriter, r *http.Request) *FlashMessage {
	cookie, err := r.Cookie(flashCookieName)
	if err != nil {
		return nil
	}

	// Clear the cookie immediately
	setCookie(r, w, flashCookieName, "", "/", false, time.Time{}, -1, http.SameSiteLaxMode)

	// Decode
	jsonData, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}

	var flash FlashMessage
	if err := json.Unmarshal(jsonData, &flash); err != nil {
		return nil
	}

	return &flash
}

// Helper functions for common use cases
func SetSuccessFlash(w http.ResponseWriter, r *http.Request, message string) {
	SetFlash(w, r, FlashSuccess, message)
}

func SetErrorFlash(w http.ResponseWriter, r *http.Request, message string) {
	SetFlash(w, r, FlashError, message)
}
