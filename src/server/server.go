package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prinzpiuz/Clappy/src/server/handlers"
)

type Server struct {
	Router *mux.Router
	Config ServerConfig
}

type ServerConfig struct {
	Host string
	Port int
	DB   map[string]string
}

func (server *Server) Start() {
	err := http.ListenAndServe(server.Config.Addr(), server.Router)
	print(err)
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func New(config ServerConfig) *Server {
	router := mux.NewRouter()
	server := &Server{
		Router: router,
		Config: config,
	}
	server.RegisterRoutes()
	return server
}

func (server *Server) RegisterRoutes() {
	server.Router.HandleFunc("/v1/ping", handlers.Ping).Methods("GET")
	server.Router.HandleFunc("/v1/count/{url}", handlers.GetClaps).Methods("GET")
}
