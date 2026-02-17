package app

import (
	"BloTils/src/models"
	"fmt"
	"log"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/knadh/koanf/v2"
)

var logoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))

const appLogo = `
██████  ██    ██████ ██████ ██ ██    ███████
██  ██  ██    █    █   ██   ██ ██    ██
██████  ██    █    █   ██   ██ ██    ███████
██  ██  ██    █    █   ██   ██ ██         ██
██████  █████ ██████   ██   ██ █████ ███████
`

// getversion returns a formatted string containing the version information from the provided AppConfig.
// The returned string is in the format: "Version : <version>\n".
func getversion(ac models.AppConfig) string {
	return fmt.Sprintf("Version : %s\n", ac.Version)
}

// getPort returns a formatted string indicating the server's host and port
// based on the provided ServerConfig. The output follows the format:
// "Server Running On <host>:<port>\n".
func getPort(sc models.ServerConfig) string {
	return fmt.Sprintf("Server Running On %s:%d\n", sc.Host, sc.Port)
}

// envVarMapper maps environment variable names by removing the "BT_" prefix
// and converting to the format expected by koanf.
// Example: BT_SENDGRIDAPIKEY → SENDGRIDAPIKEY
func envVarMapper(s string) string {
	return strings.TrimPrefix(s, "BT_")
}

// getConfigString returns env var value if set, otherwise falls back to config file value.
// Trims whitespace from all values.
func getConfigString(config *koanf.Koanf, envKey, configKey, defaultVal string) string {
	// First, check env var
	if val := strings.TrimSpace(config.String(envKey)); val != "" {
		return val
	}
	// Fallback to config file path
	if val := strings.TrimSpace(config.String(configKey)); val != "" {
		return val
	}
	return defaultVal
}

// getConfigInt returns env var value if set, otherwise falls back to config file value.
func getConfigInt(config *koanf.Koanf, envKey, configKey string, defaultVal int) int {
	if val := config.Int(envKey); val != 0 {
		return val
	}
	if val := config.Int(configKey); val != 0 {
		return val
	}
	return defaultVal
}

// getConfigBool retrieves a boolean value from a Koanf configuration, checking first
// for the presence of an environment variable (envKey), then a configuration file
// key (configKey), and falling back to a default value if neither exists.
//
// The environment variable value is normalized (lowercase, trimmed whitespace)
// and supports the following boolean representations:
// - True: "true", "1", "yes", "on"
// - False: "false", "0", "no", "off"
//
// If the environment variable value is not a recognized boolean string, it falls
// back to checking the configuration file key, then the default value.
func getConfigBool(config *koanf.Koanf, envKey, configKey string, defaultVal bool) bool {
	if config.Exists(envKey) {
		val := strings.ToLower(strings.TrimSpace(config.String(envKey)))
		switch val {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	if config.Exists(configKey) {
		return config.Bool(configKey)
	}

	return defaultVal
}

// getSecret retrieves a secret from environment only (no config file fallback)
// and trims whitespace
func getSecret(config *koanf.Koanf, envKey string) string {
	return strings.TrimSpace(config.String(envKey))
}

// validateConfig checks for required configuration values and logs warnings.
func validateConfig(app *models.App) {

	if !app.SMTPConfig.Enabled {
		log.Println("Info: SMTP disabled. Email features will not work.")
		return
	}

	var missingConfigs []string

	if app.SMTPConfig.Host == "" {
		missingConfigs = append(missingConfigs, "BT_SMTP_HOST")
	}
	if app.SMTPConfig.Port == 0 {
		missingConfigs = append(missingConfigs, "BT_SMTP_PORT")
	}
	if app.SMTPConfig.Username == "" {
		missingConfigs = append(missingConfigs, "BT_SMTP_USERNAME")
	}
	if app.SMTPConfig.Password == "" {
		missingConfigs = append(missingConfigs, "BT_SMTP_PASSWORD")
	}
	if app.SMTPConfig.FromEmail == "" {
		missingConfigs = append(missingConfigs, "BT_SMTP_FROM")
	}
	if app.AppConfig.BaseURL == "" {
		log.Println("Warning: BT_BASE_URL not set. Email links may not work correctly.")
	}
	if app.SentryConfig.Enabled && app.SentryConfig.DSN == "" {
		log.Println("Warning: Sentry enabled but DSN not configured. Sentry will be disabled.")
		app.SentryConfig.Enabled = false
	}
	if len(missingConfigs) > 0 {
		log.Printf("Warning: SMTP enabled but missing required config: %s",
			strings.Join(missingConfigs, ", "))
		log.Println("Email features will be disabled until all SMTP settings are configured.")
		app.SMTPConfig.Enabled = false // Disable to prevent errors
	}
}

// getBaseDir returns the current working directory.
func getBaseDir() string {
	baseDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return baseDir
}

// getDBBaseDir returns the path to the database directory.
func getDBBaseDir() string {
	return fp.Join(getBaseDir(), "src/db")
}
