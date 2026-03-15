package issue

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestViewCmd_ShowsIssue(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier:    "ENG-1",
			Title:         "Fix login bug",
			Description:   "The login page is broken.",
			Priority:      1,
			PriorityLabel: "Urgent",
			State:         api.WorkflowState{Name: "In Progress"},
			Team:          api.Team{Key: "ENG", Name: "Engineering"},
			Assignee:      &api.User{Name: "Jane"},
			Labels:        []api.IssueLabel{{Name: "bug"}},
			Comments: []api.Comment{
				{Body: "Working on it", User: &api.User{Name: "Jane"}, CreatedAt: now},
			},
			CreatedAt: now,
			UpdatedAt: now,
			URL:       "https://linear.app/issue/ENG-1",
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newViewCmd(f)
	cmd.SetArgs([]string{"ENG-1"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "ENG-1")
	assert.Contains(t, out, "Fix login bug")
	assert.Contains(t, out, "The login page is broken.")
	assert.Contains(t, out, "Urgent")
	assert.Contains(t, out, "Jane")
	assert.Contains(t, out, "bug")
	assert.Contains(t, out, "Working on it")
}

func TestViewCmd_NoArgs(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{}, nil
		},
	}

	cmd := newViewCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
}
