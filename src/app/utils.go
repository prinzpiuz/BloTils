package app

import (
	"BloTils/src/models"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var logoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))

const appLogo = `
██████  ██    ██████ ██████ ██ ██    ███████
██  ██  ██    █    █   ██   ██ ██    ██
██████  ██    █    █   ██   ██ ██    ███████
██  ██  ██    █    █   ██   ██ ██         ██
██████  █████ ██████   ██   ██ █████ ███████
`

func getversion(ac models.AppConfig) string {
	return fmt.Sprintf("Version : %s\n", ac.Version)
}

func getPort(sc models.ServerConfig) string {
	return fmt.Sprintf("Server Running On %s:%d\n", sc.Host, sc.Port)
}
