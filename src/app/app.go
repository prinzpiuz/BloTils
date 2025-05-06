package app

import (
	"BloTils/src/db"
	"BloTils/src/email"
	"BloTils/src/server"
	"log"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type App struct {
	AppConfig     AppConfig
	ServerConfig  server.ServerConfig
	DBConfig      db.DBConfig
	EmailSettings email.EmailSettings
}

type AppConfig struct {
	Name          string
	Version       string
	Env           string
	BaseDirectory string
}

func InitApp(config *koanf.Koanf) *App {

	err := config.Load(file.Provider("config.json"), json.Parser())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	app := &App{
		AppConfig: AppConfig{
			Name:          config.String("AppConfig.Name"),
			Version:       config.String("AppConfig.Version"),
			Env:           config.String("AppConfig.Env"),
			BaseDirectory: config.String("AppConfig.BaseDirectory"),
		},
		ServerConfig: server.ServerConfig{
			Host:        config.String("ServerConfig.Host"),
			Port:        config.Int("ServerConfig.Port"),
			StaticFiles: config.String("ServerConfig.StaticFiles"),
		},
		DBConfig: db.DBConfig{
			DBLocation:  config.String("DBConfig.DBLocation"),
			Vacuum:      config.String("DBConfig.Vacuum"),
			ForeignKeys: config.Bool("DBConfig.ForeignKeys"),
		},
		EmailSettings: email.EmailSettings{
			FromMail:       config.String("EmailSettings.FromMail"),
			SendgridApiKey: config.String("EmailSettings.SendgridApiKey"),
		},
	}
	return app
}

func (app App) Start() {
	logo(app)
	err := db.InitDB(app.DBConfig)
	if err != nil {
		log.Fatalf("Error Initializing DB: %v", err)
	}
	server := server.InitServer(app.ServerConfig)
	// setup the middlewares
	server.Start()
}

func (app App) CreateAdmin() {
	logo(app)
	// handlers.CreateAdmin(app.Server.Config.DB.Connection)
}

// func LoadConfig(filepath string) Config {
// 	config, err := loadConfigFromFile(filepath)
// 	if err != nil {
// 		fmt.Println(err)
// 		return Config{}
// 	}
// 	baseDir, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}
// 	config.BaseDirectory = baseDir
// 	config.ServerConfig.BaseDirectory = baseDir
// 	config.ServerConfig.DB.DBBaseDirectory = fp.Join(baseDir, "src/db")
// 	config.ServerConfig.IsProduction = config.Env == "production"
// 	return config
// }

// var ErrLoadingConfig = fmt.Errorf("error loading config from file")

// func loadConfigFromFile(fpath string) (Config, error) {
// 	fp := fp.Clean(fpath)
// 	jsonFile, err := os.Open(fp)
// 	if err != nil {
// 		return Config{}, fmt.Errorf("%w: %v", ErrLoadingConfig, err)
// 	}
// 	defer func() {
// 		if err := jsonFile.Close(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()

// 	var config Config
// 	jsonParser := json.NewDecoder(jsonFile)
// 	err = jsonParser.Decode(&config)
// 	if err != nil {
// 		return Config{}, fmt.Errorf("%w: %v", ErrLoadingConfig, err)
// 	}
// 	return config, nil
// }

// func New(config Config) *App {
// 	server := server.New(config.ServerConfig)
// 	app := &App{
// 		Server: *server,
// 		Config: config,
// 	}
// 	routes.RegisterRoutes(server)
// 	routes.ServeStaticFiles(server)
// 	return app
// }
