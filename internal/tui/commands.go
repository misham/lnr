package tui

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/ui"
)

func loadIssues(ctx context.Context, client api.Client, teamID string, first int, after string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return IssuesLoadedMsg{Err: ctx.Err()}
		}
		result, err := client.ListIssues(ctx, teamID, first, after)
		if err != nil {
			return IssuesLoadedMsg{Err: err}
		}
		return IssuesLoadedMsg{Issues: result.Issues, PageInfo: result.PageInfo}
	}
}

func loadIssueDetail(ctx context.Context, client api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return IssueDetailLoadedMsg{Err: ctx.Err()}
		}
		issue, err := client.GetIssue(ctx, id)
		if err != nil {
			return IssueDetailLoadedMsg{Err: err}
		}
		return IssueDetailLoadedMsg{Issue: issue}
	}
}

func loadCycles(ctx context.Context, client api.Client, teamID string, first int, after string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return CyclesLoadedMsg{Err: ctx.Err()}
		}
		result, err := client.ListCycles(ctx, teamID, true, first, after)
		if err != nil {
			return CyclesLoadedMsg{Err: err}
		}
		return CyclesLoadedMsg{Cycles: result.Cycles, PageInfo: result.PageInfo}
	}
}

func loadCycleIssues(ctx context.Context, client api.Client, cycleID string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return CycleIssuesLoadedMsg{Err: ctx.Err()}
		}
		result, err := client.ListCycleIssues(ctx, cycleID, 50, "")
		if err != nil {
			return CycleIssuesLoadedMsg{Err: err}
		}
		return CycleIssuesLoadedMsg{Issues: result.Issues}
	}
}

func loadCycleDetail(ctx context.Context, client api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return CycleDetailLoadedMsg{Err: ctx.Err()}
		}
		cycle, err := client.GetCycle(ctx, id)
		if err != nil {
			return CycleDetailLoadedMsg{Err: err}
		}
		return CycleDetailLoadedMsg{Cycle: cycle}
	}
}

func loadProjects(ctx context.Context, client api.Client, teamID string, first int, after string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return ProjectsLoadedMsg{Err: ctx.Err()}
		}
		result, err := client.ListProjects(ctx, teamID, "", first, after)
		if err != nil {
			return ProjectsLoadedMsg{Err: err}
		}
		return ProjectsLoadedMsg{Projects: result.Projects, PageInfo: result.PageInfo}
	}
}

func loadProjectIssues(ctx context.Context, client api.Client, projectID string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return ProjectIssuesLoadedMsg{Err: ctx.Err()}
		}
		result, err := client.ListProjectIssues(ctx, projectID, 50, "")
		if err != nil {
			return ProjectIssuesLoadedMsg{Err: err}
		}
		return ProjectIssuesLoadedMsg{Issues: result.Issues}
	}
}

func loadProjectDetail(ctx context.Context, client api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return ProjectDetailLoadedMsg{Err: ctx.Err()}
		}
		project, err := client.GetProject(ctx, id)
		if err != nil {
			return ProjectDetailLoadedMsg{Err: err}
		}
		return ProjectDetailLoadedMsg{Project: project}
	}
}

func loadInitiatives(ctx context.Context, client api.Client, first int, after string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return InitiativesLoadedMsg{Err: ctx.Err()}
		}
		result, err := client.ListInitiatives(ctx, "", first, after)
		if err != nil {
			return InitiativesLoadedMsg{Err: err}
		}
		return InitiativesLoadedMsg{Initiatives: result.Initiatives, PageInfo: result.PageInfo}
	}
}

func loadInitiativeDetail(ctx context.Context, client api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return InitiativeDetailLoadedMsg{Err: ctx.Err()}
		}
		initiative, err := client.GetInitiative(ctx, id)
		if err != nil {
			return InitiativeDetailLoadedMsg{Err: err}
		}
		return InitiativeDetailLoadedMsg{Initiative: initiative}
	}
}

func loadInitiativeProjects(ctx context.Context, client api.Client, id string, first int, after string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return InitiativeProjectsLoadedMsg{Err: ctx.Err()}
		}
		result, err := client.ListInitiativeProjects(ctx, id, first, after)
		if err != nil {
			return InitiativeProjectsLoadedMsg{Err: err}
		}
		return InitiativeProjectsLoadedMsg{Projects: result.Projects, PageInfo: result.PageInfo}
	}
}

func fetchImage(ctx context.Context, client api.Client, url string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return ImageFetchedMsg{URL: url, Err: ctx.Err()}
		}
		body, err := client.DownloadURL(ctx, url)
		if err != nil {
			return ImageFetchedMsg{URL: url, Err: err}
		}
		defer func() { _ = body.Close() }()
		data, err := io.ReadAll(body)
		if err != nil {
			return ImageFetchedMsg{URL: url, Err: err}
		}
		return ImageFetchedMsg{URL: url, Data: data}
	}
}

func openFile(ctx context.Context, client api.Client, url, filename string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return FileOpenedMsg{Err: ctx.Err()}
		}
		body, err := client.DownloadURL(ctx, url)
		if err != nil {
			return FileOpenedMsg{Err: err}
		}
		defer func() { _ = body.Close() }()

		f, err := os.CreateTemp("", "lnr-*-"+filepath.Base(filename))
		if err != nil {
			return FileOpenedMsg{Err: err}
		}
		if _, err := io.Copy(f, body); err != nil {
			_ = f.Close()
			return FileOpenedMsg{Err: err}
		}
		tmp := f.Name()
		_ = f.Close()

		opener := openCommand()
		if err := exec.Command(opener, tmp).Start(); err != nil { //nolint:gosec // user-initiated open of downloaded file
			return FileOpenedMsg{Err: err}
		}
		return FileOpenedMsg{}
	}
}

func uploadFile(ctx context.Context, client api.Client, issue *api.Issue, filePath string) tea.Cmd {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return FileUploadedMsg{Err: ctx.Err()}
		}

		filePath = filepath.Clean(filePath)
		file, err := os.Open(filePath)
		if err != nil {
			return FileUploadedMsg{Err: err}
		}
		defer func() { _ = file.Close() }()

		info, err := file.Stat()
		if err != nil {
			return FileUploadedMsg{Err: err}
		}

		filename := filepath.Base(filePath)
		contentType := "application/octet-stream"
		if ext := filepath.Ext(filename); ext != "" {
			if ct := mime.TypeByExtension(ext); ct != "" {
				contentType = ct
			}
		}

		result, err := client.FileUpload(ctx, contentType, filename, info.Size())
		if err != nil {
			return FileUploadedMsg{Err: err}
		}

		if _, err := file.Seek(0, 0); err != nil {
			return FileUploadedMsg{Err: err}
		}
		headers := append(result.Headers, api.UploadHeader{Key: "Content-Type", Value: contentType})
		if err := client.UploadToURL(ctx, result.UploadURL, headers, file); err != nil {
			return FileUploadedMsg{Err: err}
		}

		prefix := ""
		if ui.IsImageFile(filename) {
			prefix = "!"
		}
		link := fmt.Sprintf("%s[%s](%s)", prefix, filename, result.AssetURL)
		desc := issue.Description
		if desc != "" {
			desc += "\n\n"
		}
		desc += link
		if _, err := client.UpdateIssue(ctx, issue.ID, api.IssueUpdateInput{
			Description: &desc,
		}); err != nil {
			return FileUploadedMsg{Err: err}
		}

		return FileUploadedMsg{}
	}
}

func openCommand() string {
	switch runtime.GOOS {
	case "darwin":
		return "open"
	case "windows":
		return "cmd"
	default:
		return "xdg-open"
	}
}
