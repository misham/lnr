package tui

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSI removes ANSI escape sequences for test assertions.
func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func newTestIssueTab(fc *fakeClient) *issueTab {
	ctx := context.Background()
	tab := NewIssueTab(ctx, fc, "team-1")
	tab.size = tea.WindowSizeMsg{Width: 80, Height: 20}
	return tab
}

func updateIssueTab(t *testing.T, tab *issueTab, msg tea.Msg) (*issueTab, tea.Cmd) {
	t.Helper()
	m, cmd := tab.Update(msg)
	result, ok := m.(*issueTab)
	require.True(t, ok)
	return result, cmd
}

func TestIssueTab_Loading(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	cmd := tab.Init()
	assert.NotNil(t, cmd)
	assert.Equal(t, stateLoading, tab.state)
	assert.Contains(t, tab.View(), "Loading issues")
}

func TestIssueTab_ListRendering(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "Fix bug", State: api.WorkflowState{Name: "In Progress"}, PriorityLabel: "Urgent", CreatedAt: now},
			{ID: "2", Identifier: "ENG-2", Title: "Add feature", State: api.WorkflowState{Name: "Backlog"}, PriorityLabel: "Normal", CreatedAt: now},
		},
	})

	assert.Equal(t, stateList, tab.state)
	view := tab.View()
	assert.Contains(t, view, "ENG-1")
	assert.Contains(t, view, "Fix bug")
	assert.Contains(t, view, "ENG-2")
	assert.Contains(t, view, "Add feature")
}

func TestIssueTab_CursorNavigation(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "First"},
			{ID: "2", Identifier: "ENG-2", Title: "Second"},
			{ID: "3", Identifier: "ENG-3", Title: "Third"},
		},
	})

	assert.Equal(t, 0, tab.cursor)

	// Move down
	tab, _ = updateIssueTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Move down again
	tab, _ = updateIssueTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move down at bottom — should stay
	tab, _ = updateIssueTab(t, tab, keyMsg("j"))
	assert.Equal(t, 2, tab.cursor)

	// Move up
	tab, _ = updateIssueTab(t, tab, keyMsg("k"))
	assert.Equal(t, 1, tab.cursor)

	// View should highlight current row with accent bar
	view := tab.View()
	assert.Contains(t, view, "▎")
}

func TestIssueTab_DetailView(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "issue-1", Identifier: "ENG-1", Title: "Fix bug"},
		},
	})

	// Press Enter to load detail
	tab, cmd := updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, cmd)

	// Simulate detail loaded
	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "issue-1", Identifier: "ENG-1", Title: "Fix bug",
			Description:   "This is a detailed description",
			State:         api.WorkflowState{Name: "In Progress"},
			PriorityLabel: "Urgent",
			CreatedAt:     now,
			Comments: []api.Comment{
				{Body: "Great work!", User: &api.User{DisplayName: "Alice"}, CreatedAt: now},
			},
			Labels: []api.IssueLabel{
				{Name: "bug", Color: "#ff0000"},
			},
		},
	})

	assert.Equal(t, stateDetail, tab.state)
	view := tab.View()
	plain := stripANSI(view)
	assert.Contains(t, plain, "ENG-1")
	assert.Contains(t, plain, "Fix bug")
	assert.Contains(t, plain, "detailed description")
	assert.Contains(t, plain, "Great work!")
	assert.Contains(t, plain, "Alice")
	assert.Contains(t, plain, "bug")
}

func TestIssueTab_BackToList(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "First"},
			{ID: "2", Identifier: "ENG-2", Title: "Second"},
		},
	})

	// Move cursor to second item
	tab, _ = updateIssueTab(t, tab, keyMsg("j"))
	assert.Equal(t, 1, tab.cursor)

	// Go to detail
	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{ID: "2", Identifier: "ENG-2", Title: "Second", CreatedAt: time.Now()},
	})
	assert.Equal(t, stateDetail, tab.state)

	// Press Esc to go back
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, stateList, tab.state)
	assert.Equal(t, 1, tab.cursor, "cursor should be preserved")
}

func TestIssueTab_Pagination(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	// First page
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "1", Identifier: "ENG-1", Title: "First"},
		},
		PageInfo: api.PageInfo{HasNextPage: true, EndCursor: "cursor-1"},
	})

	assert.Len(t, tab.items, 1)
	assert.True(t, tab.pageInfo.HasNextPage)

	// Press n to load next page
	tab, cmd := updateIssueTab(t, tab, keyMsg("n"))
	require.NotNil(t, cmd)

	// Simulate next page loaded
	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{
			{ID: "2", Identifier: "ENG-2", Title: "Second"},
		},
		PageInfo: api.PageInfo{HasNextPage: false},
	})

	assert.Len(t, tab.items, 2, "items should accumulate")
	assert.False(t, tab.pageInfo.HasNextPage)
}

func TestIssueTab_Error(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Err: fmt.Errorf("network error"),
	})

	assert.Equal(t, stateError, tab.state)
	view := tab.View()
	assert.Contains(t, view, "network error")
	assert.Contains(t, view, "retry")

	// Press r to retry
	tab, cmd := updateIssueTab(t, tab, keyMsg("r"))
	assert.Equal(t, stateLoading, tab.state)
	assert.NotNil(t, cmd)
}

func TestIssueTab_EmptyList(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{},
	})

	assert.Equal(t, stateList, tab.state)
	assert.Contains(t, tab.View(), "No issues found")
}

func TestIssueTab_TopBottom(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	issues := make([]api.Issue, 10)
	for i := range issues {
		issues[i] = api.Issue{ID: fmt.Sprintf("%d", i), Identifier: fmt.Sprintf("ENG-%d", i+1), Title: fmt.Sprintf("Issue %d", i+1)}
	}

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{Issues: issues})

	// Jump to bottom
	tab, _ = updateIssueTab(t, tab, keyMsg("G"))
	assert.Equal(t, 9, tab.cursor)

	// Jump to top
	tab, _ = updateIssueTab(t, tab, keyMsg("g"))
	assert.Equal(t, 0, tab.cursor)
}

func TestIssueTab_DetailScrolling(t *testing.T) {
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "1", Identifier: "ENG-1", Title: "Test",
			Description: "Long description", CreatedAt: time.Now(),
		},
	})

	assert.Equal(t, stateDetail, tab.state)

	// Scrolling should not error
	tab, _ = updateIssueTab(t, tab, keyMsg("j"))
	_, _ = updateIssueTab(t, tab, keyMsg("k"))
}

func TestRenderMarkdown_Basic(t *testing.T) {
	input := "# Heading\n\nSome **bold** text."
	out := renderMarkdown(input, 80)
	// Glamour should transform the input (at minimum, the output differs from raw text)
	assert.NotEqual(t, input, out)
	assert.Contains(t, out, "Heading")
	assert.Contains(t, out, "bold")
}

func TestRenderMarkdown_Fallback(t *testing.T) {
	// Empty string should return empty (or near-empty) output
	out := renderMarkdown("", 80)
	assert.Empty(t, strings.TrimSpace(out))
}

func TestRenderMarkdown_WordWrap(t *testing.T) {
	long := "This is a very long sentence that should be wrapped at the specified width to ensure readability in narrow terminal windows."
	out := renderMarkdown(long, 40)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// With wrap at 40, output should have more than one line
	assert.Greater(t, len(lines), 1, "expected word wrap to produce multiple lines")
}

func TestIssueTab_DetailRendersMarkdown(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "1", Identifier: "ENG-1", Title: "Test",
			Description: "## Section\n\nSome **bold** text.",
			CreatedAt:   now,
		},
	})

	assert.Equal(t, stateDetail, tab.state)
	detail := tab.renderDetail()
	// Glamour renders with ANSI styling, so output should contain content
	assert.Contains(t, detail, "Section")
	assert.Contains(t, detail, "bold")
	// Bold markers should be replaced with ANSI escape sequences
	assert.NotContains(t, detail, "**bold**")
}

func TestIssueTab_DetailFilesSection(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "1", Identifier: "ENG-1", Title: "Test",
			Description: "See [report.pdf](https://uploads.linear.app/org/report.pdf)",
			CreatedAt:   now,
		},
	})

	detail := tab.renderDetail()
	plain := stripANSI(detail)
	assert.Contains(t, plain, "Files")
	assert.Contains(t, plain, "report.pdf")
}

func TestIssueTab_DetailNoFilesSection(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "1", Identifier: "ENG-1", Title: "Test",
			Description: "No file links here",
			CreatedAt:   now,
		},
	})

	detail := tab.renderDetail()
	plain := stripANSI(detail)
	assert.NotContains(t, plain, "Files")
}

func TestIssueTab_OpenKeyBinding(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "1", Identifier: "ENG-1", Title: "Test",
			Description: "See [report.pdf](https://uploads.linear.app/org/report.pdf)",
			CreatedAt:   now,
		},
	})

	// Press 'o' should dispatch open command
	_, cmd := updateIssueTab(t, tab, keyMsg("o"))
	assert.NotNil(t, cmd, "expected open command to be dispatched")
}

func TestIssueTab_UploadKeyBinding(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "1", Identifier: "ENG-1", Title: "Test",
			CreatedAt: now,
		},
	})

	// Press 'a' should enter upload mode
	tab, _ = updateIssueTab(t, tab, keyMsg("a"))
	assert.True(t, tab.uploading, "expected uploading mode to be active")

	// Esc should cancel upload mode
	tab, _ = updateIssueTab(t, tab, tea.KeyMsg{Type: tea.KeyEscape})
	assert.False(t, tab.uploading, "expected uploading mode to be cancelled")
}

func TestIssueTab_ImageFetchedMsg(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{}
	tab := newTestIssueTab(fc)

	tab, _ = updateIssueTab(t, tab, IssuesLoadedMsg{
		Issues: []api.Issue{{ID: "1", Identifier: "ENG-1", Title: "Test"}},
	})

	tab, _ = updateIssueTab(t, tab, IssueDetailLoadedMsg{
		Issue: &api.Issue{
			ID: "1", Identifier: "ENG-1", Title: "Test",
			Description: "![img](https://uploads.linear.app/org/image.png)",
			CreatedAt:   now,
		},
	})

	// Simulate image fetch completion
	tab, _ = updateIssueTab(t, tab, ImageFetchedMsg{
		URL:  "https://uploads.linear.app/org/image.png",
		Data: []byte("fake-image-data"),
	})

	assert.Contains(t, tab.fetchedImages, "https://uploads.linear.app/org/image.png")
}
