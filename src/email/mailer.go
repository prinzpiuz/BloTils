package mailer

import (
	"BloTils/src/models"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Email represents an email to be sent
type Email struct {
	To          string
	Subject     string
	HTMLContent string
	TextContent string
}

// Mailer handles sending emails via SMTP
type Mailer struct {
	config *models.SMTPConfig
}

// Global mailer instance
var defaultMailer *Mailer

// Initialize sets up the global mailer
func Initialize(config *models.SMTPConfig) {
	defaultMailer = &Mailer{config: config}

	if !config.Enabled {
		log.Println("Email service disabled")
		return
	}

	if config.Host == "" {
		log.Println("Warning: SMTP not configured. Emails will be logged only.")
		return
	}

	log.Printf("Email service initialized: %s:%d", config.Host, config.Port)
}

// Send sends an email using SMTP
func (m *Mailer) Send(email Email) error {
	if m == nil || !m.config.Enabled {
		log.Printf("[EMAIL-DISABLED] To: %s, Subject: %s", email.To, email.Subject)
		return nil
	}

	if m.config.Host == "" {
		log.Printf("[EMAIL-LOG] To: %s, Subject: %s", email.To, email.Subject)
		log.Printf("[EMAIL-LOG] Content: %s", truncate(email.HTMLContent, 200))
		return nil
	}

	return m.sendSMTP(email)
}

func (m *Mailer) sendSMTP(email Email) error {
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	// Build from address
	from := m.config.FromEmail
	if m.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", m.config.FromName, m.config.FromEmail)
	}

	// Build message with proper headers
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", email.To))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(email.HTMLContent)

	// Setup authentication
	var auth smtp.Auth
	if m.config.Username != "" && m.config.Password != "" {
		auth = smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	}

	// Send with TLS
	err := m.sendWithTLS(addr, auth, m.config.FromEmail, email.To, msg.String())
	if err != nil {
		return fmt.Errorf("SMTP send failed: %w", err)
	}

	return nil
}

// sendWithTLS sends email using STARTTLS (port 587)
func (m *Mailer) sendWithTLS(addr string, auth smtp.Auth, from, to, msg string) error {
	// Connect to server
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing SMTP connection: %v", err)
		}
	}()

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: m.config.Host,
		MinVersion: tls.VersionTLS12,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	// Authenticate if credentials provided
	if auth != nil {
		if err = conn.Auth(auth); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
	}

	// Send email
	if err = conn.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	if err = conn.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := conn.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close failed: %w", err)
	}

	return conn.Quit()
}

// GetBaseURL returns the configured base URL
func (m *Mailer) GetBaseURL() string {
	if m.config.BaseURL != "" {
		return m.config.BaseURL
	}
	return "http://localhost:8000"
}

// Helper to truncate strings for logging
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Send sends an email using the default mailer
func Send(email Email) error {
	if defaultMailer == nil {
		log.Println("Warning: Mailer not initialized")
		return nil
	}
	return defaultMailer.Send(email)
}

// GetBaseURL returns the base URL from default mailer
func GetBaseURL() string {
	if defaultMailer == nil {
		return "http://localhost:8000"
	}
	return defaultMailer.GetBaseURL()
}
