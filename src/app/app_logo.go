package app

import (
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

func getversion(c Config) string {
	return fmt.Sprintf("Version : %s\n", c.Version)
}

func getPort(c Config) string {
	return fmt.Sprintf("Server Running On %s:%d\n", c.ServerConfig.Host, c.ServerConfig.Port)
}

func logo(c Config) {
	fmt.Print(logoStyle.Render(appLogo))
	fmt.Println()
	fmt.Print("Utilities For Your Blog Engine\n")
	fmt.Print(getversion(c))
	fmt.Print(getPort(c))
	fmt.Println()
	fmt.Println()
}
