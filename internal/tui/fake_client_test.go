package tui

import (
	"context"
	"io"

	"github.com/misham/linear-cli/internal/api"
)

// fakeClient implements api.Client for testing.
type fakeClient struct {
	issues       *api.IssueListResult
	issue        *api.Issue
	cycles       *api.CycleListResult
	cycle        *api.Cycle
	projects     *api.ProjectListResult
	project      *api.Project
	initiatives  *api.InitiativeListResult
	initiative   *api.Initiative
	initProjects *api.ProjectListResult
	err          error
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

func (c *fakeClient) CreateIssue(_ context.Context, _ api.IssueCreateInput) (*api.Issue, error) {
	return nil, nil
}

func (c *fakeClient) UpdateIssue(_ context.Context, _ string, _ api.IssueUpdateInput) (*api.Issue, error) {
	return nil, nil
}

func (c *fakeClient) ArchiveIssue(_ context.Context, _ string) error {
	return nil
}

func (c *fakeClient) ListWorkflowStates(_ context.Context, _ string) ([]api.WorkflowState, error) {
	return nil, nil
}

func (c *fakeClient) ListComments(_ context.Context, _ string, _ int, _ string) ([]api.Comment, api.PageInfo, error) {
	return nil, api.PageInfo{}, nil
}

func (c *fakeClient) CreateComment(_ context.Context, _, _ string) (*api.Comment, error) {
	return nil, nil
}

func (c *fakeClient) AddIssueLabel(_ context.Context, _, _ string) error {
	return nil
}

func (c *fakeClient) RemoveIssueLabel(_ context.Context, _, _ string) error {
	return nil
}

func (c *fakeClient) ListLabels(_ context.Context, _ string) ([]api.IssueLabel, error) {
	return nil, nil
}

func (c *fakeClient) ListCycles(_ context.Context, _ string, _ bool, _ int, _ string) (*api.CycleListResult, error) {
	return c.cycles, c.err
}

func (c *fakeClient) GetCycle(_ context.Context, _ string) (*api.Cycle, error) {
	return c.cycle, c.err
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
	return c.issues, c.err
}

func (c *fakeClient) ListProjects(_ context.Context, _, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return c.projects, c.err
}

func (c *fakeClient) GetProject(_ context.Context, _ string) (*api.Project, error) {
	return c.project, c.err
}

func (c *fakeClient) ListProjectIssues(_ context.Context, _ string, _ int, _ string) (*api.IssueListResult, error) {
	return nil, nil
}

func (c *fakeClient) ListInitiatives(_ context.Context, _ string, _ int, _ string) (*api.InitiativeListResult, error) {
	return c.initiatives, c.err
}

func (c *fakeClient) GetInitiative(_ context.Context, _ string) (*api.Initiative, error) {
	return c.initiative, c.err
}

func (c *fakeClient) ListInitiativeProjects(_ context.Context, _ string, _ int, _ string) (*api.ProjectListResult, error) {
	return c.initProjects, c.err
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
