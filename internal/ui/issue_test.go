package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestPrintIssues_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	issues := []api.Issue{
		{
			Identifier:    "ENG-1",
			Title:         "Fix login bug",
			Priority:      1,
			PriorityLabel: "Urgent",
			State:         api.WorkflowState{Name: "In Progress", Type: "started"},
			Assignee:      &api.User{Name: "Jane"},
			CreatedAt:     now,
		},
		{
			Identifier:    "ENG-2",
			Title:         "Add search feature",
			Priority:      2,
			PriorityLabel: "High",
			State:         api.WorkflowState{Name: "Backlog", Type: "backlog"},
			CreatedAt:     now,
		},
	}

	err := PrintIssues(ios, issues)
	require.NoError(t, err)

	out := outString(t, ios)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 2)
	assert.Contains(t, lines[0], "ENG-1")
	assert.Contains(t, lines[0], "Fix login bug")
	assert.Contains(t, lines[0], "In Progress")
	assert.Contains(t, lines[0], "Urgent")
	assert.Contains(t, lines[0], "Jane")
	assert.Contains(t, lines[1], "ENG-2")
	assert.Contains(t, lines[1], "Add search feature")
}

func TestPrintIssues_Empty(t *testing.T) {
	ios := NewTestIOStreams()
	err := PrintIssues(ios, nil)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "No issues found")
}

func TestPrintIssues_NilAssignee(t *testing.T) {
	ios := NewTestIOStreams()
	issues := []api.Issue{
		{
			Identifier: "ENG-1",
			Title:      "Unassigned",
			State:      api.WorkflowState{Name: "Backlog"},
			CreatedAt:  time.Now(),
		},
	}

	err := PrintIssues(ios, issues)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "ENG-1")
	assert.NotContains(t, out, "<nil>")
}

func TestPrintIssueDetail_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	completedAt := time.Date(2026, 1, 16, 10, 0, 0, 0, time.UTC)
	issue := &api.Issue{
		Identifier:    "ENG-1",
		Title:         "Fix login bug",
		Description:   "The login page has a critical bug.",
		Priority:      1,
		PriorityLabel: "Urgent",
		Estimate:      3,
		DueDate:       "2026-02-01",
		State:         api.WorkflowState{Name: "Done", Type: "completed"},
		Team:          api.Team{Key: "ENG", Name: "Engineering"},
		Assignee:      &api.User{Name: "Jane"},
		Labels:        []api.IssueLabel{{Name: "bug"}, {Name: "critical"}},
		Comments: []api.Comment{
			{Body: "Looks good", User: &api.User{Name: "Alice"}, CreatedAt: now},
		},
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: &completedAt,
		URL:         "https://linear.app/issue/ENG-1",
	}

	err := PrintIssueDetail(ios, issue)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "ENG-1")
	assert.Contains(t, out, "Fix login bug")
	assert.Contains(t, out, "The login page has a critical bug.")
	assert.Contains(t, out, "Urgent")
	assert.Contains(t, out, "Done")
	assert.Contains(t, out, "Jane")
	assert.Contains(t, out, "Engineering")
	assert.Contains(t, out, "bug")
	assert.Contains(t, out, "critical")
	assert.Contains(t, out, "Looks good")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "https://linear.app/issue/ENG-1")
}

func TestPrintIssueDetail_FilesPlain(t *testing.T) {
	ios := NewTestIOStreams()
	issue := &api.Issue{
		Identifier:  "ENG-1",
		Title:       "Issue with file",
		Description: "See [chart-reviews.pdf](https://uploads.linear.app/org/abc/chart-reviews.pdf)",
		State:       api.WorkflowState{Name: "Todo"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := PrintIssueDetail(ios, issue)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "chart-reviews.pdf\thttps://uploads.linear.app/org/abc/chart-reviews.pdf")
}

func TestPrintIssueDetail_FilesStyled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	issue := &api.Issue{
		Identifier:  "ENG-1",
		Title:       "Issue with file",
		Description: "See [chart-reviews.pdf](https://uploads.linear.app/org/abc/chart-reviews.pdf)",
		State:       api.WorkflowState{Name: "Todo"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := PrintIssueDetail(ios, issue)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "chart-reviews.pdf")
	assert.Contains(t, out, "Files")
}

func TestPrintIssueDetail_NoFiles(t *testing.T) {
	ios := NewTestIOStreams()
	issue := &api.Issue{
		Identifier:  "ENG-1",
		Title:       "Plain issue",
		Description: "No file links here",
		State:       api.WorkflowState{Name: "Todo"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := PrintIssueDetail(ios, issue)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.NotContains(t, out, "Files")
}

func TestPrintIssueDetail_FilesFromComments(t *testing.T) {
	ios := NewTestIOStreams()
	issue := &api.Issue{
		Identifier: "ENG-1",
		Title:      "Issue with comment file",
		State:      api.WorkflowState{Name: "Todo"},
		Comments: []api.Comment{
			{Body: "Here is [log.txt](https://uploads.linear.app/org/xyz/log.txt)", User: &api.User{Name: "Bob"}, CreatedAt: time.Now()},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := PrintIssueDetail(ios, issue)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "log.txt\thttps://uploads.linear.app/org/xyz/log.txt")
}

func TestPrintIssueDetail_NoAssignee(t *testing.T) {
	ios := NewTestIOStreams()
	issue := &api.Issue{
		Identifier: "ENG-2",
		Title:      "Unassigned issue",
		State:      api.WorkflowState{Name: "Backlog"},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := PrintIssueDetail(ios, issue)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "ENG-2")
	assert.Contains(t, out, "Unassigned")
}
