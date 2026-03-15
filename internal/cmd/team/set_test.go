package team

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

func TestSetCmd_WithValidKey(t *testing.T) {
	ios := ui.NewTestIOStreams()
	dir := t.TempDir()
	store := config.NewViperStore(dir)
	require.NoError(t, store.Load())

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{
				teams: []api.Team{
					{ID: "team-uuid-1", Key: "ENG", Name: "Engineering"},
					{ID: "team-uuid-2", Key: "DES", Name: "Design"},
				},
			}, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newSetCmd(f)
	cmd.SetArgs([]string{"ENG"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "team-uuid-1", store.TeamID())

	buf, ok := ios.Out.(*bytes.Buffer)
	require.True(t, ok)
	assert.Contains(t, buf.String(), "ENG")
	assert.Contains(t, buf.String(), "Engineering")
}

func TestSetCmd_WithValidKey_CaseInsensitive(t *testing.T) {
	ios := ui.NewTestIOStreams()
	dir := t.TempDir()
	store := config.NewViperStore(dir)
	require.NoError(t, store.Load())

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{
				teams: []api.Team{
					{ID: "team-uuid-1", Key: "ENG", Name: "Engineering"},
				},
			}, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newSetCmd(f)
	cmd.SetArgs([]string{"eng"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "team-uuid-1", store.TeamID())
}

func TestSetCmd_WithInvalidKey(t *testing.T) {
	ios := ui.NewTestIOStreams()

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{
				teams: []api.Team{
					{ID: "team-uuid-1", Key: "ENG", Name: "Engineering"},
					{ID: "team-uuid-2", Key: "DES", Name: "Design"},
				},
			}, nil
		},
		Config: func() (config.Store, error) {
			return config.NewViperStore(t.TempDir()), nil
		},
	}

	cmd := newSetCmd(f)
	cmd.SetArgs([]string{"INVALID"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ENG")
	assert.Contains(t, err.Error(), "DES")
}

func TestSetCmd_NoArgNonTTY(t *testing.T) {
	ios := ui.NewTestIOStreams()

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{
				teams: []api.Team{
					{ID: "team-uuid-1", Key: "ENG", Name: "Engineering"},
				},
			}, nil
		},
	}

	cmd := newSetCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "interactive")
}
