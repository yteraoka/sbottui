package tui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/yteraoka/sbottui/internal/domain"
)

// listView renders the device/scene list.
type listView struct {
	items     []domain.ListItem
	cursor    int
	sortOrder domain.SortOrder
	width     int
	height    int
}

func newListView() *listView {
	return &listView{}
}

func (l *listView) setItems(items []domain.ListItem) {
	l.items = items
	if l.cursor >= len(items) && len(items) > 0 {
		l.cursor = len(items) - 1
	}
}

func (l *listView) selected() *domain.ListItem {
	if len(l.items) == 0 || l.cursor < 0 || l.cursor >= len(l.items) {
		return nil
	}
	item := l.items[l.cursor]
	return &item
}

func (l *listView) moveUp() {
	if l.cursor > 0 {
		l.cursor--
	}
}

func (l *listView) moveDown() {
	if l.cursor < len(l.items)-1 {
		l.cursor++
	}
}

func (l *listView) setSortOrder(order domain.SortOrder) {
	if l.sortOrder == order {
		return
	}
	selectedID := ""
	if item := l.selected(); item != nil {
		selectedID = item.ID
	}
	l.sortOrder = order
	domain.Sort(l.items, order)
	// Re-find cursor position after sort
	for i, item := range l.items {
		if item.ID == selectedID {
			l.cursor = i
			break
		}
	}
}

func (l *listView) view() string {
	if len(l.items) == 0 {
		return styleNormal.Render("No devices found.")
	}

	// Calculate visible range
	listHeight := l.height - 4 // header + status bar
	if listHeight < 1 {
		listHeight = 1
	}

	start := 0
	if l.cursor >= listHeight {
		start = l.cursor - listHeight + 1
	}
	end := start + listHeight
	if end > len(l.items) {
		end = len(l.items)
	}

	var sb strings.Builder

	// Header
	sortLabel := "Name"
	if l.sortOrder == domain.SortByKind {
		sortLabel = "Kind"
	}
	header := fmt.Sprintf(" sbottui  Sort: %s  [←] Name  [→] Kind  [r] Refresh  [q] Quit",
		styleSortIndicator.Render(sortLabel))
	sb.WriteString(styleHeader.Width(l.width).Render(header))
	sb.WriteString("\n")

	// Items
	nameWidth := l.width - 25
	if nameWidth < 10 {
		nameWidth = 10
	}

	for i := start; i < end; i++ {
		item := l.items[i]
		name := padRight(truncate(item.Name, nameWidth), nameWidth)
		typeLabel := styleTypeLabel.Render(padRight(truncate(item.TypeLabel(), 20), 20))

		line := " " + name + "  " + typeLabel

		if i == l.cursor {
			sb.WriteString(styleSelected.Width(l.width).Render(line))
		} else {
			sb.WriteString(styleNormal.Width(l.width).Render(line))
		}
		if i < end-1 {
			sb.WriteString("\n")
		}
	}

	// Scroll indicator
	if len(l.items) > listHeight {
		pct := 0
		if len(l.items) > 1 {
			pct = (l.cursor * 100) / (len(l.items) - 1)
		}
		sb.WriteString("\n")
		sb.WriteString(styleTypeLabel.Render(fmt.Sprintf(" %d/%d (%d%%)", l.cursor+1, len(l.items), pct)))
	}

	return sb.String()
}

func (l *listView) setSize(w, h int) {
	l.width = w
	l.height = h
}

// truncate cuts s to fit within max display columns, appending "…" if truncated.
func truncate(s string, max int) string {
	if runewidth.StringWidth(s) <= max {
		return s
	}
	if max <= 1 {
		return runewidth.Truncate(s, max, "")
	}
	return runewidth.Truncate(s, max-1, "…")
}

// padRight pads s with spaces on the right to reach exactly width display columns.
func padRight(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}
