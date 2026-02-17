package models

// Config holds SMTP configuration
type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
	BaseURL   string // For generating links in emails
	Enabled   bool   // Easy toggle for development
}
