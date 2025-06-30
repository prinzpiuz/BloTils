package app

import (
	"BloTils/src/db"
	"BloTils/src/models"
	"BloTils/src/server"
	"fmt"
	"log"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// envVarMapper maps environment variable names by removing the "BT_" prefix and replacing
// underscores with dots, to match the configuration key format.
func envVarMapper(s string) string {
	return strings.ReplaceAll(strings.TrimPrefix(s, "BT_"), "_", ".")
}

// InitApp initializes and returns a new instance of models.App using the provided Koanf configuration.
// It loads configuration values from a "config.json" file if it exists, and then overrides or supplements
// them with environment variables prefixed with "BT_". The function populates the App, ServerConfig,
// DBConfig, and EmailSettings fields of the models.App struct using the loaded configuration values.
// If any configuration loading fails, the function logs a fatal error and terminates the application.
//
// Parameters:
//   - config: A pointer to a koanf.Koanf instance containing configuration data.
//
// Returns:
//   - A pointer to an initialized models.App struct.
func InitApp(config *koanf.Koanf) *models.App {

	// Check if config.json exists before loading
	if _, err := os.Stat("config.json"); err == nil {
		if err := config.Load(file.Provider("config.json"), json.Parser()); err != nil {
			log.Fatalf("error loading config.json: %v", err)
		}
	}

	if err := config.Load(env.Provider("BT_", ".", envVarMapper), nil); err != nil {
		log.Fatalf("error loading environment variables: %v", err)
	}

	env := config.String("ENV")
	if env == "" {
		env = config.String("AppConfig.Env")
	}

	app := &models.App{
		AppConfig: models.AppConfig{
			Name:          config.String("AppConfig.Name"),
			Version:       config.String("AppConfig.Version"),
			Env:           env,
			BaseDirectory: getBaseDir(),
		},
		ServerConfig: models.ServerConfig{
			Host:        config.String("ServerConfig.Host"),
			Port:        config.Int("ServerConfig.Port"),
			StaticFiles: config.String("ServerConfig.StaticFiles"),
		},
		DBConfig: models.DBConfig{
			DBLocation:      config.String("DBConfig.DBLocation"),
			Vacuum:          config.String("DBConfig.Vacuum"),
			ForeignKeys:     config.Bool("DBConfig.ForeignKeys"),
			DBBaseDirectory: getDBBaseDir(),
			MigrationFiles:  "migrations",
		},
		EmailSettings: models.EmailSettings{
			FromMail:       config.String("EmailSettings.FromMail"),
			SendgridApiKey: config.String("SENDGRIDAPIKEY"),
		},
	}
	return app
}

// getBaseDir returns the current working directory as a string.
// It panics if there is an error retrieving the working directory.
func getBaseDir() string {
	baseDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return baseDir
}

// getDBBaseDir returns the absolute path to the database base directory by joining
// the application's base directory with the "src/db" subdirectory.
func getDBBaseDir() string {
	baseDir := getBaseDir()
	return fp.Join(baseDir, "src/db")
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
	fmt.Println()
}
