package handlers

import (
	"BloTils/src/app"
	"BloTils/src/db"
	"BloTils/src/server"
)

var TestDBConfig = db.DB{
	DBLocation:      "./Test_BloTils.db",
	Vacuum:          "full",
	ForeignKeys:     true,
	DBBaseDirectory: "../src/db",
}

var ServerConfig = server.ServerConfig{
	Host:          "",
	Port:          8080,
	DB:            TestDBConfig,
	StaticFiles:   "static",
	BaseDirectory: "../",
}

var AppConfig = app.Config{
	ServerConfig:  ServerConfig,
	Name:          "BloTils",
	Version:       "0.0.1",
	Env:           "test",
	BaseDirectory: "../",
}
