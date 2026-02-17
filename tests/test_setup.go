package handlers

import (
	"BloTils/src/models"
)

var TestDBConfig = models.DBConfig{
	DBLocation:      "./Test_BloTils.db",
	Vacuum:          "full",
	ForeignKeys:     true,
	DBBaseDirectory: "../src/db",
}

var ServerConfig = models.ServerConfig{
	Host:        "",
	Port:        8080,
	StaticFiles: "static",
}

var AppConfig = models.AppConfig{
	Name:          "BloTils",
	Version:       "0.0.1",
	Env:           "test",
	BaseDirectory: "../",
	BaseURL:       "127.0.0.1:8000",
}

var SMTPConfig = models.SMTPConfig{}

var App = &models.App{
	AppConfig:    AppConfig,
	ServerConfig: ServerConfig,
	DBConfig:     TestDBConfig,
	SMTPConfig:   SMTPConfig,
}
