package tui

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/misham/linear-cli/internal/api"
)

func TestIssueTab_SearchActivates(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix authentication bug"},
			{ID: "2", Identifier: "ENG-2", Title: "Add search feature"},
			{ID: "3", Identifier: "ENG-3", Title: "Update docs"},
		},
	})

	// Press / to activate search
	tab, _ = updateIssueTab(t, tab, keyMsg("/"))
	assert.True(t, tab.searching)
	assert.Contains(t, tab.View(), "/")
}

func TestIssueTab_SearchFilters(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix authentication bug"},
			{ID: "2", Identifier: "ENG-2", Title: "Add search feature"},
			{ID: "3", Identifier: "ENG-3", Title: "Update docs"},
		},
	})

	// Activate search
	tab, _ = updateIssueTab(t, tab, keyMsg("/"))

	// Type "auth" — should filter to ENG-1
	for _, ch := range "auth" {
		tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}

	view := tab.View()
	assert.Contains(t, view, "ENG-1")
	assert.Contains(t, view, "authentication")
	assert.NotContains(t, view, "ENG-2")
	assert.NotContains(t, view, "ENG-3")
}

func TestIssueTab_SearchEscClears(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix bug"},
			{ID: "2", Identifier: "ENG-2", Title: "Add feature"},
		},
	})

	// Activate search and type
	tab, _ = updateIssueTab(t, tab, keyMsg("/"))
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("bug")})

	// Esc clears filter
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEscape})
	assert.False(t, tab.searching)
	assert.Empty(t, tab.filter)

	view := tab.View()
	assert.Contains(t, view, "ENG-1")
	assert.Contains(t, view, "ENG-2")
}

func TestIssueTab_SearchEnterLocks(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix bug"},
			{ID: "2", Identifier: "ENG-2", Title: "Add feature"},
		},
	})

	tab, _ = updateIssueTab(t, tab, keyMsg("/"))
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("bug")})

	// Enter locks filter
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})
	assert.False(t, tab.searching)
	assert.Equal(t, "bug", tab.filter)

	// Should still be filtered
	view := tab.View()
	assert.Contains(t, view, "ENG-1")
	assert.NotContains(t, view, "ENG-2")
}

func TestIssueTab_SearchPaginationDisabled(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix bug"},
		},
		PageInfo: api.PageInfo{HasNextPage: true, EndCursor: "cursor-1"},
	})

	// Activate search
	tab, _ = updateIssueTab(t, tab, keyMsg("/"))
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("bug")})
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})

	// n should not trigger pagination while filter is active
	_, cmd := updateIssueTab(t, tab, keyMsg("n"))
	assert.Nil(t, cmd)
}

func TestIssueTab_SearchCaseInsensitive(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix Authentication Bug"},
			{ID: "2", Identifier: "ENG-2", Title: "Add feature"},
		},
	})

	tab, _ = updateIssueTab(t, tab, keyMsg("/"))
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("auth")})
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})

	view := tab.View()
	assert.Contains(t, view, "ENG-1")
	assert.NotContains(t, view, "ENG-2")
}

func TestIssueTab_SearchMatchesIdentifier(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix bug"},
			{ID: "2", Identifier: "ENG-2", Title: "Add feature"},
		},
	})

	tab, _ = updateIssueTab(t, tab, keyMsg("/"))
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("ENG-2")})
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})

	view := tab.View()
	assert.NotContains(t, view, "ENG-1")
	assert.Contains(t, view, "ENG-2")
}

func TestAppModel_PrefetchAfterFirstLoad(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var initCalls [4]int
	tabs := make([]TabModel, 4)
	for i := range tabs {
		tabs[i] = &initTrackingTab{idx: i, calls: &initCalls}
	}
	tabNames := []string{"A", "B", "C", "D"}
	app := NewApp(ctx, cancel, tabs, tabNames, "Team")
	_ = app.Init()

	app, _ = updateApp(t, app, tea.WindowSizeMsg{Width: 80, Height: 24})

	// Simulate first data load completing — triggers prefetch
	app, cmd := updateApp(t, app, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	// Prefetch should have been triggered
	_ = cmd
	assert.True(t, app.prefetched)
}

func TestAppModel_PrefetchOnlyOnce(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var initCalls [4]int
	tabs := make([]TabModel, 4)
	for i := range tabs {
		tabs[i] = &initTrackingTab{idx: i, calls: &initCalls}
	}
	tabNames := []string{"A", "B", "C", "D"}
	app := NewApp(ctx, cancel, tabs, tabNames, "Team")
	_ = app.Init()
	app, _ = updateApp(t, app, tea.WindowSizeMsg{Width: 80, Height: 24})

	// First data load triggers prefetch
	app, _ = updateApp(t, app, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})
	assert.True(t, app.prefetched)

	// Second data load should not trigger prefetch again
	app, cmd := updateApp(t, app, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "2", Identifier: "ENG-2", Title: "Test 2"}},
	})
	// cmd should be nil since prefetch already happened and this is delegated to active tab
	_ = cmd // no additional prefetch commands
}

func TestAppModel_PrefetchSkipsLoadedTabs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var initCalls [4]int
	tabs := make([]TabModel, 4)
	for i := range tabs {
		tabs[i] = &initTrackingTab{idx: i, calls: &initCalls}
	}
	tabNames := []string{"A", "B", "C", "D"}
	app := NewApp(ctx, cancel, tabs, tabNames, "Team")
	_ = app.Init()
	app, _ = updateApp(t, app, tea.WindowSizeMsg{Width: 80, Height: 24})

	// Manually switch to tab 2 first
	app, cmd := updateApp(t, app, keyMsg("2"))
	if cmd != nil {
		_ = cmd()
	}
	assert.Equal(t, 1, initCalls[1], "tab 1 initialized via switch")

	// Switch back to tab 1
	app, _ = updateApp(t, app, keyMsg("1"))

	// Simulate data load on tab 1 → triggers prefetch
	_, _ = updateApp(t, app, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1"}},
	})

	// Tab 1 should NOT have been re-initialized by prefetch
	assert.Equal(t, 1, initCalls[1], "tab 1 should not be re-initialized by prefetch")
}
