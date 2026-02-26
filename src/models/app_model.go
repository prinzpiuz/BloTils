package models

var (
	Version   = "unknown"
	CommitSHA = "unknown"
	BuildTime = "unknown"
)
var UNKNOWN = "unknown"

type AppConfig struct {
	Name          string
	Version       string
	Env           string
	BaseDirectory string
	BaseURL       string
}

type App struct {
	AppConfig    AppConfig
	ServerConfig ServerConfig
	DBConfig     DBConfig
	SMTPConfig   SMTPConfig
	SentryConfig SentryConfig
}

func (app App) IsEmpty() bool {
	return app == App{}
}

func (appConfig AppConfig) IsEmpty() bool {
	return appConfig == AppConfig{}
}

func (appConfig AppConfig) IsProduction() bool {
	return appConfig.Env == "production"
}
