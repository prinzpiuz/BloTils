package main

import (
	"BloTils/src/app"
	"path/filepath"
)

func main() {
	fp := filepath.Join(".", "config.json")
	config := app.LoadConfig(fp)
	app := app.New(config)
	app.Start()
}
