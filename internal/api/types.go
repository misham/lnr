package api

import "time"

// Team represents a Linear team.
type Team struct {
	ID          string
	Name        string
	Key         string
	Description string
	Private     bool
	Icon        string
	Color       string
	Timezone    string
}

// User represents a Linear user.
type User struct {
	ID          string
	Name        string
	DisplayName string
	Email       string
	Active      bool
}

// Issue represents a Linear issue.
type Issue struct {
	ID            string
	Identifier    string
	Title         string
	Description   string
	Priority      int
	PriorityLabel string
	Estimate      int
	DueDate       string
	State         WorkflowState
	Team          Team
	Assignee      *User
	Labels        []IssueLabel
	Comments      []Comment
	Cycle         *IssueCycle
	Project       *IssueProject
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time
	ArchivedAt    *time.Time
	URL           string
}

// IssueCycle is the cycle an issue belongs to (lightweight reference).
type IssueCycle struct {
	ID     string
	Name   string
	Number int
}

// IssueProject is the project an issue belongs to (lightweight reference).
type IssueProject struct {
	ID   string
	Name string
}

// WorkflowState represents a workflow state in Linear.
type WorkflowState struct {
	ID       string
	Name     string
	Type     string
	Color    string
	Position float64
}

// Comment represents a comment on a Linear issue.
type Comment struct {
	ID        string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
}

// IssueLabel represents a label on a Linear issue.
type IssueLabel struct {
	ID    string
	Name  string
	Color string
}

// IssueCreateInput holds parameters for creating an issue.
type IssueCreateInput struct {
	Title       string
	TeamID      string
	Description string
	Priority    int
	StateID     string
	AssigneeID  string
	LabelIDs    []string
	Estimate    int
	DueDate     string
}

// IssueUpdateInput holds parameters for updating an issue.
type IssueUpdateInput struct {
	Title       *string
	Description *string
	Priority    *int
	StateID     *string
	AssigneeID  *string
	LabelIDs    []string
	Estimate    *int
	DueDate     *string
	CycleID     *string
}

// Cycle represents a Linear cycle (sprint).
type Cycle struct {
	ID          string
	Name        string
	Number      int
	Description string
	StartsAt    time.Time
	EndsAt      time.Time
	CompletedAt *time.Time
	Progress    float64
	IsActive    bool
	IsNext      bool
	IsPast      bool
	Team        Team
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CycleListResult holds a page of cycles with pagination info.
type CycleListResult struct {
	Cycles   []Cycle
	PageInfo PageInfo
}

// IssueListResult holds a page of issues with pagination info.
type IssueListResult struct {
	Issues   []Issue
	PageInfo PageInfo
}

// ProjectStatus represents the status of a Linear project.
type ProjectStatus struct {
	Name  string
	Color string
	Type  string // ProjectStatusType: backlog, canceled, completed, paused, planned, started
}

// ProjectMilestone represents a milestone within a project.
type ProjectMilestone struct {
	ID          string
	Name        string
	Description string
	TargetDate  string
	SortOrder   float64
}

// Project represents a Linear project.
type Project struct {
	ID            string
	Name          string
	Description   string
	Content       string
	SlugID        string
	Color         string
	Icon          string
	Status        ProjectStatus
	Progress      float64
	Priority      int
	PriorityLabel string
	StartDate     string
	TargetDate    string
	Lead          *User
	URL           string
	Milestones    []ProjectMilestone
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time
	CanceledAt    *time.Time
	ArchivedAt    *time.Time
}

// ProjectListResult holds a page of projects with pagination info.
type ProjectListResult struct {
	Projects []Project
	PageInfo PageInfo
}

// Initiative represents a Linear initiative.
type Initiative struct {
	ID          string
	Name        string
	Description string
	Content     string
	SlugID      string
	Color       string
	Icon        string
	Status      string // InitiativeStatus: Planned, Active, Completed
	Health      string
	TargetDate  string
	StartedAt   *time.Time
	CompletedAt *time.Time
	Owner       *User
	URL         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ArchivedAt  *time.Time
}

// InitiativeListResult holds a page of initiatives with pagination info.
type InitiativeListResult struct {
	Initiatives []Initiative
	PageInfo    PageInfo
}

// UploadResult holds the presigned URL and asset URL from fileUpload.
type UploadResult struct {
	UploadURL string
	AssetURL  string
	Headers   []UploadHeader
}

// UploadHeader is a key-value pair for S3 presigned URL headers.
type UploadHeader struct {
	Key   string
	Value string
}
