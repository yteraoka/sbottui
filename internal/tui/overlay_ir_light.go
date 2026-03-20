package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/domain"
)

type irLightRow int

const (
	irLightRowOn irLightRow = iota
	irLightRowOff
	irLightRowCount
)

type irLightOverlay struct {
	item      domain.ListItem
	focus     irLightRow
	statusMsg string
}

func newIRLightOverlay(item domain.ListItem) *irLightOverlay {
	return &irLightOverlay{
		item:  item,
		focus: irLightRowOn,
	}
}

func (o *irLightOverlay) Title() string {
	return fmt.Sprintf("Light: %s", o.item.Name)
}

func (o *irLightOverlay) Update(msg tea.KeyMsg, client *api.Client) (Overlay, tea.Cmd) {
	switch {
	case key.Matches(msg, overlayKeys.Up):
		if o.focus > 0 {
			o.focus--
		}
	case key.Matches(msg, overlayKeys.Down):
		if o.focus < irLightRowCount-1 {
			o.focus++
		}
	case key.Matches(msg, overlayKeys.Select):
		return o, o.sendCommand(client)
	}
	return o, nil
}

func (o *irLightOverlay) sendCommand(client *api.Client) tea.Cmd {
	deviceID := o.item.ID
	var command string
	switch o.focus {
	case irLightRowOn:
		command = "turnOn"
	case irLightRowOff:
		command = "turnOff"
	}

	return func() tea.Msg {
		cmd := api.CommandRequest{
			Command:     command,
			CommandType: "command",
		}
		_, err := client.SendIRCommand(deviceID, cmd)
		return MsgCommandDone{Err: err}
	}
}

func (o *irLightOverlay) View() string {
	var sb strings.Builder

	rows := []struct {
		label string
		row   irLightRow
	}{
		{"Turn On", irLightRowOn},
		{"Turn Off", irLightRowOff},
	}

	for _, r := range rows {
		line := fmt.Sprintf("%-20s", r.label)
		if o.focus == r.row {
			sb.WriteString(styleOverlayFocused.Render("> " + line))
		} else {
			sb.WriteString(styleOverlayNormal.Render("  " + line))
		}
		sb.WriteString("\n")
	}

	if o.statusMsg != "" {
		sb.WriteString("\n")
		sb.WriteString(styleTypeLabel.Render(o.statusMsg))
	}

	sb.WriteString("\n")
	sb.WriteString(styleTypeLabel.Render("[↑↓] move  [enter] send  [esc] close"))

	return sb.String()
}
