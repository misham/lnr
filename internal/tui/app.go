package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppModel is the top-level Bubble Tea model for the TUI dashboard.
type AppModel struct {
	tabs       []TabModel
	tabNames   []string
	activeTab  int
	keys       KeyMap
	size       tea.WindowSizeMsg
	teamName   string
	ctx        context.Context
	cancel     context.CancelFunc
	loaded     [4]bool
	showHelp   bool
	prefetched bool
}

// NewApp creates a new AppModel.
func NewApp(ctx context.Context, cancel context.CancelFunc, tabs []TabModel, tabNames []string, teamName string) AppModel {
	m := AppModel{
		tabs:     tabs,
		tabNames: tabNames,
		keys:     DefaultKeyMap(),
		teamName: teamName,
		ctx:      ctx,
		cancel:   cancel,
	}
	if len(tabs) > 0 {
		m.loaded[0] = true
	}
	return m
}

func (m AppModel) Init() tea.Cmd {
	if len(m.tabs) == 0 {
		return nil
	}
	return m.tabs[m.activeTab].Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
		contentHeight := m.contentHeight()
		contentSize := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: contentHeight,
		}
		var cmds []tea.Cmd
		for i, tab := range m.tabs {
			newTab, cmd := tab.Update(contentSize)
			m.tabs[i] = newTab
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		if m.showHelp {
			return m.handleHelpKeys(msg)
		}
		return m.handleKeys(msg)
	}

	// Broadcast data-loaded messages to ALL tabs so prefetched data reaches the right tab.
	if len(m.tabs) > 0 {
		var cmds []tea.Cmd
		if isDataLoadedMsg(msg) {
			for i, tab := range m.tabs {
				newTab, cmd := tab.Update(msg)
				m.tabs[i] = newTab
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		} else {
			newTab, cmd := m.tabs[m.activeTab].Update(msg)
			m.tabs[m.activeTab] = newTab
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

		// Trigger prefetch after first data load.
		if !m.prefetched && isDataLoadedMsg(msg) {
			m.prefetched = true
			cmds = append(cmds, m.prefetchTabs()...)
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m AppModel) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Help), key.Matches(msg, m.keys.Back):
		m.showHelp = false
		return m, nil
	case key.Matches(msg, m.keys.Quit):
		m.cancel()
		return m, tea.Quit
	}
	return m, nil
}

func (m AppModel) handleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		m.cancel()
		return m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.showHelp = true
		return m, nil
	case key.Matches(msg, m.keys.Tab1):
		return m.switchTab(0)
	case key.Matches(msg, m.keys.Tab2):
		return m.switchTab(1)
	case key.Matches(msg, m.keys.Tab3):
		return m.switchTab(2)
	case key.Matches(msg, m.keys.Tab4):
		return m.switchTab(3)
	case key.Matches(msg, m.keys.Tab):
		next := (m.activeTab + 1) % len(m.tabs)
		return m.switchTab(next)
	case key.Matches(msg, m.keys.ShiftTab):
		prev := (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		return m.switchTab(prev)
	default:
		if len(m.tabs) > 0 {
			newTab, cmd := m.tabs[m.activeTab].Update(msg)
			m.tabs[m.activeTab] = newTab
			return m, cmd
		}
	}
	return m, nil
}

func (m AppModel) switchTab(idx int) (tea.Model, tea.Cmd) {
	if idx < 0 || idx >= len(m.tabs) {
		return m, nil
	}
	m.activeTab = idx
	if !m.loaded[idx] {
		m.loaded[idx] = true
		return m, m.tabs[idx].Init()
	}
	return m, nil
}

func (m AppModel) View() string {
	if m.size.Width == 0 {
		return ""
	}

	tabBar := m.renderTabBar()
	var content string
	if len(m.tabs) > 0 {
		content = m.tabs[m.activeTab].View()
	}

	var helpKeys []key.Binding
	if len(m.tabs) > 0 {
		helpKeys = m.tabs[m.activeTab].ShortHelp()
	}
	helpKeys = append(helpKeys, m.keys.Help, m.keys.Quit)
	info := fmt.Sprintf("%d/%d", m.activeTab+1, len(m.tabs))
	statusBar := renderStatusBar(m.size.Width, helpKeys, info)

	if m.showHelp {
		content = m.renderHelpOverlay()
	}

	// Force content area to exact height so tab bar and status bar are always visible.
	contentArea := lipgloss.NewStyle().
		Width(m.size.Width).
		Height(m.contentHeight()).
		MaxHeight(m.contentHeight()).
		Render(content)

	return tabBar + "\n" + contentArea + "\n" + statusBar
}

func (m AppModel) renderTabBar() string {
	// Tab names row.
	var tabs []string
	var activeStart, activeWidth int
	pos := 0
	for i, name := range m.tabNames {
		label := name
		if i == m.activeTab {
			rendered := activeTabStyle.Render(label)
			activeStart = pos
			activeWidth = lipgloss.Width(rendered)
			tabs = append(tabs, rendered)
		} else {
			rendered := inactiveTabStyle.Render(label)
			tabs = append(tabs, rendered)
		}
		pos += lipgloss.Width(tabs[len(tabs)-1])
	}
	left := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	right := lipgloss.NewStyle().Foreground(mutedColor).Render(m.teamName)
	gap := max(m.size.Width-lipgloss.Width(left)-lipgloss.Width(right), 1)
	nameRow := left + strings.Repeat(" ", gap) + right

	// Underline row: thick purple under active tab, thin line elsewhere.
	var underline strings.Builder
	for i := 0; i < m.size.Width; i++ {
		if i >= activeStart && i < activeStart+activeWidth {
			underline.WriteString(lipgloss.NewStyle().Foreground(purpleColor).Render("━"))
		} else {
			underline.WriteString(lipgloss.NewStyle().Foreground(dimColor).Render("─"))
		}
	}

	return nameRow + "\n" + underline.String()
}

func (m AppModel) renderHelpOverlay() string {
	keys := m.keys
	bindings := []key.Binding{
		keys.Up, keys.Down, keys.Enter, keys.Back,
		keys.Tab, keys.ShiftTab,
		keys.Tab1, keys.Tab2, keys.Tab3, keys.Tab4,
		keys.Search, keys.NextPage, keys.PrevPage,
		keys.Top, keys.Bottom,
		keys.HalfDown, keys.HalfUp,
		keys.Retry, keys.Help, keys.Quit,
	}

	const colW = 16
	var rows []string
	for _, b := range bindings {
		k := fmt.Sprintf("%*s", colW, b.Help().Key)
		d := fmt.Sprintf("   %-*s", colW, b.Help().Desc)
		rows = append(rows, k+d)
	}
	tableWidth := colW*2 + 3
	title := lipgloss.NewStyle().Bold(true).Foreground(purpleColor).Width(tableWidth).Align(lipgloss.Center).Render("Keybindings")
	footer := lipgloss.NewStyle().Foreground(mutedColor).Width(tableWidth).Align(lipgloss.Center).Render("Press ? or Esc to close")
	block := title + "\n\n" + strings.Join(rows, "\n") + "\n\n" + footer

	return lipgloss.Place(m.size.Width, m.contentHeight(), lipgloss.Center, lipgloss.Center, block)
}

func isDataLoadedMsg(msg tea.Msg) bool {
	switch msg.(type) {
	case IssuesLoadedMsg, CyclesLoadedMsg, ProjectsLoadedMsg, InitiativesLoadedMsg,
		CycleIssuesLoadedMsg, ProjectIssuesLoadedMsg, InitiativeProjectsLoadedMsg:
		return true
	}
	return false
}

func (m AppModel) prefetchTabs() []tea.Cmd {
	var cmds []tea.Cmd
	for i, tab := range m.tabs {
		if i != m.activeTab && !m.loaded[i] {
			cmds = append(cmds, tab.Init())
			m.loaded[i] = true
		}
	}
	return cmds
}

func (m AppModel) contentHeight() int {
	// Tab bar: 2 lines (names + underline), status bar: 1 line.
	return max(m.size.Height-3, 1)
}
