package issue

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

type fakeClient struct {
	issues         *api.IssueListResult
	issue          *api.Issue
	comment        *api.Comment
	states         []api.WorkflowState
	labels         []api.IssueLabel
	comments       []api.Comment
	commentPI      api.PageInfo
	createInput    api.IssueCreateInput
	updateInput    api.IssueUpdateInput
	updateID       string
	archiveID      string
	addLabelArgs   [2]string
	rmLabelArgs    [2]string
	uploadResult   *api.UploadResult
	uploadedURL    string
	downloadBody   io.ReadCloser
	commentIssueID string
	commentBody    string
	err            error
}

func (c *fakeClient) ListTeams(_ context.Context) ([]api.Team, error) {
	return nil, nil
}

func (c *fakeClient) Viewer(_ context.Context) (*api.User, error) {
	return nil, nil
}

func (c *fakeClient) ListIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return c.issues, c.err
}

func (c *fakeClient) GetIssue(_ context.Context, _ string) (*api.Issue, error) {
	return c.issue, c.err
}

func (c *fakeClient) SearchIssues(_ context.Context, _, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return c.issues, c.err
}

func (c *fakeClient) CreateIssue(_ context.Context, input api.IssueCreateInput) (*api.Issue, error) {
	c.createInput = input
	return c.issue, c.err
}

func (c *fakeClient) UpdateIssue(_ context.Context, id string, input api.IssueUpdateInput) (*api.Issue, error) {
	c.updateID = id
	c.updateInput = input
	return c.issue, c.err
}

func (c *fakeClient) ArchiveIssue(_ context.Context, id string) error {
	c.archiveID = id
	return c.err
}

func (c *fakeClient) ListWorkflowStates(_ context.Context, _ string) ([]api.WorkflowState, error) {
	return c.states, c.err
}

func (c *fakeClient) ListComments(_ context.Context, _ string, _ int, _ string) ([]api.Comment, api.PageInfo, error) {
	return c.comments, c.commentPI, c.err
}

func (c *fakeClient) CreateComment(_ context.Context, issueID, body string) (*api.Comment, error) {
	c.commentIssueID = issueID
	c.commentBody = body
	return c.comment, c.err
}

func (c *fakeClient) AddIssueLabel(_ context.Context, issueID, labelID string) error {
	c.addLabelArgs = [2]string{issueID, labelID}
	return c.err
}

func (c *fakeClient) RemoveIssueLabel(_ context.Context, issueID, labelID string) error {
	c.rmLabelArgs = [2]string{issueID, labelID}
	return c.err
}

func (c *fakeClient) ListLabels(_ context.Context, _ string) ([]api.IssueLabel, error) {
	return c.labels, c.err
}

func (c *fakeClient) ListCycles(_ context.Context, _ string, _ bool, _ int, _ string) (*api.CycleListResult, error) {
	return nil, nil
}

func (c *fakeClient) GetCycle(_ context.Context, _ string) (*api.Cycle, error) {
	return nil, nil
}

func (c *fakeClient) GetCycleByNumber(_ context.Context, _ string, _ int) (*api.Cycle, error) {
	return nil, nil
}

func (c *fakeClient) GetActiveCycle(_ context.Context, _ string) (*api.Cycle, error) {
	return nil, nil
}

func (c *fakeClient) RemoveIssueCycle(_ context.Context, _ string) error {
	return nil
}

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
	return c.uploadResult, c.err
}

func (c *fakeClient) UploadToURL(_ context.Context, url string, _ []api.UploadHeader, _ io.Reader) error {
	c.uploadedURL = url
	return c.err
}

func (c *fakeClient) DownloadURL(_ context.Context, _ string) (io.ReadCloser, error) {
	return c.downloadBody, c.err
}

func TestListCmd_ShowsIssues(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		issues: &api.IssueListResult{
			Issues: []api.Issue{
				{Identifier: "ENG-1", Title: "Fix bug", State: api.WorkflowState{Name: "In Progress"}, PriorityLabel: "Urgent", CreatedAt: now},
				{Identifier: "ENG-2", Title: "Add feature", State: api.WorkflowState{Name: "Backlog"}, PriorityLabel: "Normal", CreatedAt: now},
			},
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

	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "ENG-1")
	assert.Contains(t, buf.String(), "Fix bug")
	assert.Contains(t, buf.String(), "ENG-2")
	assert.Contains(t, buf.String(), "Add feature")
}

func TestListCmd_NoIssues(t *testing.T) {
	fc := &fakeClient{
		issues: &api.IssueListResult{Issues: []api.Issue{}},
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

	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "No issues found")
}

func TestListCmd_FilterByState(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		issues: &api.IssueListResult{
			Issues: []api.Issue{
				{Identifier: "ENG-1", Title: "Active", State: api.WorkflowState{Name: "In Progress", Type: "started"}, CreatedAt: now},
				{Identifier: "ENG-2", Title: "Waiting", State: api.WorkflowState{Name: "Backlog", Type: "backlog"}, CreatedAt: now},
			},
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

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--state", "started"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "ENG-1")
	assert.NotContains(t, buf.String(), "ENG-2")
}
