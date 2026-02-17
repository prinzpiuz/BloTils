package app

import (
	"BloTils/src/db"
	mailer "BloTils/src/email"
	"BloTils/src/models"
	"BloTils/src/sentry"
	"BloTils/src/server"
	"fmt"
	"log"
	"os"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// InitApp initializes and returns a new instance of models.App.
//
// Configuration precedence (highest to lowest):
//  1. Environment variables (prefixed with BT_)
//  2. config.json file
//  3. Default values
//
// Environment variables:
//   - BT_ENV              : Application environment (dev, staging, production)
//   - BT_DBLOCATION       : Database file path
//   - BT_HOST             : Server host
//   - BT_PORT             : Server port
//   - BT_SENDGRIDAPIKEY   : SendGrid API key (secret, env only)
//   - BT_FROMEMAIL        : Email sender address
func InitApp(config *koanf.Koanf) *models.App {
	// Load config.json if it exists (provides defaults)
	if _, err := os.Stat("config.json"); err == nil {
		if err := config.Load(file.Provider("config.json"), json.Parser()); err != nil {
			log.Fatalf("error loading config.json: %v", err)
		}
	}

	// Load environment variables (overrides config.json)
	if err := config.Load(env.Provider("BT_", ".", envVarMapper), nil); err != nil {
		log.Fatalf("error loading environment variables: %v", err)
	}
	app := &models.App{
		AppConfig: models.AppConfig{
			Name:          getConfigString(config, "NAME", "AppConfig.Name", "BloTils"),
			Version:       getConfigString(config, "VERSION", "AppConfig.Version", "dev"),
			Env:           getConfigString(config, "ENV", "AppConfig.Env", "dev"),
			BaseURL:       getConfigString(config, "BASE_URL", "AppConfig.BaseURL", "http://localhost:8000"),
			BaseDirectory: getBaseDir(),
		},
		ServerConfig: models.ServerConfig{
			Host:        getConfigString(config, "HOST", "ServerConfig.Host", ""),
			Port:        getConfigInt(config, "PORT", "ServerConfig.Port", 8000),
			StaticFiles: getConfigString(config, "STATICFILES", "ServerConfig.StaticFiles", "static"),
		},
		DBConfig: models.DBConfig{
			DBLocation:      getConfigString(config, "DBLOCATION", "DBConfig.DBLocation", "./BloTils.db"),
			Vacuum:          getConfigString(config, "DBVACUUM", "DBConfig.Vacuum", "full"),
			ForeignKeys:     getConfigBool(config, "DBFOREIGNKEYS", "DBConfig.ForeignKeys", true),
			DBBaseDirectory: getDBBaseDir(),
			MigrationFiles:  getConfigString(config, "MIGRATIONS", "DBConfig.MigrationFiles", "migrations"),
		},
		SMTPConfig: models.SMTPConfig{
			Enabled:   getConfigBool(config, "SMTP_ENABLED", "SMTPConfig.Enabled", false),
			Host:      getConfigString(config, "SMTP_HOST", "SMTPConfig.Host", ""),
			Port:      getConfigInt(config, "SMTP_PORT", "SMTPConfig.Port", 587),
			Username:  getSecret(config, "SMTP_USERNAME"),
			Password:  getSecret(config, "SMTP_PASSWORD"),
			FromEmail: getConfigString(config, "SMTP_FROM", "SMTPConfig.FromEmail", "noreply@blotils.com"),
			FromName:  getConfigString(config, "SMTP_FROMNAME", "SMTPConfig.FromName", "BloTils"),
			BaseURL:   getConfigString(config, "BASE_URL", "AppConfig.BaseURL", "http://localhost:8000"),
		},
		SentryConfig: models.SentryConfig{
			Enabled: getConfigBool(config, "SENTRY_ENABLED", "SentryConfig.Enabled", false),
			DSN:     getSecret(config, "SENTRY_DSN"),
			Debug:   getConfigBool(config, "SENTRY_DEBUG", "SentryConfig.Debug", false),
		},
	}

	// Validate required configuration
	validateConfig(app)

	return app
}

// Start initializes and starts the application using the provided App configuration.
// It displays the application logo, initializes the database, sets up the server with
// middlewares and routes, serves static files, and starts the server. If database
// initialization fails, the function logs a fatal error and terminates the application.
//
// Parameters:
//   - app: The application configuration containing database and server settings.
func Start(app models.App) {
	Logo(app)
	err := db.InitDB(&app.DBConfig)
	if err != nil {
		log.Fatalf("Error Initializing DB: %v", err)
	}
	mailer.Initialize(&app.SMTPConfig)
	sentry.Initialize(&app.SentryConfig)
	// anything should be done before starting the server can be added here (e.g. background tasks, etc.)
	// Initialize server and register routes
	serverConfig := server.InitServer(&app.ServerConfig)
	server.SetupMiddlewares(serverConfig.Router, app)
	server.RegisterRoutes(serverConfig)
	server.ServeStaticFiles(serverConfig)
	server.Start(serverConfig)
}

// Logo prints the application logo, version, server port, and environment information
// to the standard output. It uses the provided App model to access configuration details.
func Logo(a models.App) {
	fmt.Print(logoStyle.Render(appLogo))
	fmt.Println()
	fmt.Print("Utilities For Your Static Blog\n")
	fmt.Print(getversion(a.AppConfig))
	fmt.Print(getPort(a.ServerConfig))
	fmt.Printf("Environment: %s", a.AppConfig.Env)
	fmt.Println()
	fmt.Printf("Debug: %t", a.SentryConfig.Debug)
	fmt.Println()
	fmt.Println()
}
