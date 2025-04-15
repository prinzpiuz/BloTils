package main

import (
	"BloTils/src/app"
	"flag"
	"os"
	"path/filepath"
)

func main() {
	fp := filepath.Join(".", "config.json")
	config := app.LoadConfig(fp)
	app := app.New(config)
	createAdmin := flag.Bool("createadmin", false, "Create A Admin Account")
	flag.Parse()
	if *createAdmin {
		app.CreateAdmin()
		os.Exit(0)
	}
	app.Start()
}
