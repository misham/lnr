package initiative

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
	initiatives  *api.InitiativeListResult
	initiative   *api.Initiative
	projects     *api.ProjectListResult
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

func (c *fakeClient) ListProjects(_ context.Context, _, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return nil, nil
}

func (c *fakeClient) GetProject(_ context.Context, _ string) (*api.Project, error) {
	return nil, nil
}

func (c *fakeClient) ListProjectIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeClient) ListInitiatives(_ context.Context, status string, _ int, _ string) (*api.InitiativeListResult, error) {
	c.statusFilter = status
	return c.initiatives, c.err
}

func (c *fakeClient) GetInitiative(_ context.Context, _ string) (*api.Initiative, error) {
	return c.initiative, c.err
}

func (c *fakeClient) ListInitiativeProjects(_ context.Context, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return c.projects, c.err
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

func TestListCmd_ShowsInitiatives(t *testing.T) {
	fc := &fakeClient{
		initiatives: &api.InitiativeListResult{
			Initiatives: []api.Initiative{
				{Name: "Platform Reliability", Status: "Active"},
				{Name: "Cost Reduction", Status: "Planned"},
			},
		},
	}
	f := newTestFactory(fc)

	cmd := NewInitiativeCmd(f)
	cmd.SetArgs([]string{"list"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Platform Reliability")
	assert.Contains(t, out, "Cost Reduction")
}

func TestListCmd_NoInitiatives(t *testing.T) {
	fc := &fakeClient{
		initiatives: &api.InitiativeListResult{Initiatives: nil},
	}
	f := newTestFactory(fc)

	cmd := NewInitiativeCmd(f)
	cmd.SetArgs([]string{"list"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "No initiatives found")
}

func TestListCmd_WithStatusFilter(t *testing.T) {
	fc := &fakeClient{
		initiatives: &api.InitiativeListResult{Initiatives: nil},
	}
	f := newTestFactory(fc)

	cmd := NewInitiativeCmd(f)
	cmd.SetArgs([]string{"list", "--status", "active"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "Active", fc.statusFilter)
}
