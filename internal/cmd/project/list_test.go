package project

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

type fakeClient struct {
	projects     *api.ProjectListResult
	projectPages []*api.ProjectListResult
	projectPage  int
	project      *api.Project
	issues       *api.IssueListResult
	teamFilter   string
	statusFilter string
	err          error
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
	return nil, nil
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

func (c *fakeClient) ListProjects(_ context.Context, teamID, status string, _ int, _ string) (*api.ProjectListResult, error) {
	c.teamFilter = teamID
	c.statusFilter = status
	if len(c.projectPages) > 0 {
		page := c.projectPages[c.projectPage]
		c.projectPage++
		return page, c.err
	}
	return c.projects, c.err
}

func (c *fakeClient) GetProject(_ context.Context, _ string) (*api.Project, error) {
	return c.project, c.err
}

func (c *fakeClient) ListProjectIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return c.issues, c.err
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

func newTestFactory(fc *fakeClient) *cmdutil.Factory {
	ios := ui.NewTestIOStreams()
	dir, _ := os.MkdirTemp("", "lnr-test")
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

func TestListCmd_ShowsProjects(t *testing.T) {
	fc := &fakeClient{
		projects: &api.ProjectListResult{
			Projects: []api.Project{
				{Name: "Auth Rewrite", Status: api.ProjectStatus{Type: "started"}, Progress: 0.5},
				{Name: "Dark Mode", Status: api.ProjectStatus{Type: "planned"}, Progress: 0},
			},
		},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"list"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Auth Rewrite")
	assert.Contains(t, out, "Dark Mode")
}

func TestListCmd_NoProjects(t *testing.T) {
	fc := &fakeClient{
		projects: &api.ProjectListResult{Projects: nil},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"list"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "No projects found")
}

func TestListCmd_WithStatusFilter(t *testing.T) {
	fc := &fakeClient{
		projects: &api.ProjectListResult{Projects: nil},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"list", "--status", "started"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "started", fc.statusFilter)
}

func TestListCmd_PassesTeamID(t *testing.T) {
	fc := &fakeClient{
		projects: &api.ProjectListResult{Projects: nil},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"list"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "team-1", fc.teamFilter)
}

func TestListCmd_Pagination(t *testing.T) {
	fc := &fakeClient{
		projectPages: []*api.ProjectListResult{
			{
				Projects: []api.Project{{Name: "Project A"}},
				PageInfo: api.PageInfo{HasNextPage: true, EndCursor: "cursor-1"},
			},
			{
				Projects: []api.Project{{Name: "Project B"}},
				PageInfo: api.PageInfo{HasNextPage: false},
			},
		},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"list"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Project A")
	assert.Contains(t, out, "Project B")
}
