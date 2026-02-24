package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/tui"
)

func main() {
	token := os.Getenv("SWITCHBOT_TOKEN")
	secret := os.Getenv("SWITCHBOT_CLIENT_SECRET")

	if token == "" || secret == "" {
		fmt.Fprintln(os.Stderr, "Error: SWITCHBOT_TOKEN and SWITCHBOT_CLIENT_SECRET must be set")
		os.Exit(1)
	}

	client := api.NewClient(token, secret)
	model := tui.New(client)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
