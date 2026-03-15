package issue

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestCloseCmd_ClosesIssue(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Title:      "Fix bug",
			URL:        "https://linear.app/issue/ENG-1",
			Team:       api.Team{ID: "team-1"},
			State:      api.WorkflowState{Name: "In Progress"},
		},
		states: []api.WorkflowState{
			{ID: "state-backlog", Name: "Backlog", Type: "backlog", Position: 0},
			{ID: "state-done", Name: "Done", Type: "completed", Position: 2},
			{ID: "state-ip", Name: "In Progress", Type: "started", Position: 1},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newCloseCmd(f)
	cmd.SetArgs([]string{"ENG-1"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "ENG-1", fc.updateID)
	require.NotNil(t, fc.updateInput.StateID)
	assert.Equal(t, "state-done", *fc.updateInput.StateID)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "ENG-1")
}

func TestCloseCmd_NoArgs(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{}, nil
		},
	}

	cmd := newCloseCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
}
