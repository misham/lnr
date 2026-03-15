package cmd

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

type fakeClient struct {
	viewer *api.User
	teams  []api.Team
	err    error
}

func (c *fakeClient) ListTeams(_ context.Context) ([]api.Team, error) {
	return c.teams, c.err
}

func (c *fakeClient) Viewer(_ context.Context) (*api.User, error) {
	return c.viewer, c.err
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

func (c *fakeClient) ListCycleIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
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

func (c *fakeClient) FileUpload(_ context.Context, _, _ string, _ int64) (*api.UploadResult, error) {
	return nil, nil
}

func (c *fakeClient) UploadToURL(_ context.Context, _ string, _ []api.UploadHeader, _ io.Reader) error {
	return nil
}

func (c *fakeClient) DownloadURL(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, nil
}

func TestMeCmd_ShowsUser(t *testing.T) {
	ios := ui.NewTestIOStreams()

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{
				viewer: &api.User{
					ID:          "user-1",
					Name:        "Jane Doe",
					DisplayName: "janed",
					Email:       "jane@example.com",
					Active:      true,
				},
			}, nil
		},
	}

	cmd := newMeCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, ok := ios.Out.(*bytes.Buffer)
	require.True(t, ok)
	assert.Contains(t, buf.String(), "Jane Doe")
	assert.Contains(t, buf.String(), "jane@example.com")
	assert.Contains(t, buf.String(), "janed")
}
