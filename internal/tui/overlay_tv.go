package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/domain"
)

type tvRow int

const (
	tvRowPower tvRow = iota
	tvRowVolume
	tvRowChannel
	tvRowMute
	tvRowCount
)

type tvButton struct {
	label   string
	command string
}

var tvRows = []struct {
	name    string
	row     tvRow
	buttons []tvButton
}{
	{
		name: "Power",
		row:  tvRowPower,
		buttons: []tvButton{
			{"On", "turnOn"},
			{"Off", "turnOff"},
		},
	},
	{
		name: "Volume",
		row:  tvRowVolume,
		buttons: []tvButton{
			{"Vol-", "volumeSub"},
			{"Vol+", "volumeAdd"},
		},
	},
	{
		name: "Channel",
		row:  tvRowChannel,
		buttons: []tvButton{
			{"Ch-", "channelSub"},
			{"Ch+", "channelAdd"},
		},
	},
	{
		name: "Mute",
		row:  tvRowMute,
		buttons: []tvButton{
			{"Mute", "setMute"},
		},
	},
}

type tvOverlay struct {
	item      domain.ListItem
	focusRow  tvRow
	focusCol  int
	statusMsg string
}

func newTVOverlay(item domain.ListItem) *tvOverlay {
	return &tvOverlay{
		item:     item,
		focusRow: tvRowPower,
		focusCol: 0,
	}
}

func (o *tvOverlay) Title() string {
	return fmt.Sprintf("TV: %s", o.item.Name)
}

func (o *tvOverlay) Update(msg tea.KeyMsg, client *api.Client) (Overlay, tea.Cmd) {
	switch {
	case key.Matches(msg, overlayKeys.Up):
		if o.focusRow > 0 {
			o.focusRow--
			row := tvRows[o.focusRow]
			if o.focusCol >= len(row.buttons) {
				o.focusCol = len(row.buttons) - 1
			}
		}
	case key.Matches(msg, overlayKeys.Down):
		if o.focusRow < tvRowCount-1 {
			o.focusRow++
			row := tvRows[o.focusRow]
			if o.focusCol >= len(row.buttons) {
				o.focusCol = len(row.buttons) - 1
			}
		}
	case key.Matches(msg, overlayKeys.Left):
		if o.focusCol > 0 {
			o.focusCol--
		}
	case key.Matches(msg, overlayKeys.Right):
		row := tvRows[o.focusRow]
		if o.focusCol < len(row.buttons)-1 {
			o.focusCol++
		}
	case key.Matches(msg, overlayKeys.Select):
		return o, o.sendCommand(client)
	}
	return o, nil
}

func (o *tvOverlay) sendCommand(client *api.Client) tea.Cmd {
	deviceID := o.item.ID
	row := tvRows[o.focusRow]
	if o.focusCol >= len(row.buttons) {
		return nil
	}
	command := row.buttons[o.focusCol].command

	return func() tea.Msg {
		cmd := api.CommandRequest{
			Command:     command,
			CommandType: "command",
		}
		_, err := client.SendIRCommand(deviceID, cmd)
		return MsgCommandDone{Err: err}
	}
}

func (o *tvOverlay) View() string {
	var sb strings.Builder

	for _, r := range tvRows {
		label := styleOverlayLabel.Render(r.name)
		var btnStr strings.Builder
		for i, btn := range r.buttons {
			var style lipgloss.Style
			if o.focusRow == r.row && o.focusCol == i {
				style = styleButtonFocused
			} else {
				style = styleButton
			}
			btnStr.WriteString(style.Render(btn.label))
		}

		line := fmt.Sprintf("%s %s", label, btnStr.String())
		if o.focusRow == r.row {
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
	sb.WriteString(styleTypeLabel.Render("[↑↓] row  [←→] button  [enter] send  [esc] close"))

	return sb.String()
}

