package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/domain"
)

type bulbRow int

const (
	bulbRowPower bulbRow = iota
	bulbRowBrightness
	bulbRowColorTemp
	bulbRowCount
)

type bulbOverlay struct {
	item       domain.ListItem
	focus      bulbRow
	power      string // "on" or "off"
	brightness int    // 1-100
	colorTemp  int    // 2700-6500
	version    string
	loading    bool
	statusMsg  string
}

func newBulbOverlay(item domain.ListItem, client *api.Client) (*bulbOverlay, tea.Cmd) {
	o := &bulbOverlay{
		item:       item,
		power:      "on",
		brightness: 50,
		colorTemp:  4000,
		loading:    true,
	}
	cmd := func() tea.Msg {
		status, err := client.GetDeviceStatus(item.ID)
		return MsgDeviceStatus{Status: status, Err: err}
	}
	return o, cmd
}

func (o *bulbOverlay) ApplyStatus(status *api.DeviceStatus) {
	o.loading = false
	o.power = status.Power
	if status.Brightness > 0 {
		o.brightness = status.Brightness
	}
	if status.ColorTemp > 0 {
		o.colorTemp = status.ColorTemp
	}
	o.version = status.Version
}

func (o *bulbOverlay) Title() string {
	return fmt.Sprintf("Color Bulb: %s", o.item.Name)
}

func (o *bulbOverlay) Update(msg tea.KeyMsg, client *api.Client) (Overlay, tea.Cmd) {
	switch {
	case key.Matches(msg, overlayKeys.Up):
		if o.focus > 0 {
			o.focus--
		}
	case key.Matches(msg, overlayKeys.Down):
		if o.focus < bulbRowCount-1 {
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

func (o *bulbOverlay) adjustLeft() {
	switch o.focus {
	case bulbRowPower:
		if o.power == "on" {
			o.power = "off"
		} else {
			o.power = "on"
		}
	case bulbRowBrightness:
		o.brightness -= 5
		if o.brightness < 1 {
			o.brightness = 1
		}
	case bulbRowColorTemp:
		o.colorTemp -= 100
		if o.colorTemp < 2700 {
			o.colorTemp = 2700
		}
	}
}

func (o *bulbOverlay) adjustRight() {
	switch o.focus {
	case bulbRowPower:
		if o.power == "on" {
			o.power = "off"
		} else {
			o.power = "on"
		}
	case bulbRowBrightness:
		o.brightness += 5
		if o.brightness > 100 {
			o.brightness = 100
		}
	case bulbRowColorTemp:
		o.colorTemp += 100
		if o.colorTemp > 6500 {
			o.colorTemp = 6500
		}
	}
}

func (o *bulbOverlay) sendCommand(client *api.Client) tea.Cmd {
	deviceID := o.item.ID
	power := o.power
	brightness := o.brightness
	colorTemp := o.colorTemp
	focus := o.focus

	return func() tea.Msg {
		var cmd api.CommandRequest
		switch focus {
		case bulbRowPower:
			if power == "on" {
				cmd = api.CommandRequest{Command: "turnOn", CommandType: "command"}
			} else {
				cmd = api.CommandRequest{Command: "turnOff", CommandType: "command"}
			}
		case bulbRowBrightness:
			cmd = api.CommandRequest{
				Command:     "setBrightness",
				Parameter:   fmt.Sprintf("%d", brightness),
				CommandType: "command",
			}
		case bulbRowColorTemp:
			cmd = api.CommandRequest{
				Command:     "setColorTemperature",
				Parameter:   fmt.Sprintf("%d", colorTemp),
				CommandType: "command",
			}
		}
		_, err := client.SendCommand(deviceID, cmd)
		return MsgCommandDone{Err: err}
	}
}

func (o *bulbOverlay) View() string {
	if o.loading {
		return styleOverlayNormal.Render("Loading status...")
	}

	var sb strings.Builder

	if o.version != "" {
		sb.WriteString(styleTypeLabel.Render(fmt.Sprintf("ID: %s  Version: %s", o.item.ID, o.version)))
		sb.WriteString("\n\n")
	}

	rows := []struct {
		label string
		value string
		row   bulbRow
	}{
		{"Power", o.power, bulbRowPower},
		{"Brightness", fmt.Sprintf("%d%%  [←-5  →+5]", o.brightness), bulbRowBrightness},
		{"Color Temp", fmt.Sprintf("%dK  [←-100  →+100]", o.colorTemp), bulbRowColorTemp},
	}

	for _, r := range rows {
		label := styleOverlayLabel.Render(r.label)
		value := r.value
		line := fmt.Sprintf("%s %s", label, value)

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
	sb.WriteString(styleTypeLabel.Render("[↑↓] move  [←→] adjust  [enter] send  [esc] close"))

	return sb.String()
}
