package tui

import (
	"context"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// updateApp sends a message to the app and returns the updated AppModel.
func updateApp(t *testing.T, app AppModel, msg tea.Msg) (AppModel, tea.Cmd) {
	t.Helper()
	m, cmd := app.Update(msg)
	result, ok := m.(AppModel)
	require.True(t, ok, "Update should return AppModel")
	return result, cmd
}

func keyMsg(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func newTestApp(t *testing.T) AppModel {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	tabs := []TabModel{
		newPlaceholderTab("Issues"),
		newPlaceholderTab("Cycles"),
		newPlaceholderTab("Projects"),
		newPlaceholderTab("Initiatives"),
	}
	tabNames := []string{"Issues", "Cycles", "Projects", "Initiatives"}
	app := NewApp(ctx, cancel, tabs, tabNames, "TestTeam")
	_ = app.Init()
	app, _ = updateApp(t, app, tea.WindowSizeMsg{Width: 80, Height: 24})
	return app
}

func TestAppModel_InitialTab(t *testing.T) {
	app := newTestApp(t)
	assert.Equal(t, 0, app.activeTab)
}

func TestAppModel_TabSwitchByNumber(t *testing.T) {
	app := newTestApp(t)

	tests := []struct {
		key      string
		expected int
	}{
		{"2", 1},
		{"3", 2},
		{"4", 3},
		{"1", 0},
	}

	for _, tt := range tests {
		app, _ = updateApp(t, app, keyMsg(tt.key))
		assert.Equal(t, tt.expected, app.activeTab, "pressing %s should switch to tab %d", tt.key, tt.expected)
	}
}

func TestAppModel_TabSwitchByTab(t *testing.T) {
	app := newTestApp(t)

	app, _ = updateApp(t, app, tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, 1, app.activeTab)

	app, _ = updateApp(t, app, tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, 2, app.activeTab)
}

func TestAppModel_TabSwitchByShiftTab(t *testing.T) {
	app := newTestApp(t)

	app, _ = updateApp(t, app, tea.KeyMsg{Type: tea.KeyShiftTab})
	assert.Equal(t, 3, app.activeTab)

	app, _ = updateApp(t, app, tea.KeyMsg{Type: tea.KeyShiftTab})
	assert.Equal(t, 2, app.activeTab)
}

func TestAppModel_TabWrapAround(t *testing.T) {
	app := newTestApp(t)

	app, _ = updateApp(t, app, keyMsg("4"))
	assert.Equal(t, 3, app.activeTab)

	app, _ = updateApp(t, app, tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, 0, app.activeTab)
}

func TestAppModel_WindowSizeBroadcast(t *testing.T) {
	app := newTestApp(t)

	for i, tab := range app.tabs {
		pt, ok := tab.(*placeholderTab)
		require.True(t, ok)
		assert.Equal(t, 80, pt.size.Width, "tab %d should have width 80", i)
	}
}

func TestAppModel_QuitReturnsQuit(t *testing.T) {
	app := newTestApp(t)

	_, cmd := updateApp(t, app, keyMsg("q"))
	require.NotNil(t, cmd)
	msg := cmd()
	_, isQuit := msg.(tea.QuitMsg)
	assert.True(t, isQuit, "q should return tea.Quit")
}

func TestAppModel_ViewContainsTabBar(t *testing.T) {
	app := newTestApp(t)
	view := app.View()

	assert.Contains(t, view, "Issues")
	assert.Contains(t, view, "Cycles")
	assert.Contains(t, view, "Projects")
	assert.Contains(t, view, "Initiatives")
}

func TestAppModel_ViewContainsTeamName(t *testing.T) {
	app := newTestApp(t)
	assert.Contains(t, app.View(), "TestTeam")
}

func TestAppModel_ViewContainsStatusBar(t *testing.T) {
	app := newTestApp(t)
	assert.Contains(t, app.View(), "quit")
}

func TestAppModel_HelpToggle(t *testing.T) {
	app := newTestApp(t)

	app, _ = updateApp(t, app, keyMsg("?"))
	assert.True(t, app.showHelp)
	assert.Contains(t, app.View(), "Keybindings")

	app, _ = updateApp(t, app, keyMsg("?"))
	assert.False(t, app.showHelp)
}

func TestAppModel_HelpDismissWithEsc(t *testing.T) {
	app := newTestApp(t)

	app, _ = updateApp(t, app, keyMsg("?"))
	assert.True(t, app.showHelp)

	app, _ = updateApp(t, app, tea.KeyMsg{Type: tea.KeyEscape})
	assert.False(t, app.showHelp)
}

func TestAppModel_HelpSuppressesNavigation(t *testing.T) {
	app := newTestApp(t)

	app, _ = updateApp(t, app, keyMsg("?"))

	app, _ = updateApp(t, app, keyMsg("2"))
	assert.Equal(t, 0, app.activeTab, "tab switch should be suppressed while help is shown")
	assert.True(t, app.showHelp, "help should still be shown")
}

func TestAppModel_QuitFromHelp(t *testing.T) {
	app := newTestApp(t)

	app, _ = updateApp(t, app, keyMsg("?"))

	_, cmd := updateApp(t, app, keyMsg("q"))
	require.NotNil(t, cmd)
	msg := cmd()
	_, isQuit := msg.(tea.QuitMsg)
	assert.True(t, isQuit, "q should quit even from help overlay")
}

func TestAppModel_LazyTabInit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var initCalls [4]int
	tabs := make([]TabModel, 4)
	for i := range tabs {
		tabs[i] = &initTrackingTab{idx: i, calls: &initCalls}
	}
	tabNames := []string{"A", "B", "C", "D"}
	app := NewApp(ctx, cancel, tabs, tabNames, "Team")
	cmd := app.Init()
	if cmd != nil {
		_ = cmd()
	}

	assert.Equal(t, 1, initCalls[0], "tab 0 should be initialized")
	assert.Equal(t, 0, initCalls[1], "tab 1 should not be initialized")

	app, _ = updateApp(t, app, tea.WindowSizeMsg{Width: 80, Height: 24})

	var c tea.Cmd
	app, c = updateApp(t, app, keyMsg("2"))
	if c != nil {
		_ = c()
	}
	assert.Equal(t, 1, initCalls[1], "tab 1 should now be initialized")

	_, c = updateApp(t, app, keyMsg("1"))
	if c != nil {
		_ = c()
	}
	assert.Equal(t, 1, initCalls[0], "tab 0 should not be re-initialized")
}

func TestAppModel_EmptyView(t *testing.T) {
	app := NewApp(context.Background(), func() {}, nil, nil, "")
	assert.Equal(t, "", app.View())
}

func TestAppModel_HelpOverlayContent(t *testing.T) {
	app := newTestApp(t)

	// Use taller terminal so all help content fits.
	app, _ = updateApp(t, app, tea.WindowSizeMsg{Width: 80, Height: 40})

	app, _ = updateApp(t, app, keyMsg("?"))
	view := app.View()

	assert.Contains(t, view, "Keybindings")
	assert.Contains(t, view, "j/down")
	assert.Contains(t, view, "k/up")
	assert.Contains(t, view, "Press ? or Esc to close")
}

// initTrackingTab tracks how many times Init was called.
type initTrackingTab struct {
	idx   int
	calls *[4]int
	size  tea.WindowSizeMsg
}

func (t *initTrackingTab) Init() tea.Cmd {
	t.calls[t.idx]++
	return nil
}

func (t *initTrackingTab) Update(msg tea.Msg) (TabModel, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		t.size = msg
	}
	return t, nil
}

func (t *initTrackingTab) View() string {
	return ""
}

func (t *initTrackingTab) ShortHelp() []key.Binding {
	return nil
}
