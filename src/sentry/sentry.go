package sentry

import (
	"BloTils/src/models"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

// Global flag to check if Sentry is initialized
var isInitialized bool

// Initialize sets up Sentry SDK if enabled
func Initialize(config *models.SentryConfig) {
	if !config.Enabled {
		log.Println("Sentry disabled")
		return
	}

	if config.DSN == "" {
		log.Println("Warning: Sentry enabled but DSN not configured. Sentry will be disabled.")
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.DSN,
		Debug:            config.Debug,
		TracesSampleRate: 0.2, // Adjust based on traffic (0.0 to 1.0)
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// You can modify or filter events here
			return event
		},
	})

	if err != nil {
		log.Printf("Warning: Sentry initialization failed: %v", err)
		return // Don't fail app startup
	}

	isInitialized = true
	log.Print("Sentry Initialized")
}

// CaptureException sends an error to Sentry
func CaptureException(err error) {
	if !isInitialized || err == nil {
		return
	}
	sentry.CaptureException(err)
}

// CaptureMessage sends a message to Sentry
func CaptureMessage(message string) {
	if !isInitialized {
		return
	}
	sentry.CaptureMessage(message)
}

// CaptureExceptionWithContext sends an error with additional context
func CaptureExceptionWithContext(err error, context map[string]interface{}) {
	if !isInitialized || err == nil {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		for key, value := range context {
			scope.SetExtra(key, value)
		}
		sentry.CaptureException(err)
	})
}

// SetUser sets user information for error tracking
func SetUser(id, email string) {
	if !isInitialized {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:    id,
			Email: email,
		})
	})
}

// ClearUser clears user information
func ClearUser() {
	if !isInitialized {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{})
	})
}

// AddBreadcrumb adds a breadcrumb for debugging
func AddBreadcrumb(category, message string) {
	if !isInitialized {
		return
	}
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: category,
		Message:  message,
		Level:    sentry.LevelInfo,
	})
}

// Flush waits for all events to be sent (call before app shutdown)
func Flush(timeout time.Duration) {
	if !isInitialized {
		return
	}
	sentry.Flush(timeout)
}

// Recover captures panics and sends to Sentry
// Usage: defer sentry.Recover()
func Recover() {
	if !isInitialized {
		return
	}
	if r := recover(); r != nil {
		sentry.CurrentHub().Recover(r)
		sentry.Flush(2 * time.Second)
		panic(r) // Re-panic after capturing
	}
}

// RecoverWithCallback captures panics and calls a callback
func RecoverWithCallback(callback func(error)) {
	if r := recover(); r != nil {
		var err error
		switch x := r.(type) {
		case error:
			err = x
		case string:
			err = fmt.Errorf("%s", x)
		default:
			err = fmt.Errorf("unknown panic: %v", x)
		}

		if isInitialized {
			sentry.CurrentHub().Recover(r)
			sentry.Flush(2 * time.Second)
		}

		if callback != nil {
			callback(err)
		}
	}
}

// IsEnabled returns whether Sentry is initialized and enabled
func IsEnabled() bool {
	return isInitialized
}
