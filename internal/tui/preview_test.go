package tui

import (
	"context"
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/misham/linear-cli/internal/api"
)

func TestRenderPreview(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fc := &fakeClient{}
	tabs := []TabModel{
		NewIssueTab(ctx, fc, "team-1"),
		NewCycleTab(ctx, fc, "team-1"),
		NewProjectTab(ctx, fc, "team-1"),
		NewInitiativeTab(ctx, fc),
	}
	tabNames := []string{"Issues", "Cycles", "Projects", "Initiatives"}
	app := NewApp(ctx, cancel, tabs, tabNames, "Engineering")
	_ = app.Init()

	m, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	app, _ = m.(AppModel)

	now := time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC)
	issues := make([]api.Issue, 50)
	for i := range issues {
		issues[i] = api.Issue{
			ID: fmt.Sprintf("%d", i+1), Identifier: fmt.Sprintf("ENG-%d", i+1),
			Title: fmt.Sprintf("Issue number %d", i+1),
			State: api.WorkflowState{Name: "In Progress"}, PriorityLabel: "High", CreatedAt: now,
		}
	}
	m, _ = app.Update(IssuesLoadedMsg{
		Issues:   issues,
		PageInfo: api.PageInfo{HasNextPage: true, EndCursor: "cursor1"},
	})
	app, _ = m.(AppModel)

	fmt.Println(app.View())
}
