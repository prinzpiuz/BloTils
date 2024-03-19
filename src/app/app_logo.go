package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var logoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))
var subtitleStyle = lipgloss.NewStyle().
	Padding(0, 22).
	Bold(true)

const appLogo = `
				 ██████  ██      ███████ ███████ ███████  █    █
				 ██      ██      ██   ██ ██   ██ ██   ██  ██  ██
				 ██      ██      ███████ ██████  ██████    ████
				 ██      ██      ██   ██ ██      ██         ██
				 ██████ ███████  ██   ██ ██      ██         ██ 👏
`

func logo(c Config) {
	fmt.Print(logoStyle.Render(appLogo))
	fmt.Println()
	fmt.Printf("Version : %s", c.Version)
	fmt.Print(subtitleStyle.Render(fmt.Sprintf("Server Runing On %s:%d", c.ServerConfig.Host, c.ServerConfig.Port)))
	fmt.Println()
	fmt.Println()
}
