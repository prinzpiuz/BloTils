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

func InitApp(config *koanf.Koanf) *models.App {

	err := config.Load(file.Provider("config.json"), json.Parser())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	config.Load(env.Provider("BT_", ".", func(s string) string {
		return strings.Replace(strings.TrimPrefix(s, "BT_"), "_", ".", -1)
	}), nil)
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
func getBaseDir() string {
	baseDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return baseDir
}

func getDBBaseDir() string {
	baseDir := getBaseDir()
	return fp.Join(baseDir, "src/db")
}

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
