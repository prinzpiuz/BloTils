package main

import (
	"BloTils/src/app"
	"BloTils/src/utils"
	"flag"
	"os"

	"github.com/knadh/koanf/v2"
)

func main() {
	var config = koanf.New(".")
	appConfig := app.InitApp(config)
	createAdmin := flag.Bool("createadmin", false, "Create A Admin Account")
	flag.Parse()
	if *createAdmin {
		utils.CreateAdmin(*appConfig)
		os.Exit(0)
	}
	app.Start(*appConfig)
}
