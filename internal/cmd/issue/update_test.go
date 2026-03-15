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

func TestUpdateCmd_TitleOnly(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Title:      "Updated title",
			URL:        "https://linear.app/issue/ENG-1",
			State:      api.WorkflowState{Name: "In Progress"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newUpdateCmd(f)
	cmd.SetArgs([]string{"ENG-1", "--title", "Updated title"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "ENG-1", fc.updateID)
	require.NotNil(t, fc.updateInput.Title)
	assert.Equal(t, "Updated title", *fc.updateInput.Title)
	assert.Nil(t, fc.updateInput.Description)
	assert.Nil(t, fc.updateInput.Priority)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "ENG-1")
}

func TestUpdateCmd_MultipleFields(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Title:      "Updated",
			URL:        "https://linear.app/issue/ENG-1",
			State:      api.WorkflowState{Name: "In Progress"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newUpdateCmd(f)
	cmd.SetArgs([]string{"ENG-1", "--title", "New title", "--priority", "2", "--description", "New desc"})

	err := cmd.Execute()
	require.NoError(t, err)

	require.NotNil(t, fc.updateInput.Title)
	assert.Equal(t, "New title", *fc.updateInput.Title)
	require.NotNil(t, fc.updateInput.Priority)
	assert.Equal(t, 2, *fc.updateInput.Priority)
	require.NotNil(t, fc.updateInput.Description)
	assert.Equal(t, "New desc", *fc.updateInput.Description)
}

func TestUpdateCmd_NoArgs(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{}, nil
		},
	}

	cmd := newUpdateCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
}
