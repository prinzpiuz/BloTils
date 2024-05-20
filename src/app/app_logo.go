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

func logo(c Config) {
	fmt.Print(logoStyle.Render(appLogo))
	fmt.Println()
	fmt.Print("Utilities For Your Blog Engine\n")
	fmt.Print(fmt.Sprintf("Version : %s\n", c.Version))
	fmt.Print(fmt.Sprintf("Server Running On %s:%d\n", c.ServerConfig.Host, c.ServerConfig.Port))
	fmt.Println()
	fmt.Println()
}
