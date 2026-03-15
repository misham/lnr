package cmd

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

func TestTuiCmd_Exists(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	cmd := newTuiCmd(f)
	assert.Equal(t, "tui", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
}

func TestTuiCmd_NonTTY(t *testing.T) {
	ios := ui.NewTestIOStreams() // plain=true → non-TTY
	f := &cmdutil.Factory{IO: ios}

	cmd := newTuiCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "interactive terminal")
}

func TestTuiCmd_MissingAuth(t *testing.T) {
	ios := ui.NewTestIOStreams()
	ios.SetPlain(false) // simulate TTY

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return nil, assert.AnError
		},
	}

	cmd := newTuiCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
}

func TestTuiCmd_MissingTeam(t *testing.T) {
	ios := ui.NewTestIOStreams()
	ios.SetPlain(false) // simulate TTY

	dir := t.TempDir()
	store := config.NewViperStore(dir)
	require.NoError(t, store.Load())

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeTuiClient{}, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newTuiCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "lnr team set")
}

type fakeTuiClient struct{}

func (c *fakeTuiClient) ListTeams(_ context.Context) ([]api.Team, error) { return nil, nil }
func (c *fakeTuiClient) Viewer(_ context.Context) (*api.User, error)     { return nil, nil }

func (c *fakeTuiClient) ListIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) GetIssue(_ context.Context, _ string) (*api.Issue, error) { return nil, nil }

func (c *fakeTuiClient) SearchIssues(_ context.Context, _, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) CreateIssue(_ context.Context, _ api.IssueCreateInput) (*api.Issue, error) {
	return nil, nil
}

func (c *fakeTuiClient) UpdateIssue(_ context.Context, _ string, _ api.IssueUpdateInput) (*api.Issue, error) {
	return nil, nil
}

func (c *fakeTuiClient) ArchiveIssue(_ context.Context, _ string) error { return nil }

func (c *fakeTuiClient) ListWorkflowStates(_ context.Context, _ string) ([]api.WorkflowState, error) {
	return nil, nil
}

func (c *fakeTuiClient) ListComments(_ context.Context, _ string, _ int, _ string) ([]api.Comment, api.PageInfo, error) {
	return nil, api.PageInfo{}, nil
}

func (c *fakeTuiClient) CreateComment(_ context.Context, _, _ string) (*api.Comment, error) {
	return nil, nil
}

func (c *fakeTuiClient) AddIssueLabel(_ context.Context, _, _ string) error    { return nil }
func (c *fakeTuiClient) RemoveIssueLabel(_ context.Context, _, _ string) error { return nil }

func (c *fakeTuiClient) ListLabels(_ context.Context, _ string) ([]api.IssueLabel, error) {
	return nil, nil
}

func (c *fakeTuiClient) ListCycles(_ context.Context, _ string, _ bool, _ int, _ string) (*api.CycleListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) GetCycle(_ context.Context, _ string) (*api.Cycle, error) { return nil, nil }

func (c *fakeTuiClient) GetCycleByNumber(_ context.Context, _ string, _ int) (*api.Cycle, error) {
	return nil, nil
}

func (c *fakeTuiClient) GetActiveCycle(_ context.Context, _ string) (*api.Cycle, error) {
	return nil, nil
}

func (c *fakeTuiClient) RemoveIssueCycle(_ context.Context, _ string) error { return nil }

func (c *fakeTuiClient) ListCycleIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) ListProjects(_ context.Context, _, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) GetProject(_ context.Context, _ string) (*api.Project, error) {
	return nil, nil
}

func (c *fakeTuiClient) ListProjectIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) ListInitiatives(_ context.Context, _ string, _ int, _ string) (*api.InitiativeListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) GetInitiative(_ context.Context, _ string) (*api.Initiative, error) {
	return nil, nil
}

func (c *fakeTuiClient) ListInitiativeProjects(_ context.Context, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) FileUpload(_ context.Context, _, _ string, _ int64) (*api.UploadResult, error) {
	return nil, nil
}

func (c *fakeTuiClient) UploadToURL(_ context.Context, _ string, _ []api.UploadHeader, _ io.Reader) error {
	return nil
}

func (c *fakeTuiClient) DownloadURL(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, nil
}
