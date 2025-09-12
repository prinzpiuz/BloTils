package models

import "github.com/gorilla/mux"

type Server struct {
	Router *mux.Router
	Config *ServerConfig
}

type ServerConfig struct {
	Host        string
	Port        int
	StaticFiles string
}
