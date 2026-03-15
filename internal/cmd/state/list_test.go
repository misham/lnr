package state

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

type fakeClient struct {
	states []api.WorkflowState
	err    error
}

func (c *fakeClient) ListTeams(_ context.Context) ([]api.Team, error) { return nil, nil }
func (c *fakeClient) Viewer(_ context.Context) (*api.User, error)     { return nil, nil }
func (c *fakeClient) ListIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}
func (c *fakeClient) GetIssue(_ context.Context, _ string) (*api.Issue, error) { return nil, nil }
func (c *fakeClient) SearchIssues(_ context.Context, _, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeClient) CreateIssue(_ context.Context, _ api.IssueCreateInput) (*api.Issue, error) {
	return nil, nil
}

func (c *fakeClient) UpdateIssue(_ context.Context, _ string, _ api.IssueUpdateInput) (*api.Issue, error) {
	return nil, nil
}
func (c *fakeClient) ArchiveIssue(_ context.Context, _ string) error { return nil }
func (c *fakeClient) ListWorkflowStates(_ context.Context, _ string) ([]api.WorkflowState, error) {
	return c.states, c.err
}

func (c *fakeClient) ListComments(_ context.Context, _ string, _ int, _ string) ([]api.Comment, api.PageInfo, error) {
	return nil, api.PageInfo{}, nil
}

func (c *fakeClient) CreateComment(_ context.Context, _, _ string) (*api.Comment, error) {
	return nil, nil
}
func (c *fakeClient) AddIssueLabel(_ context.Context, _, _ string) error    { return nil }
func (c *fakeClient) RemoveIssueLabel(_ context.Context, _, _ string) error { return nil }
func (c *fakeClient) ListLabels(_ context.Context, _ string) ([]api.IssueLabel, error) {
	return nil, nil
}

func (c *fakeClient) ListCycles(_ context.Context, _ string, _ bool, _ int, _ string) (*api.CycleListResult, error) {
	return nil, nil
}
func (c *fakeClient) GetCycle(_ context.Context, _ string) (*api.Cycle, error) { return nil, nil }
func (c *fakeClient) GetCycleByNumber(_ context.Context, _ string, _ int) (*api.Cycle, error) {
	return nil, nil
}

func (c *fakeClient) GetActiveCycle(_ context.Context, _ string) (*api.Cycle, error) {
	return nil, nil
}
func (c *fakeClient) RemoveIssueCycle(_ context.Context, _ string) error { return nil }
func (c *fakeClient) ListCycleIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeClient) ListProjects(_ context.Context, _, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return nil, nil
}

func (c *fakeClient) GetProject(_ context.Context, _ string) (*api.Project, error) {
	return nil, nil
}

func (c *fakeClient) ListProjectIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeClient) ListInitiatives(_ context.Context, _ string, _ int, _ string) (*api.InitiativeListResult, error) {
	return nil, nil
}

func (c *fakeClient) GetInitiative(_ context.Context, _ string) (*api.Initiative, error) {
	return nil, nil
}

func (c *fakeClient) ListInitiativeProjects(_ context.Context, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return nil, nil
}

func (c *fakeClient) FileUpload(_ context.Context, _, _ string, _ int64) (*api.UploadResult, error) {
	return nil, nil
}

func (c *fakeClient) UploadToURL(_ context.Context, _ string, _ []api.UploadHeader, _ io.Reader) error {
	return nil
}

func (c *fakeClient) DownloadURL(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, nil
}

func newTestFactory(t *testing.T, fc *fakeClient) *cmdutil.Factory {
	ios := ui.NewTestIOStreams()
	dir := t.TempDir()
	store := config.NewViperStore(dir)
	_ = store.Load()
	store.SetTeamID("team-1")

	return &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}
}

func TestListCmd_ShowsStates(t *testing.T) {
	fc := &fakeClient{
		states: []api.WorkflowState{
			{ID: "s1", Name: "Backlog", Type: "backlog", Color: "#bec2c8", Position: 0},
			{ID: "s2", Name: "In Progress", Type: "started", Color: "#f2c94c", Position: 1},
		},
	}

	f := newTestFactory(t, fc)
	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Backlog")
	assert.Contains(t, buf.String(), "In Progress")
}

func TestListCmd_NoStates(t *testing.T) {
	fc := &fakeClient{
		states: []api.WorkflowState{},
	}

	f := newTestFactory(t, fc)
	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "No workflow states found")
}

func TestListCmd_SortsByPosition(t *testing.T) {
	fc := &fakeClient{
		states: []api.WorkflowState{
			{ID: "s3", Name: "Done", Type: "completed", Color: "#5e6ad2", Position: 3},
			{ID: "s1", Name: "Backlog", Type: "backlog", Color: "#bec2c8", Position: 0},
			{ID: "s2", Name: "In Progress", Type: "started", Color: "#f2c94c", Position: 2},
		},
	}

	f := newTestFactory(t, fc)
	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	backlogIdx := strings.Index(out, "Backlog")
	progressIdx := strings.Index(out, "In Progress")
	doneIdx := strings.Index(out, "Done")
	assert.Less(t, backlogIdx, progressIdx, "Backlog should appear before In Progress")
	assert.Less(t, progressIdx, doneIdx, "In Progress should appear before Done")
}

func TestListCmd_APIError(t *testing.T) {
	fc := &fakeClient{
		err: errors.New("api fail"),
	}

	f := newTestFactory(t, fc)
	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "list states")
}
