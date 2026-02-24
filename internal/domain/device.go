package domain

import (
	"sort"
	"strings"

	"github.com/yteraoka/sbottui/internal/api"
)

// Kind classifies a list item.
type Kind int

const (
	KindPhysical Kind = iota
	KindIR
	KindScene
)

// SortOrder determines how items are sorted in the list.
type SortOrder int

const (
	SortByName SortOrder = iota
	SortByKind
)

// ListItem is a unified item that can represent a physical device, IR device, or scene.
type ListItem struct {
	ID         string
	Name       string
	Kind       Kind
	DeviceType string // physical device type (e.g. "Color Bulb")
	RemoteType string // IR remote type (e.g. "Air Conditioner", "TV")
}

// TypeLabel returns a human-readable label for the item type.
func (i ListItem) TypeLabel() string {
	switch i.Kind {
	case KindPhysical:
		return i.DeviceType
	case KindIR:
		return "IR:" + i.RemoteType
	case KindScene:
		return "Scene"
	}
	return ""
}

// BuildList converts API responses into sorted ListItems.
func BuildList(devices *api.DevicesResponse, scenes *api.ScenesResponse, order SortOrder) []ListItem {
	var items []ListItem

	if devices != nil {
		for _, d := range devices.Body.DeviceList {
			items = append(items, ListItem{
				ID:         d.DeviceID,
				Name:       d.DeviceName,
				Kind:       KindPhysical,
				DeviceType: d.DeviceType,
			})
		}
		for _, d := range devices.Body.InfraredRemoteList {
			items = append(items, ListItem{
				ID:         d.DeviceID,
				Name:       d.DeviceName,
				Kind:       KindIR,
				RemoteType: d.RemoteType,
			})
		}
	}

	if scenes != nil {
		for _, s := range scenes.Body {
			items = append(items, ListItem{
				ID:   s.SceneID,
				Name: s.SceneName,
				Kind: KindScene,
			})
		}
	}

	Sort(items, order)
	return items
}

// Sort sorts items in-place according to order.
func Sort(items []ListItem, order SortOrder) {
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		switch order {
		case SortByKind:
			if a.Kind != b.Kind {
				return a.Kind < b.Kind
			}
			typeA := typeKey(a)
			typeB := typeKey(b)
			if typeA != typeB {
				return typeA < typeB
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		default: // SortByName
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		}
	})
}

func typeKey(i ListItem) string {
	switch i.Kind {
	case KindPhysical:
		return i.DeviceType
	case KindIR:
		return i.RemoteType
	case KindScene:
		return "~scene"
	}
	return ""
}
