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

func newTestProjectTab(fc *fakeClient) *projectTab {
	ctx := context.Background()
	tab := NewProjectTab(ctx, fc, "team-1")
	tab.size = tea.WindowSizeMsg{Width: 120, Height: 20}
	return tab
}

func updateProjectTab(t *testing.T, tab *projectTab, msg tea.Msg) (*projectTab, tea.Cmd) {
	t.Helper()
	m, cmd := tab.Update(msg)
	result, ok := m.(*projectTab)
	require.True(t, ok)
	return result, cmd
}

func TestProjectTab_Loading(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	cmd := tab.Init()
	assert.NotNil(t, cmd)
	assert.Equal(t, stateLoading, tab.state)
	assert.Contains(t, tab.View(), "Loading projects")
}

func TestProjectTab_ListRendering(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Projects: []api.Project{
			{ID: "1", Name: "Project Alpha", Status: api.ProjectStatus{Name: "Started", Type: "started"}, Progress: 0.6, Lead: &api.User{DisplayName: "Alice"}},
			{ID: "2", Name: "Project Beta", Status: api.ProjectStatus{Name: "Planned", Type: "planned"}, Progress: 0},
		},
	})

	assert.Equal(t, stateList, tab.state)
	view := tab.View()
	assert.Contains(t, view, "Project Alpha")
	assert.Contains(t, view, "Started")
	assert.Contains(t, view, "Alice")
	assert.Contains(t, view, "Project Beta")
	assert.Contains(t, view, "Planned")
}

func TestProjectTab_CursorNavigation(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Projects: []api.Project{
			{ID: "1", Name: "Project A"},
			{ID: "2", Name: "Project B"},
			{ID: "3", Name: "Project C"},
		},
	})

	assert.Equal(t, 0, tab.cursor)

	// Move down
	tab, _ = updateProjectTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Move down again
	tab, _ = updateProjectTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move down at bottom — should stay
	tab, _ = updateProjectTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move up
	tab, _ = updateProjectTab(t, tab, keyMsg("k"))
	assert.Equal(t, 1, tab.cursor)

	// View should highlight current row (accent bar marks selection)
	view := tab.View()
	assert.Contains(t, view, "▎")
}

func TestProjectTab_DetailView(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Projects: []api.Project{
			{ID: "proj-1", Name: "Project Alpha"},
		},
	})

	// Press Enter to load detail
	tab, cmd := updateProjectTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, cmd)

	// Simulate detail loaded
	tab, _ = updateProjectTab(t, tab, ProjectDetailLoadedMsg{
		Project: &api.Project{
			ID:          "proj-1",
			Name:        "Project Alpha",
			Description: "A very important project",
			Status:      api.ProjectStatus{Name: "Started", Type: "started"},
			Progress:    0.45,
			Lead:        &api.User{DisplayName: "Alice"},
			StartDate:   "2026-01-01",
			TargetDate:  "2026-06-30",
			URL:         "https://linear.app/team/project/alpha",
			Milestones: []api.ProjectMilestone{
				{ID: "m1", Name: "Phase 1", TargetDate: "2026-03-01"},
				{ID: "m2", Name: "Phase 2", TargetDate: "2026-06-01"},
			},
		},
	})

	assert.Equal(t, stateDetail, tab.state)
	view := tab.View()
	assert.Contains(t, view, "Project Alpha")
	assert.Contains(t, view, "A very important project")
	assert.Contains(t, view, "Started")
	assert.Contains(t, view, "45%")
	assert.Contains(t, view, "Alice")
	assert.Contains(t, view, "Phase 1")
	assert.Contains(t, view, "Phase 2")
	assert.Contains(t, view, "Milestones")
}

func TestProjectTab_BackToList(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Projects: []api.Project{
			{ID: "1", Name: "Project A"},
			{ID: "2", Name: "Project B"},
		},
	})

	// Move cursor to second item
	tab, _ = updateProjectTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Go to detail
	tab, _ = updateProjectTab(t, tab, ProjectDetailLoadedMsg{
		Project: &api.Project{ID: "2", Name: "Project B", Status: api.ProjectStatus{Name: "Planned"}},
	})
	assert.Equal(t, stateDetail, tab.state)

	// Press Esc to go back
	tab, _ = updateProjectTab(t, tab, tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, stateList, tab.state)
	assert.Equal(t, 1, tab.cursor, "cursor should be preserved")
}

func TestProjectTab_Pagination(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	// First page
	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Projects: []api.Project{
			{ID: "1", Name: "Project A"},
		},
		PageInfo: api.PageInfo{HasNextPage: true, EndCursor: "cursor-1"},
	})

	assert.Len(t, tab.items, 1)
	assert.True(t, tab.pageInfo.HasNextPage)

	// Press n to load next page
	tab, cmd := updateProjectTab(t, tab, keyMsg("n"))
	require.NotNil(t, cmd)

	// Simulate next page loaded
	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Projects: []api.Project{
			{ID: "2", Name: "Project B"},
		},
		PageInfo: api.PageInfo{HasNextPage: false},
	})

	assert.Len(t, tab.items, 2, "items should accumulate")
	assert.False(t, tab.pageInfo.HasNextPage)
}

func TestProjectTab_Error(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Err: fmt.Errorf("network error"),
	})

	assert.Equal(t, stateError, tab.state)
	view := tab.View()
	assert.Contains(t, view, "network error")
	assert.Contains(t, view, "retry")

	// Press r to retry
	tab, cmd := updateProjectTab(t, tab, keyMsg("r"))
	assert.Equal(t, stateLoading, tab.state)
	assert.NotNil(t, cmd)
}

func TestProjectTab_EmptyList(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestProjectTab(fc)

	tab, _ = updateProjectTab(t, tab, ProjectsLoadedMsg{
		Projects: []api.Project{},
	})

	assert.Equal(t, stateList, tab.state)
	assert.Contains(t, tab.View(), "No projects found")
}
