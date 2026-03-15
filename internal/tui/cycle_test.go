package tui

import (
	"context"
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func newTestCycleTab(fc *fakeClient) *cycleTab {
	ctx := context.Background()
	tab := NewCycleTab(ctx, fc, "team-1")
	tab.size = tea.WindowSizeMsg{Width: 120, Height: 20}
	return tab
}

func updateCycleTab(t *testing.T, tab *cycleTab, msg tea.Msg) (*cycleTab, tea.Cmd) {
	t.Helper()
	m, cmd := tab.Update(msg)
	result, ok := m.(*cycleTab)
	require.True(t, ok)
	return result, cmd
}

func TestCycleTab_Loading(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	cmd := tab.Init()
	assert.NotNil(t, cmd)
	assert.Equal(t, stateLoading, tab.state)
	assert.Contains(t, tab.View(), "Loading cycles")
}

func TestCycleTab_ListRendering(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Cycles: []api.Cycle{
			{ID: "1", Number: 1, Name: "Sprint 1", IsActive: true, StartsAt: now, EndsAt: now.Add(14 * 24 * time.Hour), Progress: 0.5, CreatedAt: now, UpdatedAt: now},
			{ID: "2", Number: 2, Name: "Sprint 2", IsNext: true, StartsAt: now.Add(14 * 24 * time.Hour), EndsAt: now.Add(28 * 24 * time.Hour), Progress: 0, CreatedAt: now, UpdatedAt: now},
		},
	})

	assert.Equal(t, stateList, tab.state)
	view := tab.View()
	assert.Contains(t, view, "Cycle #1")
	assert.Contains(t, view, "Sprint 1")
	assert.Contains(t, view, "Active")
	assert.Contains(t, view, "Cycle #2")
	assert.Contains(t, view, "Sprint 2")
	assert.Contains(t, view, "Next")
}

func TestCycleTab_CursorNavigation(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Cycles: []api.Cycle{
			{ID: "1", Number: 1, Name: "Sprint 1", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
			{ID: "2", Number: 2, Name: "Sprint 2", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
			{ID: "3", Number: 3, Name: "Sprint 3", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
		},
	})

	assert.Equal(t, 0, tab.cursor)

	// Move down
	tab, _ = updateCycleTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Move down again
	tab, _ = updateCycleTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move down at bottom — should stay
	tab, _ = updateCycleTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move up
	tab, _ = updateCycleTab(t, tab, keyMsg("k"))
	assert.Equal(t, 1, tab.cursor)

	// View should highlight current row (accent bar marks selection)
	view := tab.View()
	assert.Contains(t, view, "▎")
}

func TestCycleTab_DetailView(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Cycles: []api.Cycle{
			{ID: "cycle-1", Number: 1, Name: "Sprint 1", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
		},
	})

	// Press Enter to load detail
	tab, cmd := updateCycleTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, cmd)

	// Simulate detail loaded
	tab, _ = updateCycleTab(t, tab, CycleDetailLoadedMsg{
		Cycle: &api.Cycle{
			ID: "cycle-1", Number: 1, Name: "Sprint 1",
			Description: "First sprint of the year",
			IsActive:    true,
			Progress:    0.75,
			StartsAt:    now,
			EndsAt:      now.Add(14 * 24 * time.Hour),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	})

	assert.Equal(t, stateDetail, tab.state)
	view := tab.View()
	assert.Contains(t, view, "Cycle #1")
	assert.Contains(t, view, "Sprint 1")
	assert.Contains(t, view, "First sprint of the year")
	assert.Contains(t, view, "Active")
	assert.Contains(t, view, "75%")
}

func TestCycleTab_BackToList(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Cycles: []api.Cycle{
			{ID: "1", Number: 1, Name: "Sprint 1", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
			{ID: "2", Number: 2, Name: "Sprint 2", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
		},
	})

	// Move cursor to second item
	tab, _ = updateCycleTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Go to detail
	tab, _ = updateCycleTab(t, tab, CycleDetailLoadedMsg{
		Cycle: &api.Cycle{ID: "2", Number: 2, Name: "Sprint 2", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
	})
	assert.Equal(t, stateDetail, tab.state)

	// Press Esc to go back
	tab, _ = updateCycleTab(t, tab, tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, stateList, tab.state)
	assert.Equal(t, 1, tab.cursor, "cursor should be preserved")
}

func TestCycleTab_Pagination(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	// First page
	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Cycles: []api.Cycle{
			{ID: "1", Number: 1, Name: "Sprint 1", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
		},
		PageInfo: api.PageInfo{HasNextPage: true, EndCursor: "cursor-1"},
	})

	assert.Len(t, tab.items, 1)
	assert.True(t, tab.pageInfo.HasNextPage)

	// Press n to load next page
	tab, cmd := updateCycleTab(t, tab, keyMsg("n"))
	require.NotNil(t, cmd)

	// Simulate next page loaded
	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Cycles: []api.Cycle{
			{ID: "2", Number: 2, Name: "Sprint 2", StartsAt: now, EndsAt: now, CreatedAt: now, UpdatedAt: now},
		},
		PageInfo: api.PageInfo{HasNextPage: false},
	})

	assert.Len(t, tab.items, 2, "items should accumulate")
	assert.False(t, tab.pageInfo.HasNextPage)
}

func TestCycleTab_Error(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Err: fmt.Errorf("network error"),
	})

	assert.Equal(t, stateError, tab.state)
	view := tab.View()
	assert.Contains(t, view, "network error")
	assert.Contains(t, view, "retry")

	// Press r to retry
	tab, cmd := updateCycleTab(t, tab, keyMsg("r"))
	assert.Equal(t, stateLoading, tab.state)
	assert.NotNil(t, cmd)
}

func TestCycleTab_EmptyList(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestCycleTab(fc)

	tab, _ = updateCycleTab(t, tab, CyclesLoadedMsg{
		Cycles: []api.Cycle{},
	})

	assert.Equal(t, stateList, tab.state)
	assert.Contains(t, tab.View(), "No cycles found")
}
