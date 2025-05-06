package app

import (
	"BloTils/src/server"
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

func getversion(ac AppConfig) string {
	return fmt.Sprintf("Version : %s\n", ac.Version)
}

func getPort(sc server.ServerConfig) string {
	return fmt.Sprintf("Server Running On %s:%d\n", sc.Host, sc.Port)
}

func logo(a App) {
	fmt.Print(logoStyle.Render(appLogo))
	fmt.Println()
	fmt.Print("Utilities For Your Static Blog\n")
	fmt.Print(getversion(a.AppConfig))
	fmt.Print(getPort(a.ServerConfig))
	fmt.Println()
	fmt.Println()
}
