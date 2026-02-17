package models

// Config holds SMTP configuration
type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
	BaseURL   string
	Enabled   bool
}
