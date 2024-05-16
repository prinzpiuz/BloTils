package main

import (
	"log"
	"path/filepath"

	"BloTils/src/app"
)

func main() {
	fp := filepath.Join(".", "config.json")
	config := app.LoadConfig(fp)
	app := app.New(config)
	err := app.Config.ServerConfig.DB.Initalize()
	if err != nil {
		log.Fatal(err)
	}
	app.Start()
}
