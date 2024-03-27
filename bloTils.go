package main

import (
	"path/filepath"

	"BloTils/src/app"
)

func main() {
	fp := filepath.Join(".", "config.json")
	config := app.LoadConfig(fp)
	app := app.New(config)
	app.Start()

}
