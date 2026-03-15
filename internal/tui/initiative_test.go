package tui

import (
	"context"
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func newTestInitiativeTab(fc *fakeClient) *initiativeTab {
	ctx := context.Background()
	tab := NewInitiativeTab(ctx, fc)
	tab.size = tea.WindowSizeMsg{Width: 120, Height: 20}
	return tab
}

func updateInitiativeTab(t *testing.T, tab *initiativeTab, msg tea.Msg) (*initiativeTab, tea.Cmd) {
	t.Helper()
	m, cmd := tab.Update(msg)
	result, ok := m.(*initiativeTab)
	require.True(t, ok)
	return result, cmd
}

func TestInitiativeTab_Loading(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	cmd := tab.Init()
	assert.NotNil(t, cmd)
	assert.Equal(t, stateLoading, tab.state)
	assert.Contains(t, tab.View(), "Loading initiatives")
}

func TestInitiativeTab_ListRendering(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Initiatives: []api.Initiative{
			{ID: "1", Name: "Q1 Goals", Status: "Active"},
			{ID: "2", Name: "Q2 Planning", Status: "Planned"},
		},
	})

	assert.Equal(t, stateList, tab.state)
	view := tab.View()
	assert.Contains(t, view, "Q1 Goals")
	assert.Contains(t, view, "Active")
	assert.Contains(t, view, "Q2 Planning")
	assert.Contains(t, view, "Planned")
}

func TestInitiativeTab_CursorNavigation(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Initiatives: []api.Initiative{
			{ID: "1", Name: "Initiative A"},
			{ID: "2", Name: "Initiative B"},
			{ID: "3", Name: "Initiative C"},
		},
	})

	assert.Equal(t, 0, tab.cursor)

	// Move down
	tab, _ = updateInitiativeTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Move down again
	tab, _ = updateInitiativeTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move down at bottom — should stay
	tab, _ = updateInitiativeTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move up
	tab, _ = updateInitiativeTab(t, tab, keyMsg("k"))
	assert.Equal(t, 1, tab.cursor)

	// View should highlight current row (accent bar marks selection)
	view := tab.View()
	assert.Contains(t, view, "▎")
}

func TestInitiativeTab_DetailView(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Initiatives: []api.Initiative{
			{ID: "init-1", Name: "Q1 Goals"},
		},
	})

	// Press Enter to load detail
	tab, cmd := updateInitiativeTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, cmd)

	// Simulate detail loaded
	tab, cmd = updateInitiativeTab(t, tab, InitiativeDetailLoadedMsg{
		Initiative: &api.Initiative{
			ID:          "init-1",
			Name:        "Q1 Goals",
			Description: "Goals for the first quarter",
			Status:      "Active",
			TargetDate:  "2026-03-31",
			Owner:       &api.User{DisplayName: "Bob"},
			URL:         "https://linear.app/initiative/q1-goals",
		},
	})

	assert.Equal(t, stateDetail, tab.state)
	// Detail loaded should also fire loadInitiativeProjects
	require.NotNil(t, cmd)

	// Simulate projects loaded
	tab, _ = updateInitiativeTab(t, tab, InitiativeProjectsLoadedMsg{
		Projects: []api.Project{
			{ID: "p1", Name: "Alpha Release", Status: api.ProjectStatus{Name: "Started"}},
			{ID: "p2", Name: "Beta Launch", Status: api.ProjectStatus{Name: "Planned"}},
		},
	})

	view := tab.View()
	assert.Contains(t, view, "Q1 Goals")
	assert.Contains(t, view, "Goals for the first quarter")
	assert.Contains(t, view, "Active")
	assert.Contains(t, view, "Bob")
	assert.Contains(t, view, "Linked Projects")
	assert.Contains(t, view, "Alpha Release")
	assert.Contains(t, view, "Beta Launch")
}

func TestInitiativeTab_BackToList(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Initiatives: []api.Initiative{
			{ID: "1", Name: "Initiative A"},
			{ID: "2", Name: "Initiative B"},
		},
	})

	// Move cursor to second item
	tab, _ = updateInitiativeTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Go to detail
	tab, _ = updateInitiativeTab(t, tab, InitiativeDetailLoadedMsg{
		Initiative: &api.Initiative{ID: "2", Name: "Initiative B", Status: "Planned"},
	})
	assert.Equal(t, stateDetail, tab.state)

	// Press Esc to go back
	tab, _ = updateInitiativeTab(t, tab, tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, stateList, tab.state)
	assert.Equal(t, 1, tab.cursor, "cursor should be preserved")
}

func TestInitiativeTab_Pagination(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	// First page
	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Initiatives: []api.Initiative{
			{ID: "1", Name: "Initiative A"},
		},
		PageInfo: api.PageInfo{HasNextPage: true, EndCursor: "cursor-1"},
	})

	assert.Len(t, tab.items, 1)
	assert.True(t, tab.pageInfo.HasNextPage)

	// Press n to load next page
	tab, cmd := updateInitiativeTab(t, tab, keyMsg("n"))
	require.NotNil(t, cmd)

	// Simulate next page loaded
	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Initiatives: []api.Initiative{
			{ID: "2", Name: "Initiative B"},
		},
		PageInfo: api.PageInfo{HasNextPage: false},
	})

	assert.Len(t, tab.items, 2, "items should accumulate")
	assert.False(t, tab.pageInfo.HasNextPage)
}

func TestInitiativeTab_Error(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Err: fmt.Errorf("network error"),
	})

	assert.Equal(t, stateError, tab.state)
	view := tab.View()
	assert.Contains(t, view, "network error")
	assert.Contains(t, view, "retry")

	// Press r to retry
	tab, cmd := updateInitiativeTab(t, tab, keyMsg("r"))
	assert.Equal(t, stateLoading, tab.state)
	assert.NotNil(t, cmd)
}

func TestInitiativeTab_EmptyList(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestInitiativeTab(fc)

	tab, _ = updateInitiativeTab(t, tab, InitiativesLoadedMsg{
		Initiatives: []api.Initiative{},
	})

	assert.Equal(t, stateList, tab.state)
	assert.Contains(t, tab.View(), "No initiatives found")
}
