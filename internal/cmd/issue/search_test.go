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

func TestSearchCmd_ShowsResults(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		issues: &api.IssueListResult{
			Issues: []api.Issue{
				{Identifier: "ENG-5", Title: "Fix search bug", State: api.WorkflowState{Name: "In Progress"}, CreatedAt: now},
			},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newSearchCmd(f)
	cmd.SetArgs([]string{"--query", "search bug"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "ENG-5")
	assert.Contains(t, buf.String(), "Fix search bug")
}

func TestSearchCmd_NoResults(t *testing.T) {
	fc := &fakeClient{
		issues: &api.IssueListResult{Issues: []api.Issue{}},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newSearchCmd(f)
	cmd.SetArgs([]string{"--query", "nonexistent"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "No issues found")
}

func TestSearchCmd_NoQuery(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{}, nil
		},
	}

	cmd := newSearchCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
}
