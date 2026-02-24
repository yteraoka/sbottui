package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/domain"
)

type plugMiniRow int

const (
	plugMiniRowOn plugMiniRow = iota
	plugMiniRowOff
	plugMiniRowCount
)

type plugMiniOverlay struct {
	item        domain.ListItem
	focus       plugMiniRow
	loading     bool
	power       string
	version     string
	hubDeviceID string
	voltage     float64
	weight      float64
	statusMsg   string
}

func newPlugMiniOverlay(item domain.ListItem, client *api.Client) (*plugMiniOverlay, tea.Cmd) {
	o := &plugMiniOverlay{
		item:    item,
		focus:   plugMiniRowOn,
		loading: true,
	}
	cmd := func() tea.Msg {
		status, err := client.GetDeviceStatus(item.ID)
		return MsgDeviceStatus{Status: status, Err: err}
	}
	return o, cmd
}

func (o *plugMiniOverlay) ApplyStatus(status *api.DeviceStatus) {
	o.loading = false
	o.power = status.Power
	o.version = status.Version
	o.hubDeviceID = status.HubDeviceID
	o.voltage = status.Voltage
	o.weight = status.Weight
	// Set initial focus based on current power state
	if o.power == "off" {
		o.focus = plugMiniRowOn // suggest turning on
	} else {
		o.focus = plugMiniRowOff // suggest turning off
	}
}

func (o *plugMiniOverlay) Title() string {
	return fmt.Sprintf("Plug Mini: %s", o.item.Name)
}

func (o *plugMiniOverlay) Update(msg tea.KeyMsg, client *api.Client) (Overlay, tea.Cmd) {
	switch {
	case key.Matches(msg, overlayKeys.Up):
		if o.focus > 0 {
			o.focus--
		}
	case key.Matches(msg, overlayKeys.Down):
		if o.focus < plugMiniRowCount-1 {
			o.focus++
		}
	case key.Matches(msg, overlayKeys.Select):
		return o, o.sendCommand(client)
	}
	return o, nil
}

func (o *plugMiniOverlay) sendCommand(client *api.Client) tea.Cmd {
	deviceID := o.item.ID
	var command string
	switch o.focus {
	case plugMiniRowOn:
		command = "turnOn"
	case plugMiniRowOff:
		command = "turnOff"
	}

	return func() tea.Msg {
		cmd := api.CommandRequest{
			Command:     command,
			CommandType: "command",
		}
		_, err := client.SendCommand(deviceID, cmd)
		return MsgCommandDone{Err: err}
	}
}

func (o *plugMiniOverlay) View() string {
	if o.loading {
		return styleOverlayNormal.Render("Loading status...")
	}

	var sb strings.Builder

	// Info fields
	info := []struct{ label, value string }{
		{"Device ID", o.item.ID},
		{"Hub Device ID", o.hubDeviceID},
		{"Version", o.version},
		{"Power", o.power},
		{"Voltage", fmt.Sprintf("%.1f V", o.voltage)},
		{"Wattage", fmt.Sprintf("%.1f W", o.weight)},
	}
	for _, f := range info {
		sb.WriteString(styleTypeLabel.Render(fmt.Sprintf("%-16s", f.label)))
		sb.WriteString(styleOverlayNormal.Render(f.value))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Control rows
	rows := []struct {
		label string
		row   plugMiniRow
	}{
		{"Turn On", plugMiniRowOn},
		{"Turn Off", plugMiniRowOff},
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

	sb.WriteString("\n")
	sb.WriteString(styleTypeLabel.Render("[↑↓] move  [enter] send  [esc] close"))

	return sb.String()
}
