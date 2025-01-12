package app

import (
	"BloTils/src/server"
	"BloTils/src/server/routes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type App struct {
	Server server.Server
	Config Config
}

type Config struct {
	ServerConfig server.ServerConfig
	Name         string
	Version      string
	Env          string
}

func (app *App) Start() {
	logo(app.Config)
	app.Server.Start()
}

func LoadConfig(filepath string) Config {
	config, err := loadConfigFromFile(filepath)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}
	return config
}

var ErrLoadingConfig = fmt.Errorf("error loading config from file")

func loadConfigFromFile(fpath string) (Config, error) {
	fp := filepath.Clean(fpath)
	jsonFile, err := os.Open(fp)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrLoadingConfig, err)
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var config Config
	jsonParser := json.NewDecoder(jsonFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrLoadingConfig, err)
	}
	return config, nil
}

func New(config Config) *App {
	server := server.New(config.ServerConfig)
	app := &App{
		Server: *server,
		Config: config,
	}
	routes.RegisterRoutes(server)
	routes.ServeStaticFiles(server)
	return app
}
