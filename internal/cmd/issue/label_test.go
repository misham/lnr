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

func TestLabelListCmd_AllLabels(t *testing.T) {
	fc := &fakeClient{
		labels: []api.IssueLabel{
			{ID: "label-bug", Name: "Bug"},
			{ID: "label-feat", Name: "Feature"},
			{ID: "label-spike", Name: "Spike"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newLabelCmd(f)
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Bug")
	assert.Contains(t, buf.String(), "Feature")
	assert.Contains(t, buf.String(), "Spike")
}

func TestLabelListCmd_IssueLabels(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Labels: []api.IssueLabel{
				{ID: "label-bug", Name: "Bug"},
				{ID: "label-feat", Name: "Feature"},
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

	cmd := newLabelCmd(f)
	cmd.SetArgs([]string{"list", "ENG-1"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Bug")
	assert.Contains(t, buf.String(), "Feature")
	assert.NotContains(t, buf.String(), "Spike")
}

func TestLabelListCmd_IssueNoLabels(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Labels:     []api.IssueLabel{},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newLabelCmd(f)
	cmd.SetArgs([]string{"list", "ENG-1"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "No labels")
}

func TestLabelAddCmd_AddsLabel(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Team:       api.Team{ID: "team-1"},
		},
		labels: []api.IssueLabel{
			{ID: "label-bug", Name: "Bug"},
			{ID: "label-feat", Name: "Feature"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newLabelCmd(f)
	cmd.SetArgs([]string{"add", "ENG-1", "--label", "Bug"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, [2]string{"ENG-1", "label-bug"}, fc.addLabelArgs)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Bug")
}

func TestLabelAddCmd_CaseInsensitive(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Team:       api.Team{ID: "team-1"},
		},
		labels: []api.IssueLabel{
			{ID: "label-bug", Name: "Bug"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newLabelCmd(f)
	cmd.SetArgs([]string{"add", "ENG-1", "--label", "bug"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, [2]string{"ENG-1", "label-bug"}, fc.addLabelArgs)
}

func TestLabelAddCmd_InvalidLabel(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Team:       api.Team{ID: "team-1"},
		},
		labels: []api.IssueLabel{
			{ID: "label-bug", Name: "Bug"},
			{ID: "label-feat", Name: "Feature"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newLabelCmd(f)
	cmd.SetArgs([]string{"add", "ENG-1", "--label", "nonexistent"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
	assert.Contains(t, err.Error(), "Bug")
	assert.Contains(t, err.Error(), "Feature")
}

func TestLabelRemoveCmd_RemovesLabel(t *testing.T) {
	fc := &fakeClient{
		issue: &api.Issue{
			Identifier: "ENG-1",
			Team:       api.Team{ID: "team-1"},
		},
		labels: []api.IssueLabel{
			{ID: "label-bug", Name: "Bug"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newLabelCmd(f)
	cmd.SetArgs([]string{"remove", "ENG-1", "--label", "Bug"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, [2]string{"ENG-1", "label-bug"}, fc.rmLabelArgs)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Bug")
}
