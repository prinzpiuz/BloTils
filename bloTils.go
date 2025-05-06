package main

import (
	"BloTils/src/app"
	"flag"
	"os"

	"github.com/knadh/koanf/v2"
)

func main() {
	var config = koanf.New(".")
	app := app.InitApp(config)
	createAdmin := flag.Bool("createadmin", false, "Create A Admin Account")
	flag.Parse()
	if *createAdmin {
		app.CreateAdmin()
		os.Exit(0)
	}
	app.Start()
}
