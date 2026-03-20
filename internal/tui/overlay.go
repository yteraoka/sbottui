package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/domain"
)

// Overlay is implemented by all overlay types.
type Overlay interface {
	// Update handles key events while overlay is active.
	Update(msg tea.KeyMsg, client *api.Client) (Overlay, tea.Cmd)
	// View renders the overlay content (without border).
	View() string
	// Title returns the overlay title.
	Title() string
}

// NewOverlay creates the appropriate overlay for the given list item.
// Returns nil if no overlay is needed (e.g., for scenes, handled separately).
func NewOverlay(item domain.ListItem, client *api.Client) (Overlay, tea.Cmd) {
	switch item.Kind {
	case domain.KindPhysical:
		switch item.DeviceType {
		case "Color Bulb", "Strip Light", "Color Bulb 1M", "Color Bulb 2M":
			return newBulbOverlay(item, client)
		case "Plug Mini (US)", "Plug Mini (JP)":
			return newPlugMiniOverlay(item, client)
		}
	case domain.KindIR:
		switch item.RemoteType {
		case "Air Conditioner":
			return newACOverlay(item), nil
		case "TV":
			return newTVOverlay(item), nil
		case "Light", "Fan", "DIY Fan", "DIY Light":
			return newIRLightOverlay(item), nil
		}
	}
	return newGenericOverlay(item), nil
}

// genericOverlay is a fallback overlay for unsupported device types.
type genericOverlay struct {
	item domain.ListItem
}

func newGenericOverlay(item domain.ListItem) *genericOverlay {
	return &genericOverlay{item: item}
}

func (o *genericOverlay) Title() string {
	return o.item.Name
}

func (o *genericOverlay) Update(msg tea.KeyMsg, _ *api.Client) (Overlay, tea.Cmd) {
	return o, nil
}

func (o *genericOverlay) View() string {
	return styleOverlayNormal.Render("No controls available for " + o.item.TypeLabel())
}
