package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/misham/linear-cli/internal/api/linear_graphql"
)

// ErrNoActiveCycle is returned when no active cycle exists for a team.
var ErrNoActiveCycle = errors.New("no active cycle")

// Client defines the interface for interacting with the Linear API.
type Client interface {
	ListTeams(ctx context.Context) ([]Team, error)
	Viewer(ctx context.Context) (*User, error)
	ListIssues(ctx context.Context, teamID string, first int, after string) (*IssueListResult, error)
	GetIssue(ctx context.Context, id string) (*Issue, error)
	SearchIssues(ctx context.Context, term, teamID string, first int, after string) (*IssueListResult, error)
	CreateIssue(ctx context.Context, input IssueCreateInput) (*Issue, error)
	UpdateIssue(ctx context.Context, id string, input IssueUpdateInput) (*Issue, error)
	ArchiveIssue(ctx context.Context, id string) error
	ListWorkflowStates(ctx context.Context, teamID string) ([]WorkflowState, error)
	ListComments(ctx context.Context, issueID string, first int, after string) ([]Comment, PageInfo, error)
	CreateComment(ctx context.Context, issueID, body string) (*Comment, error)
	AddIssueLabel(ctx context.Context, issueID, labelID string) error
	RemoveIssueLabel(ctx context.Context, issueID, labelID string) error
	ListLabels(ctx context.Context, teamID string) ([]IssueLabel, error)
	ListCycles(ctx context.Context, teamID string, includeAll bool, first int, after string) (*CycleListResult, error)
	GetCycle(ctx context.Context, id string) (*Cycle, error)
	GetCycleByNumber(ctx context.Context, teamID string, number int) (*Cycle, error)
	GetActiveCycle(ctx context.Context, teamID string) (*Cycle, error)
	ListCycleIssues(ctx context.Context, cycleID string, first int, after string) (*IssueListResult, error)
	RemoveIssueCycle(ctx context.Context, issueID string) error
	ListProjects(ctx context.Context, teamID, statusType string, first int, after string) (*ProjectListResult, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	ListProjectIssues(ctx context.Context, projectID string, first int, after string) (*IssueListResult, error)
	ListInitiatives(ctx context.Context, status string, first int, after string) (*InitiativeListResult, error)
	GetInitiative(ctx context.Context, id string) (*Initiative, error)
	ListInitiativeProjects(ctx context.Context, initiativeID string, first int, after string) (*ProjectListResult, error)
	FileUpload(ctx context.Context, contentType, filename string, size int64) (*UploadResult, error)
	UploadToURL(ctx context.Context, url string, headers []UploadHeader, body io.Reader) error
	DownloadURL(ctx context.Context, url string) (io.ReadCloser, error)
}

type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req)
}

// GraphQLClient implements Client using genqlient.
type GraphQLClient struct {
	gql        graphql.Client
	httpClient *http.Client
}

// NewGraphQLClient creates a Client that talks to the given GraphQL endpoint.
func NewGraphQLClient(endpoint, token string) *GraphQLClient {
	httpClient := &http.Client{
		Transport: &authTransport{
			token: token,
			base:  http.DefaultTransport,
		},
	}
	return &GraphQLClient{
		gql:        graphql.NewClient(endpoint, httpClient),
		httpClient: httpClient,
	}
}

func (c *GraphQLClient) ListTeams(ctx context.Context) ([]Team, error) {
	resp, err := linear_graphql.ListTeams(ctx, c.gql)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	teams := make([]Team, len(resp.Teams.Nodes))
	for i, n := range resp.Teams.Nodes {
		teams[i] = Team{
			ID:          n.Id,
			Name:        n.Name,
			Key:         n.Key,
			Description: derefStr(n.Description),
			Private:     n.Private,
			Icon:        derefStr(n.Icon),
			Color:       derefStr(n.Color),
			Timezone:    n.Timezone,
		}
	}
	return teams, nil
}

func (c *GraphQLClient) Viewer(ctx context.Context) (*User, error) {
	resp, err := linear_graphql.Viewer(ctx, c.gql)
	if err != nil {
		return nil, fmt.Errorf("viewer: %w", err)
	}

	return &User{
		ID:          resp.Viewer.Id,
		Name:        resp.Viewer.Name,
		DisplayName: resp.Viewer.DisplayName,
		Email:       resp.Viewer.Email,
		Active:      resp.Viewer.Active,
	}, nil
}

func (c *GraphQLClient) ListIssues(ctx context.Context, teamID string, first int, after string) (*IssueListResult, error) {
	var filter *linear_graphql.IssueFilter
	if teamID != "" {
		filter = &linear_graphql.IssueFilter{
			Team: &linear_graphql.TeamFilter{
				Id: &linear_graphql.IDComparator{Eq: &teamID},
			},
		}
	}
	resp, err := linear_graphql.ListIssues(ctx, c.gql, filter, &first, strPtrOrNil(after))
	if err != nil {
		return nil, fmt.Errorf("list issues: %w", err)
	}

	issues := make([]Issue, len(resp.Issues.Nodes))
	for i, n := range resp.Issues.Nodes {
		labels := make([]IssueLabel, len(n.Labels.Nodes))
		for j, l := range n.Labels.Nodes {
			labels[j] = IssueLabel{ID: l.Id, Name: l.Name, Color: l.Color}
		}
		var assignee *User
		if n.Assignee != nil {
			assignee = &User{ID: n.Assignee.Id, Name: n.Assignee.Name, DisplayName: n.Assignee.DisplayName, Email: n.Assignee.Email, Active: n.Assignee.Active}
		}
		issue := Issue{
			ID: n.Id, Identifier: n.Identifier, Title: n.Title, Description: derefStr(n.Description),
			Priority: int(n.Priority), PriorityLabel: n.PriorityLabel,
			Estimate: derefIntFromFloat(n.Estimate), DueDate: derefStr(n.DueDate), URL: n.Url,
			CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
			CompletedAt: n.CompletedAt, ArchivedAt: n.ArchivedAt,
			State:    WorkflowState{ID: n.State.Id, Name: n.State.Name, Type: n.State.Type, Color: n.State.Color, Position: n.State.Position},
			Team:     Team{ID: n.Team.Id, Name: n.Team.Name, Key: n.Team.Key},
			Assignee: assignee,
			Labels:   labels,
		}
		if n.Cycle != nil {
			issue.Cycle = &IssueCycle{ID: n.Cycle.Id, Name: derefStr(n.Cycle.Name), Number: int(n.Cycle.Number)}
		}
		if n.Project != nil {
			issue.Project = &IssueProject{ID: n.Project.Id, Name: n.Project.Name}
		}
		issues[i] = issue
	}
	return &IssueListResult{
		Issues: issues,
		PageInfo: PageInfo{
			HasNextPage: resp.Issues.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.Issues.PageInfo.EndCursor),
		},
	}, nil
}

func (c *GraphQLClient) GetIssue(ctx context.Context, id string) (*Issue, error) {
	resp, err := linear_graphql.GetIssue(ctx, c.gql, id)
	if err != nil {
		return nil, fmt.Errorf("get issue: %w", err)
	}
	n := resp.Issue
	labels := make([]IssueLabel, len(n.Labels.Nodes))
	for j, l := range n.Labels.Nodes {
		labels[j] = IssueLabel{ID: l.Id, Name: l.Name, Color: l.Color}
	}
	comments := make([]Comment, len(n.Comments.Nodes))
	for j, cm := range n.Comments.Nodes {
		var u *User
		if cm.User != nil {
			u = &User{ID: cm.User.Id, Name: cm.User.Name, DisplayName: cm.User.DisplayName, Email: cm.User.Email, Active: cm.User.Active}
		}
		comments[j] = Comment{
			ID: cm.Id, Body: cm.Body,
			CreatedAt: cm.CreatedAt, UpdatedAt: cm.UpdatedAt,
			User: u,
		}
	}
	var assignee *User
	if n.Assignee != nil {
		assignee = &User{ID: n.Assignee.Id, Name: n.Assignee.Name, DisplayName: n.Assignee.DisplayName, Email: n.Assignee.Email, Active: n.Assignee.Active}
	}
	issue := &Issue{
		ID: n.Id, Identifier: n.Identifier, Title: n.Title, Description: derefStr(n.Description),
		Priority: int(n.Priority), PriorityLabel: n.PriorityLabel,
		Estimate: derefIntFromFloat(n.Estimate), DueDate: derefStr(n.DueDate), URL: n.Url,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
		CompletedAt: n.CompletedAt, ArchivedAt: n.ArchivedAt,
		State:    WorkflowState{ID: n.State.Id, Name: n.State.Name, Type: n.State.Type, Color: n.State.Color, Position: n.State.Position},
		Team:     Team{ID: n.Team.Id, Name: n.Team.Name, Key: n.Team.Key},
		Assignee: assignee,
		Labels:   labels,
		Comments: comments,
	}
	if n.Cycle != nil {
		issue.Cycle = &IssueCycle{ID: n.Cycle.Id, Name: derefStr(n.Cycle.Name), Number: int(n.Cycle.Number)}
	}
	if n.Project != nil {
		issue.Project = &IssueProject{ID: n.Project.Id, Name: n.Project.Name}
	}
	return issue, nil
}

func (c *GraphQLClient) SearchIssues(ctx context.Context, term, teamID string, first int, after string) (*IssueListResult, error) {
	resp, err := linear_graphql.SearchIssues(ctx, c.gql, term, &linear_graphql.IssueFilter{}, &first, strPtrOrNil(after), strPtrOrNil(teamID))
	if err != nil {
		return nil, fmt.Errorf("search issues: %w", err)
	}

	issues := make([]Issue, len(resp.SearchIssues.Nodes))
	for i, n := range resp.SearchIssues.Nodes {
		labels := make([]IssueLabel, len(n.Labels.Nodes))
		for j, l := range n.Labels.Nodes {
			labels[j] = IssueLabel{ID: l.Id, Name: l.Name, Color: l.Color}
		}
		var assignee *User
		if n.Assignee != nil {
			assignee = &User{ID: n.Assignee.Id, Name: n.Assignee.Name, DisplayName: n.Assignee.DisplayName, Email: n.Assignee.Email, Active: n.Assignee.Active}
		}
		issue := Issue{
			ID: n.Id, Identifier: n.Identifier, Title: n.Title, Description: derefStr(n.Description),
			Priority: int(n.Priority), PriorityLabel: n.PriorityLabel,
			Estimate: derefIntFromFloat(n.Estimate), DueDate: derefStr(n.DueDate), URL: n.Url,
			CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
			CompletedAt: n.CompletedAt, ArchivedAt: n.ArchivedAt,
			State:    WorkflowState{ID: n.State.Id, Name: n.State.Name, Type: n.State.Type, Color: n.State.Color, Position: n.State.Position},
			Team:     Team{ID: n.Team.Id, Name: n.Team.Name, Key: n.Team.Key},
			Assignee: assignee,
			Labels:   labels,
		}
		if n.Cycle != nil {
			issue.Cycle = &IssueCycle{ID: n.Cycle.Id, Name: derefStr(n.Cycle.Name), Number: int(n.Cycle.Number)}
		}
		if n.Project != nil {
			issue.Project = &IssueProject{ID: n.Project.Id, Name: n.Project.Name}
		}
		issues[i] = issue
	}
	return &IssueListResult{
		Issues: issues,
		PageInfo: PageInfo{
			HasNextPage: resp.SearchIssues.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.SearchIssues.PageInfo.EndCursor),
		},
	}, nil
}

func (c *GraphQLClient) CreateIssue(ctx context.Context, input IssueCreateInput) (*Issue, error) {
	genInput := linear_graphql.IssueCreateInput{
		Title:  &input.Title,
		TeamId: input.TeamID,
	}
	if input.Description != "" {
		genInput.Description = &input.Description
	}
	if input.Priority != 0 {
		genInput.Priority = &input.Priority
	}
	if input.StateID != "" {
		genInput.StateId = &input.StateID
	}
	if input.AssigneeID != "" {
		genInput.AssigneeId = &input.AssigneeID
	}
	if len(input.LabelIDs) > 0 {
		genInput.LabelIds = input.LabelIDs
	}
	if input.Estimate != 0 {
		genInput.Estimate = &input.Estimate
	}
	if input.DueDate != "" {
		genInput.DueDate = &input.DueDate
	}

	resp, err := linear_graphql.CreateIssue(ctx, c.gql, genInput)
	if err != nil {
		return nil, fmt.Errorf("create issue: %w", err)
	}
	n := resp.IssueCreate.Issue
	return &Issue{
		ID: n.Id, Identifier: n.Identifier, Title: n.Title, URL: n.Url,
		State: WorkflowState{ID: n.State.Id, Name: n.State.Name, Type: n.State.Type},
	}, nil
}

func (c *GraphQLClient) UpdateIssue(ctx context.Context, id string, input IssueUpdateInput) (*Issue, error) {
	genInput := linear_graphql.IssueUpdateInput{}
	if input.Title != nil {
		genInput.Title = input.Title
	}
	if input.Description != nil {
		genInput.Description = input.Description
	}
	if input.Priority != nil {
		genInput.Priority = input.Priority
	}
	if input.StateID != nil {
		genInput.StateId = input.StateID
	}
	if input.AssigneeID != nil {
		genInput.AssigneeId = input.AssigneeID
	}
	if len(input.LabelIDs) > 0 {
		genInput.LabelIds = input.LabelIDs
	}
	if input.Estimate != nil {
		genInput.Estimate = input.Estimate
	}
	if input.DueDate != nil {
		genInput.DueDate = input.DueDate
	}
	if input.CycleID != nil {
		genInput.CycleId = input.CycleID
	}

	resp, err := linear_graphql.UpdateIssue(ctx, c.gql, id, genInput)
	if err != nil {
		return nil, fmt.Errorf("update issue: %w", err)
	}
	n := resp.IssueUpdate.Issue
	return &Issue{
		ID: n.Id, Identifier: n.Identifier, Title: n.Title, URL: n.Url,
		State: WorkflowState{ID: n.State.Id, Name: n.State.Name, Type: n.State.Type},
	}, nil
}

func (c *GraphQLClient) ArchiveIssue(ctx context.Context, id string) error {
	_, err := linear_graphql.ArchiveIssue(ctx, c.gql, id)
	if err != nil {
		return fmt.Errorf("archive issue: %w", err)
	}
	return nil
}

func (c *GraphQLClient) ListWorkflowStates(ctx context.Context, teamID string) ([]WorkflowState, error) {
	var filter *linear_graphql.WorkflowStateFilter
	if teamID != "" {
		filter = &linear_graphql.WorkflowStateFilter{
			Team: &linear_graphql.TeamFilter{
				Id: &linear_graphql.IDComparator{Eq: &teamID},
			},
		}
	}
	resp, err := linear_graphql.ListWorkflowStates(ctx, c.gql, filter)
	if err != nil {
		return nil, fmt.Errorf("list workflow states: %w", err)
	}
	states := make([]WorkflowState, len(resp.WorkflowStates.Nodes))
	for i, n := range resp.WorkflowStates.Nodes {
		states[i] = WorkflowState{ID: n.Id, Name: n.Name, Type: n.Type, Color: n.Color, Position: n.Position}
	}
	return states, nil
}

func (c *GraphQLClient) ListComments(ctx context.Context, issueID string, first int, after string) ([]Comment, PageInfo, error) {
	resp, err := linear_graphql.ListComments(ctx, c.gql, issueID, &first, strPtrOrNil(after))
	if err != nil {
		return nil, PageInfo{}, fmt.Errorf("list comments: %w", err)
	}
	nodes := resp.Issue.Comments.Nodes
	comments := make([]Comment, len(nodes))
	for i, n := range nodes {
		var u *User
		if n.User != nil {
			u = &User{ID: n.User.Id, Name: n.User.Name, DisplayName: n.User.DisplayName, Email: n.User.Email, Active: n.User.Active}
		}
		comments[i] = Comment{
			ID: n.Id, Body: n.Body,
			CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
			User: u,
		}
	}
	pi := resp.Issue.Comments.PageInfo
	return comments, PageInfo{HasNextPage: pi.HasNextPage, EndCursor: derefStr(pi.EndCursor)}, nil
}

func (c *GraphQLClient) CreateComment(ctx context.Context, issueID, body string) (*Comment, error) {
	resp, err := linear_graphql.CreateComment(ctx, c.gql, issueID, body)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	n := resp.CommentCreate.Comment
	var u *User
	if n.User != nil {
		u = &User{ID: n.User.Id, Name: n.User.Name, DisplayName: n.User.DisplayName, Email: n.User.Email, Active: n.User.Active}
	}
	return &Comment{
		ID: n.Id, Body: n.Body, CreatedAt: n.CreatedAt,
		User: u,
	}, nil
}

func (c *GraphQLClient) AddIssueLabel(ctx context.Context, issueID, labelID string) error {
	_, err := linear_graphql.AddIssueLabel(ctx, c.gql, issueID, labelID)
	if err != nil {
		return fmt.Errorf("add issue label: %w", err)
	}
	return nil
}

func (c *GraphQLClient) RemoveIssueLabel(ctx context.Context, issueID, labelID string) error {
	_, err := linear_graphql.RemoveIssueLabel(ctx, c.gql, issueID, labelID)
	if err != nil {
		return fmt.Errorf("remove issue label: %w", err)
	}
	return nil
}

func (c *GraphQLClient) ListLabels(ctx context.Context, teamID string) ([]IssueLabel, error) {
	// Fetch all labels (workspace + team) since both are valid for any issue.
	resp, err := linear_graphql.ListLabels(ctx, c.gql, nil)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	labels := make([]IssueLabel, len(resp.IssueLabels.Nodes))
	for i, n := range resp.IssueLabels.Nodes {
		labels[i] = IssueLabel{ID: n.Id, Name: n.Name, Color: n.Color}
	}
	return labels, nil
}

func (c *GraphQLClient) ListCycles(ctx context.Context, teamID string, includeAll bool, first int, after string) (*CycleListResult, error) {
	filter := linear_graphql.CycleFilter{
		Team: &linear_graphql.TeamFilter{
			Id: &linear_graphql.IDComparator{Eq: &teamID},
		},
	}
	if !includeAll {
		notPast := false
		filter.IsPast = &linear_graphql.BooleanComparator{Eq: &notPast}
	}
	resp, err := linear_graphql.ListCycles(ctx, c.gql, &filter, &first, strPtrOrNil(after))
	if err != nil {
		return nil, fmt.Errorf("list cycles: %w", err)
	}

	cycles := make([]Cycle, len(resp.Cycles.Nodes))
	for i, n := range resp.Cycles.Nodes {
		cycles[i] = cycleFromListNode(n)
	}
	return &CycleListResult{
		Cycles: cycles,
		PageInfo: PageInfo{
			HasNextPage: resp.Cycles.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.Cycles.PageInfo.EndCursor),
		},
	}, nil
}

func (c *GraphQLClient) GetCycle(ctx context.Context, id string) (*Cycle, error) {
	resp, err := linear_graphql.GetCycle(ctx, c.gql, id)
	if err != nil {
		return nil, fmt.Errorf("get cycle: %w", err)
	}
	n := resp.Cycle
	return &Cycle{
		ID: n.Id, Name: derefStr(n.Name), Number: int(n.Number),
		Description: derefStr(n.Description),
		StartsAt:    n.StartsAt, EndsAt: n.EndsAt, CompletedAt: n.CompletedAt,
		Progress: n.Progress, IsActive: n.IsActive, IsNext: n.IsNext, IsPast: n.IsPast,
		Team:      Team{ID: n.Team.Id, Name: n.Team.Name, Key: n.Team.Key},
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}, nil
}

func (c *GraphQLClient) GetCycleByNumber(ctx context.Context, teamID string, number int) (*Cycle, error) {
	numberFloat := float64(number)
	filter := linear_graphql.CycleFilter{
		Team: &linear_graphql.TeamFilter{
			Id: &linear_graphql.IDComparator{Eq: &teamID},
		},
		Number: &linear_graphql.NumberComparator{Eq: &numberFloat},
	}
	first := 1
	resp, err := linear_graphql.ListCycles(ctx, c.gql, &filter, &first, nil)
	if err != nil {
		return nil, fmt.Errorf("get cycle by number: %w", err)
	}
	if len(resp.Cycles.Nodes) == 0 {
		return nil, fmt.Errorf("no cycle with number %d found", number)
	}
	cy := cycleFromListNode(resp.Cycles.Nodes[0])
	return &cy, nil
}

func (c *GraphQLClient) GetActiveCycle(ctx context.Context, teamID string) (*Cycle, error) {
	active := true
	filter := linear_graphql.CycleFilter{
		Team: &linear_graphql.TeamFilter{
			Id: &linear_graphql.IDComparator{Eq: &teamID},
		},
		IsActive: &linear_graphql.BooleanComparator{Eq: &active},
	}
	first := 1
	resp, err := linear_graphql.ListCycles(ctx, c.gql, &filter, &first, nil)
	if err != nil {
		return nil, fmt.Errorf("get active cycle: %w", err)
	}
	if len(resp.Cycles.Nodes) == 0 {
		return nil, ErrNoActiveCycle
	}
	cy := cycleFromListNode(resp.Cycles.Nodes[0])
	return &cy, nil
}

func (c *GraphQLClient) ListCycleIssues(ctx context.Context, cycleID string, first int, after string) (*IssueListResult, error) {
	resp, err := linear_graphql.ListCycleIssues(ctx, c.gql, cycleID, &first, strPtrOrNil(after))
	if err != nil {
		return nil, fmt.Errorf("list cycle issues: %w", err)
	}

	nodes := resp.Cycle.Issues.Nodes
	issues := make([]Issue, len(nodes))
	for i, n := range nodes {
		var assignee *User
		if n.Assignee != nil {
			assignee = &User{ID: n.Assignee.Id, Name: n.Assignee.Name, DisplayName: n.Assignee.DisplayName, Email: n.Assignee.Email, Active: n.Assignee.Active}
		}
		issues[i] = Issue{
			ID: n.Id, Identifier: n.Identifier, Title: n.Title,
			Priority: int(n.Priority), PriorityLabel: n.PriorityLabel,
			State:    WorkflowState{ID: n.State.Id, Name: n.State.Name, Type: n.State.Type, Color: n.State.Color, Position: n.State.Position},
			Assignee: assignee,
		}
	}
	return &IssueListResult{
		Issues: issues,
		PageInfo: PageInfo{
			HasNextPage: resp.Cycle.Issues.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.Cycle.Issues.PageInfo.EndCursor),
		},
	}, nil
}

// issueUpdateNullCycle is a minimal input that explicitly sends cycleId: null.
// The generated IssueUpdateInput uses omitempty, which omits nil pointers entirely.
// Linear requires an explicit null to unset the cycle.
type issueUpdateNullCycle struct {
	CycleID interface{} `json:"cycleId"`
}

func (c *GraphQLClient) RemoveIssueCycle(ctx context.Context, issueID string) error {
	input := issueUpdateNullCycle{CycleID: nil}
	req := &graphql.Request{
		OpName: "RemoveIssueCycle",
		Query: `mutation RemoveIssueCycle($id: String!, $input: IssueUpdateInput!) {
			issueUpdate(id: $id, input: $input) {
				issue { id }
			}
		}`,
		Variables: map[string]interface{}{
			"id":    issueID,
			"input": input,
		},
	}
	resp := &graphql.Response{Data: &struct{}{}}
	err := c.gql.MakeRequest(ctx, req, resp)
	if err != nil {
		return fmt.Errorf("remove issue cycle: %w", err)
	}
	return nil
}

func cycleFromListNode(n linear_graphql.ListCyclesCyclesCycleConnectionNodesCycle) Cycle {
	return Cycle{
		ID: n.Id, Name: derefStr(n.Name), Number: int(n.Number),
		Description: derefStr(n.Description),
		StartsAt:    n.StartsAt, EndsAt: n.EndsAt, CompletedAt: n.CompletedAt,
		Progress: n.Progress, IsActive: n.IsActive, IsNext: n.IsNext, IsPast: n.IsPast,
		Team:      Team{ID: n.Team.Id, Name: n.Team.Name, Key: n.Team.Key},
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}
}

func (c *GraphQLClient) ListProjects(ctx context.Context, teamID, statusType string, first int, after string) (*ProjectListResult, error) {
	var filter *linear_graphql.ProjectFilter
	if statusType != "" || teamID != "" {
		filter = &linear_graphql.ProjectFilter{}
		if statusType != "" {
			filter.Status = &linear_graphql.ProjectStatusFilter{
				Type: &linear_graphql.StringComparator{Eq: &statusType},
			}
		}
		if teamID != "" {
			filter.AccessibleTeams = &linear_graphql.TeamCollectionFilter{
				Some: &linear_graphql.TeamFilter{
					Id: &linear_graphql.IDComparator{Eq: &teamID},
				},
			}
		}
	}
	resp, err := linear_graphql.ListProjects(ctx, c.gql, filter, &first, strPtrOrNil(after))
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	projects := make([]Project, len(resp.Projects.Nodes))
	for i, n := range resp.Projects.Nodes {
		projects[i] = projectFromListNode(n)
	}
	return &ProjectListResult{
		Projects: projects,
		PageInfo: PageInfo{
			HasNextPage: resp.Projects.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.Projects.PageInfo.EndCursor),
		},
	}, nil
}

func projectFromListNode(n linear_graphql.ListProjectsProjectsProjectConnectionNodesProject) Project {
	var lead *User
	if n.Lead != nil {
		lead = &User{ID: n.Lead.Id, Name: n.Lead.Name, DisplayName: n.Lead.DisplayName, Email: n.Lead.Email, Active: n.Lead.Active}
	}
	return Project{
		ID: n.Id, Name: n.Name, Description: n.Description,
		SlugID: n.SlugId, Color: n.Color, Icon: derefStr(n.Icon),
		Status:   ProjectStatus{Name: n.Status.Name, Color: n.Status.Color, Type: string(n.Status.Type)},
		Progress: n.Progress, Priority: n.Priority, PriorityLabel: n.PriorityLabel,
		StartDate: derefStr(n.StartDate), TargetDate: derefStr(n.TargetDate),
		Lead: lead, URL: n.Url,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
		CompletedAt: n.CompletedAt, CanceledAt: n.CanceledAt, ArchivedAt: n.ArchivedAt,
	}
}

func (c *GraphQLClient) GetProject(ctx context.Context, id string) (*Project, error) {
	resp, err := linear_graphql.GetProject(ctx, c.gql, id)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	n := resp.Project
	var lead *User
	if n.Lead != nil {
		lead = &User{ID: n.Lead.Id, Name: n.Lead.Name, DisplayName: n.Lead.DisplayName, Email: n.Lead.Email, Active: n.Lead.Active}
	}
	milestones := make([]ProjectMilestone, len(n.ProjectMilestones.Nodes))
	for i, m := range n.ProjectMilestones.Nodes {
		milestones[i] = ProjectMilestone{
			ID: m.Id, Name: m.Name, Description: derefStr(m.Description),
			TargetDate: derefStr(m.TargetDate), SortOrder: m.SortOrder,
		}
	}
	sort.Slice(milestones, func(i, j int) bool {
		return milestones[i].SortOrder < milestones[j].SortOrder
	})
	return &Project{
		ID: n.Id, Name: n.Name, Description: n.Description, Content: derefStr(n.Content),
		SlugID: n.SlugId, Color: n.Color, Icon: derefStr(n.Icon),
		Status:   ProjectStatus{Name: n.Status.Name, Color: n.Status.Color, Type: string(n.Status.Type)},
		Progress: n.Progress, Priority: n.Priority, PriorityLabel: n.PriorityLabel,
		StartDate: derefStr(n.StartDate), TargetDate: derefStr(n.TargetDate),
		Lead: lead, URL: n.Url, Milestones: milestones,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
		CompletedAt: n.CompletedAt, CanceledAt: n.CanceledAt, ArchivedAt: n.ArchivedAt,
	}, nil
}

func (c *GraphQLClient) ListProjectIssues(ctx context.Context, projectID string, first int, after string) (*IssueListResult, error) {
	resp, err := linear_graphql.ListProjectIssues(ctx, c.gql, projectID, &first, strPtrOrNil(after))
	if err != nil {
		return nil, fmt.Errorf("list project issues: %w", err)
	}

	nodes := resp.Project.Issues.Nodes
	issues := make([]Issue, len(nodes))
	for i, n := range nodes {
		labels := make([]IssueLabel, len(n.Labels.Nodes))
		for j, l := range n.Labels.Nodes {
			labels[j] = IssueLabel{ID: l.Id, Name: l.Name, Color: l.Color}
		}
		var assignee *User
		if n.Assignee != nil {
			assignee = &User{ID: n.Assignee.Id, Name: n.Assignee.Name, DisplayName: n.Assignee.DisplayName, Email: n.Assignee.Email, Active: n.Assignee.Active}
		}
		issue := Issue{
			ID: n.Id, Identifier: n.Identifier, Title: n.Title, Description: derefStr(n.Description),
			Priority: int(n.Priority), PriorityLabel: n.PriorityLabel,
			Estimate: derefIntFromFloat(n.Estimate), DueDate: derefStr(n.DueDate), URL: n.Url,
			CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
			CompletedAt: n.CompletedAt, ArchivedAt: n.ArchivedAt,
			State:    WorkflowState{ID: n.State.Id, Name: n.State.Name, Type: n.State.Type, Color: n.State.Color, Position: n.State.Position},
			Team:     Team{ID: n.Team.Id, Name: n.Team.Name, Key: n.Team.Key},
			Assignee: assignee,
			Labels:   labels,
		}
		if n.Cycle != nil {
			issue.Cycle = &IssueCycle{ID: n.Cycle.Id, Name: derefStr(n.Cycle.Name), Number: int(n.Cycle.Number)}
		}
		if n.Project != nil {
			issue.Project = &IssueProject{ID: n.Project.Id, Name: n.Project.Name}
		}
		issues[i] = issue
	}
	return &IssueListResult{
		Issues: issues,
		PageInfo: PageInfo{
			HasNextPage: resp.Project.Issues.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.Project.Issues.PageInfo.EndCursor),
		},
	}, nil
}

func (c *GraphQLClient) ListInitiatives(ctx context.Context, status string, first int, after string) (*InitiativeListResult, error) {
	var filter *linear_graphql.InitiativeFilter
	if status != "" {
		filter = &linear_graphql.InitiativeFilter{
			Status: &linear_graphql.StringComparator{Eq: &status},
		}
	}
	resp, err := linear_graphql.ListInitiatives(ctx, c.gql, filter, &first, strPtrOrNil(after))
	if err != nil {
		return nil, fmt.Errorf("list initiatives: %w", err)
	}

	initiatives := make([]Initiative, len(resp.Initiatives.Nodes))
	for i, n := range resp.Initiatives.Nodes {
		initiatives[i] = initiativeFromListNode(n)
	}
	return &InitiativeListResult{
		Initiatives: initiatives,
		PageInfo: PageInfo{
			HasNextPage: resp.Initiatives.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.Initiatives.PageInfo.EndCursor),
		},
	}, nil
}

func initiativeFromListNode(n linear_graphql.ListInitiativesInitiativesInitiativeConnectionNodesInitiative) Initiative {
	var owner *User
	if n.Owner != nil {
		owner = &User{ID: n.Owner.Id, Name: n.Owner.Name, DisplayName: n.Owner.DisplayName, Email: n.Owner.Email, Active: n.Owner.Active}
	}
	return Initiative{
		ID: n.Id, Name: n.Name, Description: derefStr(n.Description),
		SlugID: n.SlugId, Color: derefStr(n.Color), Icon: derefStr(n.Icon),
		Status: string(n.Status), TargetDate: derefStr(n.TargetDate),
		Owner: owner, URL: n.Url,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt, ArchivedAt: n.ArchivedAt,
	}
}

func (c *GraphQLClient) GetInitiative(ctx context.Context, id string) (*Initiative, error) {
	resp, err := linear_graphql.GetInitiative(ctx, c.gql, id)
	if err != nil {
		return nil, fmt.Errorf("get initiative: %w", err)
	}
	n := resp.Initiative
	var owner *User
	if n.Owner != nil {
		owner = &User{ID: n.Owner.Id, Name: n.Owner.Name, DisplayName: n.Owner.DisplayName, Email: n.Owner.Email, Active: n.Owner.Active}
	}
	var health string
	if n.Health != nil {
		health = string(*n.Health)
	}
	return &Initiative{
		ID: n.Id, Name: n.Name, Description: derefStr(n.Description), Content: derefStr(n.Content),
		SlugID: n.SlugId, Color: derefStr(n.Color), Icon: derefStr(n.Icon),
		Status: string(n.Status), Health: health,
		TargetDate: derefStr(n.TargetDate),
		StartedAt:  n.StartedAt, CompletedAt: n.CompletedAt,
		Owner: owner, URL: n.Url,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt, ArchivedAt: n.ArchivedAt,
	}, nil
}

func (c *GraphQLClient) ListInitiativeProjects(ctx context.Context, initiativeID string, first int, after string) (*ProjectListResult, error) {
	resp, err := linear_graphql.ListInitiativeProjects(ctx, c.gql, initiativeID, nil, &first, strPtrOrNil(after))
	if err != nil {
		return nil, fmt.Errorf("list initiative projects: %w", err)
	}

	nodes := resp.Initiative.Projects.Nodes
	projects := make([]Project, len(nodes))
	for i, n := range nodes {
		var lead *User
		if n.Lead != nil {
			lead = &User{ID: n.Lead.Id, Name: n.Lead.Name, DisplayName: n.Lead.DisplayName, Email: n.Lead.Email, Active: n.Lead.Active}
		}
		projects[i] = Project{
			ID: n.Id, Name: n.Name, Description: n.Description,
			SlugID: n.SlugId, Color: n.Color, Icon: derefStr(n.Icon),
			Status:   ProjectStatus{Name: n.Status.Name, Color: n.Status.Color, Type: string(n.Status.Type)},
			Progress: n.Progress, Priority: n.Priority, PriorityLabel: n.PriorityLabel,
			StartDate: derefStr(n.StartDate), TargetDate: derefStr(n.TargetDate),
			Lead: lead, URL: n.Url,
			CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
			CompletedAt: n.CompletedAt, CanceledAt: n.CanceledAt, ArchivedAt: n.ArchivedAt,
		}
	}
	return &ProjectListResult{
		Projects: projects,
		PageInfo: PageInfo{
			HasNextPage: resp.Initiative.Projects.PageInfo.HasNextPage,
			EndCursor:   derefStr(resp.Initiative.Projects.PageInfo.EndCursor),
		},
	}, nil
}

func (c *GraphQLClient) FileUpload(ctx context.Context, contentType, filename string, size int64) (*UploadResult, error) {
	if size > math.MaxInt32 {
		return nil, fmt.Errorf("file size %d exceeds maximum of %d bytes", size, math.MaxInt32)
	}
	sizeInt := int(size)
	resp, err := linear_graphql.FileUpload(ctx, c.gql, contentType, filename, sizeInt)
	if err != nil {
		return nil, fmt.Errorf("file upload: %w", err)
	}
	if resp.FileUpload.UploadFile == nil {
		return nil, fmt.Errorf("file upload: no upload file returned")
	}
	uf := resp.FileUpload.UploadFile
	headers := make([]UploadHeader, len(uf.Headers))
	for i, h := range uf.Headers {
		headers[i] = UploadHeader{Key: h.Key, Value: h.Value}
	}
	return &UploadResult{
		UploadURL: uf.UploadUrl,
		AssetURL:  uf.AssetUrl,
		Headers:   headers,
	}, nil
}

func (c *GraphQLClient) UploadToURL(ctx context.Context, url string, headers []UploadHeader, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}
	for _, h := range headers {
		req.Header.Set(h.Key, h.Value)
	}
	// Use http.DefaultClient — presigned S3 URLs must not include Bearer auth.
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // presigned URL from Linear API
	if err != nil {
		return fmt.Errorf("upload to URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (c *GraphQLClient) DownloadURL(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create download request: %w", err)
	}
	// Use authenticated client for Linear upload URLs, default client for others.
	client := http.DefaultClient //nolint:gosec // external CDN URLs
	if strings.Contains(url, "uploads.linear.app") {
		client = c.httpClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download URL: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefIntFromFloat(p *float64) int {
	if p == nil {
		return 0
	}
	return int(*p)
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
