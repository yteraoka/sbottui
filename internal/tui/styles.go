package tui

import "charm.land/lipgloss/v2"

var (
	// List styles
	styleSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("236"))

	styleNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	styleTypeLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)

	styleStatusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Background(lipgloss.Color("234")).
			Padding(0, 1)

	styleStatusOk = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	styleStatusErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	// Overlay styles
	styleOverlayBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("212")).
				Padding(1, 2)

	styleOverlayTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")).
				MarginBottom(1)

	styleOverlayFocused = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")).
				Background(lipgloss.Color("236"))

	styleOverlayNormal = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	styleOverlayLabel = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Width(20)
	styleButton = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("238")).
			Padding(0, 1).
			MarginRight(1)

	styleButtonFocused = lipgloss.NewStyle().
				Foreground(lipgloss.Color("231")).
				Background(lipgloss.Color("212")).
				Padding(0, 1).
				MarginRight(1)

	styleSortIndicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")).
				Bold(true)

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			MarginBottom(1)
)
