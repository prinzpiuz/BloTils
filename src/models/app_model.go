package models

type App struct {
	AppConfig     AppConfig
	ServerConfig  ServerConfig
	DBConfig      DBConfig
	EmailSettings EmailSettings
}

func (app App) IsEmpty() bool {
	return app == App{}
}

type AppConfig struct {
	Name          string
	Version       string
	Env           string
	BaseDirectory string
}

func (appConfig AppConfig) IsEmpty() bool {
	return appConfig == AppConfig{}
}

func (appConfig AppConfig) IsProduction() bool {
	return appConfig.Env == "production"
}
