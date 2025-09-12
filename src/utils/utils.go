package utils

import (
	"BloTils/src/app"
	"BloTils/src/db"
	"BloTils/src/models"
	"BloTils/src/server"
	"log"
)

// CreateAdmin initializes the database and creates an admin user for the application.
// It takes an application configuration, displays the logo, sets up the database connection,
// and then creates an admin user using the provided database connection.
func CreateAdmin(appConfig models.App) {
	app.Logo(appConfig)
	err := db.InitDB(&appConfig.DBConfig)
	if err != nil {
		log.Fatalf("Error Initializing DB: %v", err)
	}
	server.CreateAdmin(appConfig.DBConfig.Connection)
}
