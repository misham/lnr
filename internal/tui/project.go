package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

type projectTab struct {
	client       api.Client
	ctx          context.Context
	teamID       string
	state        tabState
	items        []api.Project
	detail       *api.Project
	detailIssues []api.Issue
	cursor       int
	offset       int
	pageInfo     api.PageInfo
	size         tea.WindowSizeMsg
	spinner      spinner.Model
	viewport     viewport.Model
	errMsg       string
	keys         KeyMap
	loadingMore  bool
}

// NewProjectTab creates a new project tab.
func NewProjectTab(ctx context.Context, client api.Client, teamID string) *projectTab {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &projectTab{
		client:  client,
		ctx:     ctx,
		teamID:  teamID,
		state:   stateLoading,
		spinner: s,
		keys:    DefaultKeyMap(),
	}
}

func (t *projectTab) Init() tea.Cmd {
	return tea.Batch(
		t.spinner.Tick,
		loadProjects(t.ctx, t.client, t.teamID, defaultPageSize, ""),
	)
}

func (t *projectTab) Update(msg tea.Msg) (TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.size = msg
		t.viewport.Width = msg.Width
		t.viewport.Height = msg.Height
		return t, nil

	case ProjectsLoadedMsg:
		if msg.Err != nil {
			t.state = stateError
			t.errMsg = msg.Err.Error()
			t.loadingMore = false
			return t, nil
		}
		prevLen := len(t.items)
		t.items = append(t.items, msg.Projects...)
		t.pageInfo = msg.PageInfo
		t.state = stateList
		if t.loadingMore && prevLen > 0 {
			t.cursor = prevLen
			t.adjustOffset()
		}
		t.loadingMore = false
		return t, nil

	case ProjectDetailLoadedMsg:
		if msg.Err != nil {
			t.state = stateError
			t.errMsg = msg.Err.Error()
			return t, nil
		}
		t.detail = msg.Project
		t.detailIssues = nil
		t.state = stateDetail
		t.viewport.SetContent(t.renderDetail())
		t.viewport.GotoTop()
		return t, loadProjectIssues(t.ctx, t.client, msg.Project.ID)

	case ProjectIssuesLoadedMsg:
		if msg.Err != nil {
			return t, nil
		}
		t.detailIssues = msg.Issues
		if t.state == stateDetail {
			t.viewport.SetContent(t.renderDetail())
		}
		return t, nil

	case spinner.TickMsg:
		if t.state == stateLoading {
			var cmd tea.Cmd
			t.spinner, cmd = t.spinner.Update(msg)
			return t, cmd
		}
		return t, nil

	case tea.KeyMsg:
		return t.handleKey(msg)
	}

	return t, nil
}

func (t *projectTab) handleKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch t.state {
	case stateList:
		return t.handleListKey(msg)
	case stateDetail:
		return t.handleDetailKey(msg)
	case stateError:
		if key.Matches(msg, t.keys.Retry) {
			t.state = stateLoading
			t.items = nil
			t.cursor = 0
			t.offset = 0
			return t, tea.Batch(
				t.spinner.Tick,
				loadProjects(t.ctx, t.client, t.teamID, defaultPageSize, ""),
			)
		}
		return t, nil
	}
	return t, nil
}

func (t *projectTab) handleListKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch {
	case key.Matches(msg, t.keys.Down):
		if t.cursor < len(t.items)-1 {
			t.cursor++
			t.adjustOffset()
		}
	case key.Matches(msg, t.keys.Up):
		if t.cursor > 0 {
			t.cursor--
			t.adjustOffset()
		}
	case key.Matches(msg, t.keys.Enter):
		if len(t.items) > 0 {
			return t, loadProjectDetail(t.ctx, t.client, t.items[t.cursor].ID)
		}
	case key.Matches(msg, t.keys.NextPage):
		if t.pageInfo.HasNextPage {
			t.loadingMore = true
			return t, loadProjects(t.ctx, t.client, t.teamID, defaultPageSize, t.pageInfo.EndCursor)
		}
	case key.Matches(msg, t.keys.PrevPage):
		visible := t.visibleRows()
		t.cursor = max(t.cursor-visible, 0)
		t.adjustOffset()
	case key.Matches(msg, t.keys.HalfDown):
		half := t.visibleRows() / 2
		t.cursor = min(t.cursor+half, len(t.items)-1)
		t.adjustOffset()
	case key.Matches(msg, t.keys.HalfUp):
		half := t.visibleRows() / 2
		t.cursor = max(t.cursor-half, 0)
		t.adjustOffset()
	case key.Matches(msg, t.keys.Top):
		t.cursor = 0
		t.offset = 0
	case key.Matches(msg, t.keys.Bottom):
		if len(t.items) > 0 {
			t.cursor = len(t.items) - 1
			t.adjustOffset()
		}
	}
	return t, nil
}

func (t *projectTab) handleDetailKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch {
	case key.Matches(msg, t.keys.Back):
		t.state = stateList
		t.detail = nil
		t.detailIssues = nil
		return t, nil
	case key.Matches(msg, t.keys.Down):
		t.viewport.ScrollDown(1)
	case key.Matches(msg, t.keys.Up):
		t.viewport.ScrollUp(1)
	}
	return t, nil
}

func (t *projectTab) visibleRows() int {
	// -1 for header, -1 for separator, -1 for pagination hint.
	return max(t.size.Height-3, 1)
}

func (t *projectTab) adjustOffset() {
	visible := t.visibleRows()
	if t.cursor < t.offset {
		t.offset = t.cursor
	}
	if t.cursor >= t.offset+visible {
		t.offset = t.cursor - visible + 1
	}
}

func (t *projectTab) View() string {
	switch t.state {
	case stateLoading:
		return t.spinner.View() + " Loading projects..."
	case stateError:
		return errorStyle.Render("Error: "+t.errMsg) + "\n\nPress r to retry"
	case stateDetail:
		if t.viewport.Width == 0 {
			return t.renderDetail()
		}
		return t.viewport.View()
	case stateList:
		return t.renderList()
	}
	return ""
}

func (t *projectTab) renderList() string {
	if len(t.items) == 0 {
		return "No projects found"
	}

	var b strings.Builder

	header := fmt.Sprintf("   %-25s %-12s %-20s %s", "Name", "Status", "Lead", "Progress")
	b.WriteString(listHeaderStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(listSepStyle.Render(strings.Repeat("─", max(t.size.Width-4, 40))))
	b.WriteString("\n")

	visible := t.visibleRows()
	end := min(t.offset+visible, len(t.items))

	for i := t.offset; i < end; i++ {
		project := t.items[i]

		lead := ""
		if project.Lead != nil {
			lead = project.Lead.DisplayName
		}
		progress := fmt.Sprintf("%.0f%%", project.Progress*100)

		var statusDot string
		switch project.Status.Type {
		case "started":
			statusDot = lipgloss.NewStyle().Foreground(greenColor).Render("●")
		case "planned":
			statusDot = lipgloss.NewStyle().Foreground(yellowColor).Render("●")
		case "backlog":
			statusDot = lipgloss.NewStyle().Foreground(cyanColor).Render("●")
		default:
			statusDot = lipgloss.NewStyle().Foreground(mutedColor).Render("●")
		}

		stateCol := fmt.Sprintf("%s %-10s", statusDot, truncate(project.Status.Name, 9))
		leadCol := fmt.Sprintf("%-20s", truncate(lead, 18))

		if i == t.cursor {
			name := fmt.Sprintf("%-25s", truncate(project.Name, 23))
			line := fmt.Sprintf(" %s %s %s %s %s",
				accentBarStyle.Render("▎"),
				selectedTitleStyle.Render(name),
				stateCol,
				leadCol,
				progress,
			)
			b.WriteString(selectedRowStyle.Width(t.size.Width).Render(line))
		} else {
			line := fmt.Sprintf("   %-25s %s %s %s",
				truncate(project.Name, 23),
				stateCol,
				leadCol,
				progress,
			)
			b.WriteString(listRowStyle.Render(line))
		}
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	if t.loadingMore {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("   Loading more..."))
	} else if t.pageInfo.HasNextPage {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(fmt.Sprintf("   %d loaded  •  n next page", len(t.items))))
	}

	return b.String()
}

func (t *projectTab) renderDetail() string {
	if t.detail == nil {
		return ""
	}
	d := t.detail
	var b strings.Builder

	b.WriteString(detailTitleStyle.Render(d.Name))
	b.WriteString("\n\n")

	if d.Description != "" {
		b.WriteString(d.Description)
		b.WriteString("\n\n")
	}

	fmt.Fprintf(&b, "  Status:   %s\n", d.Status.Name)
	fmt.Fprintf(&b, "  Progress: %.0f%%\n", d.Progress*100)
	if d.Lead != nil {
		fmt.Fprintf(&b, "  Lead:     %s\n", d.Lead.DisplayName)
	}
	if d.StartDate != "" {
		fmt.Fprintf(&b, "  Start:    %s\n", d.StartDate)
	}
	if d.TargetDate != "" {
		fmt.Fprintf(&b, "  Target:   %s\n", d.TargetDate)
	}
	if d.URL != "" {
		fmt.Fprintf(&b, "  URL:      %s\n", d.URL)
	}

	if len(d.Milestones) > 0 {
		b.WriteString("\n")
		fmt.Fprintf(&b, "Milestones (%d)\n", len(d.Milestones))
		for _, m := range d.Milestones {
			target := ""
			if m.TargetDate != "" {
				target = fmt.Sprintf(" (due %s)", m.TargetDate)
			}
			fmt.Fprintf(&b, "  - %s%s\n", m.Name, target)
		}
	}

	if len(t.detailIssues) > 0 {
		b.WriteString("\n")
		fmt.Fprintf(&b, "Issues (%d)\n", len(t.detailIssues))
		for _, issue := range t.detailIssues {
			dot := stateIcon(issue.State.Type)
			fmt.Fprintf(&b, "  %s %s  %s  %s\n", dot, issue.Identifier, truncate(issue.State.Name, 12), issue.Title)
		}
	}

	return b.String()
}

func (t *projectTab) ShortHelp() []key.Binding {
	switch t.state {
	case stateList:
		return []key.Binding{t.keys.Up, t.keys.Down, t.keys.Enter, t.keys.NextPage}
	case stateDetail:
		return []key.Binding{t.keys.Up, t.keys.Down, t.keys.Back}
	case stateError:
		return []key.Binding{t.keys.Retry}
	}
	return nil
}
