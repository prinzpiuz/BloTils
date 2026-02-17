package models

// Config holds Sentry configuration
type SentryConfig struct {
	Enabled bool
	DSN     string
	Debug   bool
}
