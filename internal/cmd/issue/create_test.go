package issue

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

func TestCreateCmd_WithFlags(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-10",
			Title:      "New bug",
			URL:        "https://linear.app/issue/ENG-10",
			State:      api.WorkflowState{Name: "Backlog"},
		},
	}

	ios := ui.NewTestIOStreams()
	dir := t.TempDir()
	store := config.NewViperStore(dir)
	require.NoError(t, store.Load())
	store.SetTeamID("team-1")

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--title", "New bug", "--description", "Something broke", "--priority", "1"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "New bug", fc.createInput.Title)
	assert.Equal(t, "team-1", fc.createInput.TeamID)
	assert.Equal(t, "Something broke", fc.createInput.Description)
	assert.Equal(t, 1, fc.createInput.Priority)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "ENG-10")
}

func TestCreateCmd_NoTitle_NonTTY(t *testing.T) {
	ios := ui.NewTestIOStreams()
	dir := t.TempDir()
	store := config.NewViperStore(dir)
	require.NoError(t, store.Load())
	store.SetTeamID("team-1")

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{}, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}
