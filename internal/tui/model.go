package tui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/cache"
	"github.com/yteraoka/sbottui/internal/domain"
)

type state int

const (
	stateLoading state = iota
	stateList
	stateOverlay
	stateError
)

// Model is the root Bubble Tea model.
type Model struct {
	state     state
	client    *api.Client
	list      *listView
	overlay   Overlay
	spinner   spinner.Model
	err       error
	statusMsg string
	statusErr bool
	width     int
	height    int
}

// New creates a new root model.
func New(client *api.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	return Model{
		state:   stateLoading,
		client:  client,
		list:    newListView(),
		spinner: s,
	}
}

// Init starts the spinner and initiates loading.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadData(),
	)
}

func (m Model) loadData() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		var devResp *api.DevicesResponse
		var sceneResp *api.ScenesResponse

		// Try cache first for devices
		if err := cache.Load("devices", &devResp); err != nil {
			var err2 error
			devResp, err2 = client.GetDevices()
			if err2 != nil {
				return MsgLoadError{Err: fmt.Errorf("fetch devices: %w", err2)}
			}
			_ = cache.Save("devices", devResp)
		}

		// Try cache first for scenes
		if err := cache.Load("scenes", &sceneResp); err != nil {
			var err2 error
			sceneResp, err2 = client.GetScenes()
			if err2 != nil {
				return MsgLoadError{Err: fmt.Errorf("fetch scenes: %w", err2)}
			}
			_ = cache.Save("scenes", sceneResp)
		}

		items := domain.BuildList(devResp, sceneResp, domain.SortByName)
		return MsgLoaded{Items: items}
	}
}

func (m Model) refreshData() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		_ = cache.ClearAll()
		devResp, err := client.GetDevices()
		if err != nil {
			return MsgLoadError{Err: fmt.Errorf("fetch devices: %w", err)}
		}
		_ = cache.Save("devices", devResp)

		sceneResp, err := client.GetScenes()
		if err != nil {
			return MsgLoadError{Err: fmt.Errorf("fetch scenes: %w", err)}
		}
		_ = cache.Save("scenes", sceneResp)

		items := domain.BuildList(devResp, sceneResp, domain.SortByName)
		return MsgLoaded{Items: items}
	}
}

// Update handles all messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.setSize(msg.Width, msg.Height)
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case MsgLoaded:
		m.state = stateList
		m.list.setItems(msg.Items)
		m.list.setSortOrder(m.list.sortOrder) // re-apply sort after refresh
		return m, nil

	case MsgLoadError:
		m.state = stateError
		m.err = msg.Err
		return m, nil

	case MsgDeviceStatus:
		if m.state == stateOverlay && msg.Err == nil {
			switch o := m.overlay.(type) {
			case *bulbOverlay:
				o.ApplyStatus(msg.Status)
			case *plugMiniOverlay:
				o.ApplyStatus(msg.Status)
			}
		}
		return m, nil

	case MsgCommandDone:
		if msg.Err != nil {
			m.statusMsg = "Error: " + msg.Err.Error()
			m.statusErr = true
		} else {
			m.statusMsg = "Command sent successfully"
			m.statusErr = false
		}
		return m, nil

	case MsgSceneDone:
		if msg.Err != nil {
			m.statusMsg = "Scene error: " + msg.Err.Error()
			m.statusErr = true
		} else {
			m.statusMsg = fmt.Sprintf("Scene '%s' executed", msg.Name)
			m.statusErr = false
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if key.Matches(msg, listKeys.Quit) {
		return m, tea.Quit
	}

	switch m.state {
	case stateError:
		// r to retry
		if key.Matches(msg, listKeys.Refresh) {
			m.state = stateLoading
			return m, tea.Batch(m.spinner.Tick, m.loadData())
		}

	case stateList:
		return m.handleListKey(msg)

	case stateOverlay:
		return m.handleOverlayKey(msg)
	}

	return m, nil
}

func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, listKeys.Up):
		m.list.moveUp()
	case key.Matches(msg, listKeys.Down):
		m.list.moveDown()
	case key.Matches(msg, listKeys.SortName):
		m.list.setSortOrder(domain.SortByName)
	case key.Matches(msg, listKeys.SortKind):
		m.list.setSortOrder(domain.SortByKind)
	case key.Matches(msg, listKeys.Refresh):
		m.state = stateLoading
		m.statusMsg = ""
		return m, tea.Batch(m.spinner.Tick, m.refreshData())
	case key.Matches(msg, listKeys.Select):
		return m.selectItem()
	}
	return m, nil
}

func (m Model) selectItem() (tea.Model, tea.Cmd) {
	item := m.list.selected()
	if item == nil {
		return m, nil
	}

	// Scenes are executed immediately
	if item.Kind == domain.KindScene {
		client := m.client
		sceneID := item.ID
		sceneName := item.Name
		return m, func() tea.Msg {
			_, err := client.ExecuteScene(sceneID)
			return MsgSceneDone{Name: sceneName, Err: err}
		}
	}

	overlay, cmd := NewOverlay(*item, m.client)
	m.overlay = overlay
	m.state = stateOverlay
	return m, cmd
}

func (m Model) handleOverlayKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, overlayKeys.Close) {
		m.state = stateList
		m.overlay = nil
		return m, nil
	}

	if m.overlay != nil {
		var cmd tea.Cmd
		m.overlay, cmd = m.overlay.Update(msg, m.client)
		return m, cmd
	}
	return m, nil
}

// View renders the full UI.
func (m Model) View() tea.View {
	var v tea.View
	v.AltScreen = true
	switch m.state {
	case stateLoading:
		v.SetContent(fmt.Sprintf("\n  %s Loading devices and scenes...\n", m.spinner.View()))
	case stateError:
		v.SetContent(fmt.Sprintf("\n  Error: %v\n\n  Press r to retry or q to quit.\n", m.err))
	case stateList:
		v.SetContent(m.viewList())
	case stateOverlay:
		v.SetContent(m.viewOverlay())
	}
	return v
}

func (m Model) viewList() string {
	listContent := m.list.view()
	statusBar := m.renderStatusBar()
	return listContent + "\n" + statusBar
}

func (m Model) viewOverlay() string {
	if m.overlay == nil {
		return m.viewList()
	}

	// Render the overlay centered on top of the list
	overlayContent := m.overlay.View()
	title := m.overlay.Title()

	content := styleOverlayTitle.Render(title) + "\n" + overlayContent
	box := styleOverlayBorder.Render(content)

	// Center the box
	boxWidth := lipgloss.Width(box)
	boxHeight := lipgloss.Height(box)

	leftPad := (m.width - boxWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	topPad := (m.height - boxHeight - 2) / 2
	if topPad < 0 {
		topPad = 0
	}

	padStyle := lipgloss.NewStyle().MarginLeft(leftPad).MarginTop(topPad)
	return padStyle.Render(box) + "\n" + m.renderStatusBar()
}

func (m Model) renderStatusBar() string {
	if m.statusMsg == "" {
		return styleStatusBar.Width(m.width).Render("Ready")
	}
	var msg string
	if m.statusErr {
		msg = styleStatusErr.Render(m.statusMsg)
	} else {
		msg = styleStatusOk.Render(m.statusMsg)
	}
	return styleStatusBar.Width(m.width).Render(msg)
}
