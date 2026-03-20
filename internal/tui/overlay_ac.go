package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/domain"
)

type acRow int

const (
	acRowPower acRow = iota
	acRowTemp
	acRowMode
	acRowFan
	acRowCount
)

var acModes = []string{"auto", "cool", "dry", "fan", "heat"}
var acFanSpeeds = []string{"auto", "low", "medium", "high"}

type acOverlay struct {
	item      domain.ListItem
	focus     acRow
	power     string // "on" or "off"
	temp      int    // 16-30
	modeIdx   int    // 0-4 (index into acModes)
	fanIdx    int    // 0-3 (index into acFanSpeeds)
	statusMsg string
}

func newACOverlay(item domain.ListItem) *acOverlay {
	return &acOverlay{
		item:    item,
		power:   "on",
		temp:    26,
		modeIdx: 1, // cool
		fanIdx:  0, // auto
	}
}

func (o *acOverlay) Title() string {
	return fmt.Sprintf("Air Conditioner: %s", o.item.Name)
}

func (o *acOverlay) Update(msg tea.KeyMsg, client *api.Client) (Overlay, tea.Cmd) {
	switch {
	case key.Matches(msg, overlayKeys.Up):
		if o.focus > 0 {
			o.focus--
		}
	case key.Matches(msg, overlayKeys.Down):
		if o.focus < acRowCount-1 {
			o.focus++
		}
	case key.Matches(msg, overlayKeys.Left):
		o.adjustLeft()
	case key.Matches(msg, overlayKeys.Right):
		o.adjustRight()
	case key.Matches(msg, overlayKeys.Select):
		return o, o.sendCommand(client)
	}
	return o, nil
}

func (o *acOverlay) adjustLeft() {
	switch o.focus {
	case acRowPower:
		if o.power == "on" {
			o.power = "off"
		} else {
			o.power = "on"
		}
	case acRowTemp:
		if o.temp > 16 {
			o.temp--
		}
	case acRowMode:
		if o.modeIdx > 0 {
			o.modeIdx--
		}
	case acRowFan:
		if o.fanIdx > 0 {
			o.fanIdx--
		}
	}
}

func (o *acOverlay) adjustRight() {
	switch o.focus {
	case acRowPower:
		if o.power == "on" {
			o.power = "off"
		} else {
			o.power = "on"
		}
	case acRowTemp:
		if o.temp < 30 {
			o.temp++
		}
	case acRowMode:
		if o.modeIdx < len(acModes)-1 {
			o.modeIdx++
		}
	case acRowFan:
		if o.fanIdx < len(acFanSpeeds)-1 {
			o.fanIdx++
		}
	}
}

func (o *acOverlay) sendCommand(client *api.Client) tea.Cmd {
	deviceID := o.item.ID
	// SwitchBot AC setAll: "temp,mode,fan,power"
	// mode: 1=auto,2=cool,3=dry,4=fan,5=heat
	// fan: 1=auto,2=low,3=medium,4=high
	// power: on=1, off=0
	temp := o.temp
	modeNum := o.modeIdx + 1
	fanNum := o.fanIdx + 1
	power := "on"
	if o.power == "off" {
		power = "off"
	}

	return func() tea.Msg {
		parameter := fmt.Sprintf("%d,%d,%d,%s", temp, modeNum, fanNum, power)
		cmd := api.CommandRequest{
			Command:     "setAll",
			Parameter:   parameter,
			CommandType: "command",
		}
		_, err := client.SendIRCommand(deviceID, cmd)
		return MsgCommandDone{Err: err}
	}
}

func (o *acOverlay) View() string {
	var sb strings.Builder

	rows := []struct {
		label string
		value string
		row   acRow
	}{
		{"Power", o.power + "  [←→ toggle]", acRowPower},
		{"Temperature", fmt.Sprintf("%d°C  [←-1  →+1]  (16-30)", o.temp), acRowTemp},
		{"Mode", fmt.Sprintf("%s  [←→]  (%s)", acModes[o.modeIdx], strings.Join(acModes, "/")), acRowMode},
		{"Fan Speed", fmt.Sprintf("%s  [←→]  (%s)", acFanSpeeds[o.fanIdx], strings.Join(acFanSpeeds, "/")), acRowFan},
	}

	for _, r := range rows {
		label := styleOverlayLabel.Render(r.label)
		line := fmt.Sprintf("%s %s", label, r.value)
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
	sb.WriteString(styleTypeLabel.Render("[↑↓] move  [←→] adjust  [enter] send setAll  [esc] close"))

	return sb.String()
}
