package cmdutil

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

type fakeClient struct {
	teams []api.Team
	err   error
}

func (c *fakeClient) ListTeams(_ context.Context) ([]api.Team, error) {
	return c.teams, c.err
}

func (c *fakeClient) Viewer(_ context.Context) (*api.User, error) {
	return nil, nil
}

func (c *fakeClient) ListIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	panic("not implemented")
}

func (c *fakeClient) GetIssue(_ context.Context, _ string) (*api.Issue, error) {
	panic("not implemented")
}

func (c *fakeClient) SearchIssues(_ context.Context, _, _ string, _ int, _ string) (*api.IssueListResult, error) {
	panic("not implemented")
}

func (c *fakeClient) CreateIssue(_ context.Context, _ api.IssueCreateInput) (*api.Issue, error) {
	panic("not implemented")
}

func (c *fakeClient) UpdateIssue(_ context.Context, _ string, _ api.IssueUpdateInput) (*api.Issue, error) {
	panic("not implemented")
}

func (c *fakeClient) ArchiveIssue(_ context.Context, _ string) error {
	panic("not implemented")
}

func (c *fakeClient) ListWorkflowStates(_ context.Context, _ string) ([]api.WorkflowState, error) {
	panic("not implemented")
}

func (c *fakeClient) ListComments(_ context.Context, _ string, _ int, _ string) ([]api.Comment, api.PageInfo, error) {
	panic("not implemented")
}

func (c *fakeClient) CreateComment(_ context.Context, _, _ string) (*api.Comment, error) {
	panic("not implemented")
}

func (c *fakeClient) AddIssueLabel(_ context.Context, _, _ string) error {
	panic("not implemented")
}

func (c *fakeClient) RemoveIssueLabel(_ context.Context, _, _ string) error {
	panic("not implemented")
}

func (c *fakeClient) ListLabels(_ context.Context, _ string) ([]api.IssueLabel, error) {
	panic("not implemented")
}

func (c *fakeClient) ListCycles(_ context.Context, _ string, _ bool, _ int, _ string) (*api.CycleListResult, error) {
	panic("not implemented")
}

func (c *fakeClient) GetCycle(_ context.Context, _ string) (*api.Cycle, error) {
	panic("not implemented")
}

func (c *fakeClient) GetCycleByNumber(_ context.Context, _ string, _ int) (*api.Cycle, error) {
	panic("not implemented")
}

func (c *fakeClient) GetActiveCycle(_ context.Context, _ string) (*api.Cycle, error) {
	panic("not implemented")
}

func (c *fakeClient) RemoveIssueCycle(_ context.Context, _ string) error {
	panic("not implemented")
}

func (c *fakeClient) ListProjects(_ context.Context, _, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	panic("not implemented")
}

func (c *fakeClient) GetProject(_ context.Context, _ string) (*api.Project, error) {
	panic("not implemented")
}

func (c *fakeClient) ListProjectIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	panic("not implemented")
}

func (c *fakeClient) ListInitiatives(_ context.Context, _ string, _ int, _ string) (*api.InitiativeListResult, error) {
	panic("not implemented")
}

func (c *fakeClient) GetInitiative(_ context.Context, _ string) (*api.Initiative, error) {
	panic("not implemented")
}

func (c *fakeClient) ListInitiativeProjects(_ context.Context, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	panic("not implemented")
}

func (c *fakeClient) ListCycleIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	panic("not implemented")
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

func TestResolveTeamID_FromFlagUUID(t *testing.T) {
	f := &Factory{
		IO:      ui.NewTestIOStreams(),
		TeamKey: "550e8400-e29b-41d4-a716-446655440000",
	}

	id, err := ResolveTeamID(context.Background(), f)
	require.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id)
}

func TestResolveTeamID_FromFlagKey(t *testing.T) {
	fc := &fakeClient{
		teams: []api.Team{
			{ID: "team-uuid-1", Key: "ENG", Name: "Engineering"},
			{ID: "team-uuid-2", Key: "DES", Name: "Design"},
		},
	}

	f := &Factory{
		IO:      ui.NewTestIOStreams(),
		TeamKey: "ENG",
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	id, err := ResolveTeamID(context.Background(), f)
	require.NoError(t, err)
	assert.Equal(t, "team-uuid-1", id)
}

func TestResolveTeamID_FromConfig(t *testing.T) {
	dir := t.TempDir()
	store := config.NewViperStore(dir)
	require.NoError(t, store.Load())
	store.SetTeamID("config-team-uuid")

	f := &Factory{
		IO: ui.NewTestIOStreams(),
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	id, err := ResolveTeamID(context.Background(), f)
	require.NoError(t, err)
	assert.Equal(t, "config-team-uuid", id)
}

func TestResolveTeamID_NoneSet(t *testing.T) {
	dir := t.TempDir()
	store := config.NewViperStore(dir)
	require.NoError(t, store.Load())

	f := &Factory{
		IO: ui.NewTestIOStreams(),
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	_, err := ResolveTeamID(context.Background(), f)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no team set")
}
