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

// getversion returns a formatted string containing the version information from the provided AppConfig.
// The returned string is in the format: "Version : <version>\n".
func getversion(ac models.AppConfig) string {
	return fmt.Sprintf("Version : %s\n", ac.Version)
}

// getPort returns a formatted string indicating the server's host and port
// based on the provided ServerConfig. The output follows the format:
// "Server Running On <host>:<port>\n".
func getPort(sc models.ServerConfig) string {
	return fmt.Sprintf("Server Running On %s:%d\n", sc.Host, sc.Port)
}
