package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/ui"
)

type tabState int

const (
	stateLoading tabState = iota
	stateList
	stateDetail
	stateError
)

const defaultPageSize = 50

type issueTab struct {
	client        api.Client
	ctx           context.Context
	teamID        string
	state         tabState
	items         []api.Issue
	detail        *api.Issue
	cursor        int
	offset        int
	pageInfo      api.PageInfo
	size          tea.WindowSizeMsg
	spinner       spinner.Model
	viewport      viewport.Model
	errMsg        string
	keys          KeyMap
	filter        string
	filtered      []api.Issue
	searching     bool
	loadingMore   bool
	uploading     bool
	uploadInput   string
	fetchedImages map[string][]byte
}

// NewIssueTab creates a new issue tab.
func NewIssueTab(ctx context.Context, client api.Client, teamID string) *issueTab {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &issueTab{
		client:  client,
		ctx:     ctx,
		teamID:  teamID,
		state:   stateLoading,
		spinner: s,
		keys:    DefaultKeyMap(),
	}
}

func (t *issueTab) Init() tea.Cmd {
	return tea.Batch(
		t.spinner.Tick,
		loadIssues(t.ctx, t.client, t.teamID, defaultPageSize, ""),
	)
}

func (t *issueTab) Update(msg tea.Msg) (TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.size = msg
		t.viewport.Width = msg.Width
		t.viewport.Height = msg.Height
		if t.state == stateDetail {
			t.viewport.SetContent(t.renderDetail())
		}
		return t, nil

	case IssuesLoadedMsg:
		if msg.Err != nil {
			t.state = stateError
			t.errMsg = msg.Err.Error()
			t.loadingMore = false
			return t, nil
		}
		prevLen := len(t.items)
		t.items = append(t.items, msg.Issues...)
		t.pageInfo = msg.PageInfo
		t.state = stateList
		if t.loadingMore && prevLen > 0 {
			t.cursor = prevLen
			t.adjustOffset()
		}
		t.loadingMore = false
		return t, nil

	case IssueDetailLoadedMsg:
		if msg.Err != nil {
			t.state = stateError
			t.errMsg = msg.Err.Error()
			return t, nil
		}
		t.detail = msg.Issue
		t.fetchedImages = nil
		t.state = stateDetail
		t.viewport.SetContent(t.renderDetail())
		t.viewport.GotoTop()

		// If on iTerm2, fetch images for inline preview.
		var cmds []tea.Cmd
		if ui.IsITerm2() {
			files := t.detailFiles()
			for _, f := range files {
				if f.IsImage {
					cmds = append(cmds, fetchImage(t.ctx, t.client, f.URL))
				}
			}
		}
		if len(cmds) > 0 {
			return t, tea.Batch(cmds...)
		}
		return t, nil

	case ImageFetchedMsg:
		if msg.Err == nil && len(msg.Data) > 0 {
			if t.fetchedImages == nil {
				t.fetchedImages = make(map[string][]byte)
			}
			t.fetchedImages[msg.URL] = msg.Data
			if t.state == stateDetail {
				t.viewport.SetContent(t.renderDetail())
			}
		}
		return t, nil

	case FileOpenedMsg:
		return t, nil

	case FileUploadedMsg:
		t.uploading = false
		t.uploadInput = ""
		if msg.Err != nil {
			t.errMsg = msg.Err.Error()
			return t, nil
		}
		if t.detail != nil {
			return t, loadIssueDetail(t.ctx, t.client, t.detail.ID)
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

func (t *issueTab) handleKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch t.state {
	case stateList:
		if t.searching {
			return t.handleSearchKey(msg)
		}
		return t.handleListKey(msg)
	case stateDetail:
		if t.uploading {
			return t.handleUploadKey(msg)
		}
		return t.handleDetailKey(msg)
	case stateError:
		if key.Matches(msg, t.keys.Retry) {
			t.state = stateLoading
			t.items = nil
			t.filter = ""
			t.filtered = nil
			t.cursor = 0
			t.offset = 0
			return t, tea.Batch(
				t.spinner.Tick,
				loadIssues(t.ctx, t.client, t.teamID, defaultPageSize, ""),
			)
		}
		return t, nil
	}
	return t, nil
}

func (t *issueTab) handleSearchKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		t.searching = false
		t.filter = ""
		t.filtered = nil
		t.cursor = 0
		t.offset = 0
		return t, nil
	case tea.KeyEnter:
		t.searching = false
		return t, nil
	case tea.KeyBackspace:
		if len(t.filter) > 0 {
			t.filter = t.filter[:len(t.filter)-1]
			t.applyFilter()
		}
		return t, nil
	case tea.KeyRunes:
		t.filter += string(msg.Runes)
		t.applyFilter()
		return t, nil
	}
	return t, nil
}

func (t *issueTab) applyFilter() {
	if t.filter == "" {
		t.filtered = nil
		t.cursor = 0
		t.offset = 0
		return
	}
	query := strings.ToLower(t.filter)
	t.filtered = nil
	for _, issue := range t.items {
		if strings.Contains(strings.ToLower(issue.Title), query) ||
			strings.Contains(strings.ToLower(issue.Identifier), query) ||
			strings.Contains(strings.ToLower(issue.State.Name), query) ||
			strings.Contains(strings.ToLower(issue.PriorityLabel), query) {
			t.filtered = append(t.filtered, issue)
		}
	}
	t.cursor = 0
	t.offset = 0
}

func (t *issueTab) displayItems() []api.Issue {
	if t.filter != "" {
		return t.filtered
	}
	return t.items
}

func (t *issueTab) handleListKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	items := t.displayItems()
	switch {
	case key.Matches(msg, t.keys.Search):
		t.searching = true
		return t, nil
	case key.Matches(msg, t.keys.Down):
		if t.cursor < len(items)-1 {
			t.cursor++
			t.adjustOffset()
		}
	case key.Matches(msg, t.keys.Up):
		if t.cursor > 0 {
			t.cursor--
			t.adjustOffset()
		}
	case key.Matches(msg, t.keys.Enter):
		if len(items) > 0 {
			return t, loadIssueDetail(t.ctx, t.client, items[t.cursor].ID)
		}
	case key.Matches(msg, t.keys.NextPage):
		if t.filter == "" && t.pageInfo.HasNextPage {
			t.loadingMore = true
			return t, loadIssues(t.ctx, t.client, t.teamID, defaultPageSize, t.pageInfo.EndCursor)
		}
	case key.Matches(msg, t.keys.PrevPage):
		visible := t.visibleRows()
		t.cursor = max(t.cursor-visible, 0)
		t.adjustOffset()
	case key.Matches(msg, t.keys.HalfDown):
		half := t.visibleRows() / 2
		t.cursor = min(t.cursor+half, len(items)-1)
		t.adjustOffset()
	case key.Matches(msg, t.keys.HalfUp):
		half := t.visibleRows() / 2
		t.cursor = max(t.cursor-half, 0)
		t.adjustOffset()
	case key.Matches(msg, t.keys.Top):
		t.cursor = 0
		t.offset = 0
	case key.Matches(msg, t.keys.Bottom):
		if len(items) > 0 {
			t.cursor = len(items) - 1
			t.adjustOffset()
		}
	}
	return t, nil
}

func (t *issueTab) handleDetailKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch {
	case key.Matches(msg, t.keys.Back):
		t.state = stateList
		t.detail = nil
		return t, nil
	case key.Matches(msg, t.keys.Down):
		t.viewport.ScrollDown(1)
	case key.Matches(msg, t.keys.Up):
		t.viewport.ScrollUp(1)
	case key.Matches(msg, t.keys.Open):
		files := t.detailFiles()
		if len(files) > 0 {
			return t, openFile(t.ctx, t.client, files[0].URL, files[0].Name)
		}
	case key.Matches(msg, t.keys.Upload):
		t.uploading = true
		t.uploadInput = ""
		return t, nil
	}
	return t, nil
}

func (t *issueTab) handleUploadKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		t.uploading = false
		t.uploadInput = ""
		return t, nil
	case tea.KeyEnter:
		if t.uploadInput != "" && t.detail != nil {
			path := t.uploadInput
			t.uploading = false
			t.uploadInput = ""
			return t, uploadFile(t.ctx, t.client, t.detail, path)
		}
		return t, nil
	case tea.KeyBackspace:
		if len(t.uploadInput) > 0 {
			t.uploadInput = t.uploadInput[:len(t.uploadInput)-1]
		}
		return t, nil
	case tea.KeyRunes:
		t.uploadInput += string(msg.Runes)
		return t, nil
	}
	return t, nil
}

func (t *issueTab) visibleRows() int {
	// -1 for header, -1 for separator, -1 for pagination/search hint.
	return max(t.size.Height-3, 1)
}

func (t *issueTab) adjustOffset() {
	visible := t.visibleRows()
	if t.cursor < t.offset {
		t.offset = t.cursor
	}
	if t.cursor >= t.offset+visible {
		t.offset = t.cursor - visible + 1
	}
}

func (t *issueTab) View() string {
	switch t.state {
	case stateLoading:
		return t.spinner.View() + " Loading issues..."
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

func (t *issueTab) renderList() string {
	items := t.displayItems()
	if len(items) == 0 && t.filter == "" {
		return "No issues found"
	}

	var b strings.Builder

	// Search bar
	if t.searching || t.filter != "" {
		b.WriteString("/" + t.filter)
		if t.searching {
			b.WriteString("_")
		}
		b.WriteString("\n")
	}

	if len(items) == 0 {
		b.WriteString("No matching issues")
		return b.String()
	}

	header := fmt.Sprintf("   %-12s %-15s %-10s %s", "ID", "State", "Priority", "Title")
	b.WriteString(listHeaderStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(listSepStyle.Render(strings.Repeat("─", max(t.size.Width-4, 40))))
	b.WriteString("\n")

	visible := t.visibleRows()
	end := min(t.offset+visible, len(items))

	for i := t.offset; i < end; i++ {
		issue := items[i]
		stateDot := stateIcon(issue.State.Type)
		priStyle := priorityStyle(issue.PriorityLabel)

		id := fmt.Sprintf("%-12s", issue.Identifier)
		state := fmt.Sprintf("%s %-13s", stateDot, truncate(issue.State.Name, 12))
		pri := priStyle.Render(fmt.Sprintf("%-10s", truncate(issue.PriorityLabel, 9)))
		title := issue.Title
		badges := issueBadges(issue)

		if i == t.cursor {
			line := fmt.Sprintf(" %s %s %s %s %s%s",
				accentBarStyle.Render("▎"),
				selectedIDStyle.Render(id),
				state,
				pri,
				selectedTitleStyle.Render(title),
				badges,
			)
			b.WriteString(selectedRowStyle.Width(t.size.Width).Render(line))
		} else {
			line := fmt.Sprintf("   %s %s %s %s%s", id, state, pri, title, badges)
			b.WriteString(listRowStyle.Render(line))
		}
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	if t.loadingMore {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("   Loading more..."))
	} else if t.filter == "" && t.pageInfo.HasNextPage {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(fmt.Sprintf("   %d loaded  •  n next page", len(t.items))))
	}

	return b.String()
}

func renderMarkdown(content string, width int) string {
	if strings.TrimSpace(content) == "" {
		return ""
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return content
	}
	out, err := r.Render(content)
	if err != nil {
		return content
	}
	return out
}

func (t *issueTab) renderDetail() string {
	if t.detail == nil {
		return ""
	}
	d := t.detail
	var b strings.Builder

	width := t.viewport.Width
	if width == 0 {
		width = 80
	}

	b.WriteString(detailTitleStyle.Render(d.Identifier + " " + d.Title))
	b.WriteString("\n\n")

	fmt.Fprintf(&b, "  State:    %s\n", d.State.Name)
	fmt.Fprintf(&b, "  Priority: %s\n", d.PriorityLabel)
	if d.Assignee != nil {
		fmt.Fprintf(&b, "  Assignee: %s\n", d.Assignee.DisplayName)
	}
	if d.DueDate != "" {
		fmt.Fprintf(&b, "  Due:      %s\n", d.DueDate)
	}
	fmt.Fprintf(&b, "  Created:  %s\n", d.CreatedAt.Format("2006-01-02"))

	if len(d.Labels) > 0 {
		var labelNames []string
		for _, l := range d.Labels {
			labelNames = append(labelNames, l.Name)
		}
		fmt.Fprintf(&b, "  Labels:   %s\n", detailLabelStyle.Render(strings.Join(labelNames, ", ")))
	}

	if d.Description != "" {
		b.WriteString("\n")
		b.WriteString(renderMarkdown(d.Description, width))
	}

	if len(d.Comments) > 0 {
		b.WriteString("\n")
		fmt.Fprintf(&b, "Comments (%d)\n", len(d.Comments))
		for _, c := range d.Comments {
			author := "Unknown"
			if c.User != nil {
				author = c.User.DisplayName
			}
			fmt.Fprintf(&b, "\n  %s — %s\n", author, c.CreatedAt.Format("2006-01-02 15:04"))
			b.WriteString(renderMarkdown(c.Body, width))
		}
	}

	files := t.detailFiles()

	// iTerm2 inline image previews.
	if ui.IsITerm2() && t.fetchedImages != nil {
		for imgURL, data := range t.fetchedImages {
			filename := "image"
			for _, f := range files {
				if f.URL == imgURL {
					filename = f.Name
					break
				}
			}
			b.WriteString("\n")
			b.WriteString(ui.RenderInlineImage(data, filename))
			b.WriteString("\n")
		}
	}

	// Files section.
	if len(files) > 0 {
		b.WriteString("\n")
		fmt.Fprintf(&b, "%s (%d)\n", detailTitleStyle.Render("Files"), len(files))
		for _, f := range files {
			fmt.Fprintf(&b, "  %s\n", ui.HyperlinkOSC8(f.URL, f.Name))
		}
	}

	// Upload prompt.
	if t.uploading {
		b.WriteString("\n")
		b.WriteString("File path: " + t.uploadInput + "_\n")
	}

	return b.String()
}

func (t *issueTab) detailFiles() []ui.FileRef {
	if t.detail == nil {
		return nil
	}
	allText := []string{t.detail.Description}
	for _, c := range t.detail.Comments {
		allText = append(allText, c.Body)
	}
	return ui.ExtractFiles(allText...)
}

func (t *issueTab) ShortHelp() []key.Binding {
	switch t.state {
	case stateList:
		return []key.Binding{t.keys.Up, t.keys.Down, t.keys.Enter, t.keys.NextPage}
	case stateDetail:
		return []key.Binding{t.keys.Up, t.keys.Down, t.keys.Back, t.keys.Open, t.keys.Upload}
	case stateError:
		return []key.Binding{t.keys.Retry}
	}
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

func issueBadges(issue api.Issue) string {
	var parts []string
	if issue.Cycle != nil {
		label := fmt.Sprintf("C#%d", issue.Cycle.Number)
		parts = append(parts, cycleBadgeStyle.Render(label))
	}
	if issue.Project != nil {
		parts = append(parts, projectBadgeStyle.Render(truncate(issue.Project.Name, 15)))
	}
	if len(parts) == 0 {
		return ""
	}
	return "  " + strings.Join(parts, " ")
}

// stateIcon returns a colored dot based on workflow state type.
func stateIcon(stateType string) string {
	switch stateType {
	case "completed":
		return lipgloss.NewStyle().Foreground(greenColor).Render("●")
	case "started":
		return lipgloss.NewStyle().Foreground(yellowColor).Render("●")
	case "unstarted":
		return lipgloss.NewStyle().Foreground(cyanColor).Render("●")
	case "backlog":
		return lipgloss.NewStyle().Foreground(cyanColor).Render("○")
	case "cancelled":
		return lipgloss.NewStyle().Foreground(mutedColor).Render("○")
	default:
		return lipgloss.NewStyle().Foreground(mutedColor).Render("●")
	}
}

// priorityStyle returns a style for the priority label.
func priorityStyle(pri string) lipgloss.Style {
	switch pri {
	case "Urgent":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3B30")).Bold(true)
	case "High":
		return lipgloss.NewStyle().Foreground(pinkColor)
	case "Medium":
		return lipgloss.NewStyle().Foreground(textColor)
	case "Low":
		return lipgloss.NewStyle().Foreground(mutedColor)
	default:
		return lipgloss.NewStyle().Foreground(mutedColor)
	}
}

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3B30")).Bold(true)
	mutedStyle = lipgloss.NewStyle().Foreground(mutedColor)

	selectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1a1a3e"))

	selectedIDStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#c4a7e7"))
	selectedTitleStyle = lipgloss.NewStyle().Foreground(brightColor).Bold(true)
	accentBarStyle     = lipgloss.NewStyle().Foreground(pinkColor)

	listRowStyle      = lipgloss.NewStyle().Foreground(textColor)
	listHeaderStyle   = lipgloss.NewStyle().Foreground(mutedColor)
	listSepStyle      = lipgloss.NewStyle().Foreground(dimColor)
	detailTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(purpleColor)
	detailLabelStyle  = lipgloss.NewStyle().Foreground(yellowColor)
	cycleBadgeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#1a1a2e")).Background(cyanColor)
	projectBadgeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1a1a2e")).Background(yellowColor)
)
