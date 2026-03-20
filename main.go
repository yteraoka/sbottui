package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/tui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("version: %s\ncommit:  %s\ndate:    %s\n", version, commit, date)
		return
	}

	token := os.Getenv("SWITCHBOT_TOKEN")
	secret := os.Getenv("SWITCHBOT_CLIENT_SECRET")

	if token == "" || secret == "" {
		fmt.Fprintln(os.Stderr, "Error: SWITCHBOT_TOKEN and SWITCHBOT_CLIENT_SECRET must be set")
		os.Exit(1)
	}

	client := api.NewClient(token, secret)
	model := tui.New(client)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
