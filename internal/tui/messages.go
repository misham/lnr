package tui

import "github.com/misham/linear-cli/internal/api"

// IssuesLoadedMsg is sent when issues have been loaded from the API.
type IssuesLoadedMsg struct {
	Issues   []api.Issue
	PageInfo api.PageInfo
	Err      error
}

// IssueDetailLoadedMsg is sent when a single issue's detail has been loaded.
type IssueDetailLoadedMsg struct {
	Issue *api.Issue
	Err   error
}

// CyclesLoadedMsg is sent when cycles have been loaded from the API.
type CyclesLoadedMsg struct {
	Cycles   []api.Cycle
	PageInfo api.PageInfo
	Err      error
}

// CycleDetailLoadedMsg is sent when a single cycle's detail has been loaded.
type CycleDetailLoadedMsg struct {
	Cycle *api.Cycle
	Err   error
}

// CycleIssuesLoadedMsg is sent when a cycle's issues have been loaded.
type CycleIssuesLoadedMsg struct {
	Issues []api.Issue
	Err    error
}

// ProjectsLoadedMsg is sent when projects have been loaded from the API.
type ProjectsLoadedMsg struct {
	Projects []api.Project
	PageInfo api.PageInfo
	Err      error
}

// ProjectDetailLoadedMsg is sent when a single project's detail has been loaded.
type ProjectDetailLoadedMsg struct {
	Project *api.Project
	Err     error
}

// ProjectIssuesLoadedMsg is sent when a project's issues have been loaded.
type ProjectIssuesLoadedMsg struct {
	Issues []api.Issue
	Err    error
}

// InitiativesLoadedMsg is sent when initiatives have been loaded from the API.
type InitiativesLoadedMsg struct {
	Initiatives []api.Initiative
	PageInfo    api.PageInfo
	Err         error
}

// InitiativeDetailLoadedMsg is sent when a single initiative's detail has been loaded.
type InitiativeDetailLoadedMsg struct {
	Initiative *api.Initiative
	Err        error
}

// InitiativeProjectsLoadedMsg is sent when an initiative's linked projects have been loaded.
type InitiativeProjectsLoadedMsg struct {
	Projects []api.Project
	PageInfo api.PageInfo
	Err      error
}

// ImageFetchedMsg is sent when an image has been downloaded for inline preview.
type ImageFetchedMsg struct {
	URL  string
	Data []byte
	Err  error
}

// FileOpenedMsg is sent after attempting to open a file with the system default app.
type FileOpenedMsg struct {
	Err error
}

// FileUploadedMsg is sent after a file upload + comment creation completes.
type FileUploadedMsg struct {
	Err error
}
