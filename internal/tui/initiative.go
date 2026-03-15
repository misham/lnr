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

type initiativeTab struct {
	client         api.Client
	ctx            context.Context
	state          tabState
	items          []api.Initiative
	detail         *api.Initiative
	detailProjects []api.Project
	cursor         int
	offset         int
	pageInfo       api.PageInfo
	size           tea.WindowSizeMsg
	spinner        spinner.Model
	viewport       viewport.Model
	errMsg         string
	keys           KeyMap
	loadingMore    bool
}

// NewInitiativeTab creates a new initiative tab.
func NewInitiativeTab(ctx context.Context, client api.Client) *initiativeTab {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &initiativeTab{
		client:  client,
		ctx:     ctx,
		state:   stateLoading,
		spinner: s,
		keys:    DefaultKeyMap(),
	}
}

func (t *initiativeTab) Init() tea.Cmd {
	return tea.Batch(
		t.spinner.Tick,
		loadInitiatives(t.ctx, t.client, defaultPageSize, ""),
	)
}

func (t *initiativeTab) Update(msg tea.Msg) (TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.size = msg
		t.viewport.Width = msg.Width
		t.viewport.Height = msg.Height
		return t, nil

	case InitiativesLoadedMsg:
		if msg.Err != nil {
			t.state = stateError
			t.errMsg = msg.Err.Error()
			t.loadingMore = false
			return t, nil
		}
		prevLen := len(t.items)
		t.items = append(t.items, msg.Initiatives...)
		t.pageInfo = msg.PageInfo
		t.state = stateList
		if t.loadingMore && prevLen > 0 {
			t.cursor = prevLen
			t.adjustOffset()
		}
		t.loadingMore = false
		return t, nil

	case InitiativeDetailLoadedMsg:
		if msg.Err != nil {
			t.state = stateError
			t.errMsg = msg.Err.Error()
			return t, nil
		}
		t.detail = msg.Initiative
		t.detailProjects = nil
		t.state = stateDetail
		t.viewport.SetContent(t.renderDetail())
		t.viewport.GotoTop()
		return t, loadInitiativeProjects(t.ctx, t.client, msg.Initiative.ID, defaultPageSize, "")

	case InitiativeProjectsLoadedMsg:
		if msg.Err != nil {
			return t, nil
		}
		t.detailProjects = msg.Projects
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

func (t *initiativeTab) handleKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
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
				loadInitiatives(t.ctx, t.client, defaultPageSize, ""),
			)
		}
		return t, nil
	}
	return t, nil
}

func (t *initiativeTab) handleListKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
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
			return t, loadInitiativeDetail(t.ctx, t.client, t.items[t.cursor].ID)
		}
	case key.Matches(msg, t.keys.NextPage):
		if t.pageInfo.HasNextPage {
			t.loadingMore = true
			return t, loadInitiatives(t.ctx, t.client, defaultPageSize, t.pageInfo.EndCursor)
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

func (t *initiativeTab) handleDetailKey(msg tea.KeyMsg) (TabModel, tea.Cmd) {
	switch {
	case key.Matches(msg, t.keys.Back):
		t.state = stateList
		t.detail = nil
		t.detailProjects = nil
		return t, nil
	case key.Matches(msg, t.keys.Down):
		t.viewport.ScrollDown(1)
	case key.Matches(msg, t.keys.Up):
		t.viewport.ScrollUp(1)
	}
	return t, nil
}

func (t *initiativeTab) visibleRows() int {
	// -1 for header, -1 for separator, -1 for pagination hint.
	return max(t.size.Height-3, 1)
}

func (t *initiativeTab) adjustOffset() {
	visible := t.visibleRows()
	if t.cursor < t.offset {
		t.offset = t.cursor
	}
	if t.cursor >= t.offset+visible {
		t.offset = t.cursor - visible + 1
	}
}

func (t *initiativeTab) View() string {
	switch t.state {
	case stateLoading:
		return t.spinner.View() + " Loading initiatives..."
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

func (t *initiativeTab) renderList() string {
	if len(t.items) == 0 {
		return "No initiatives found"
	}

	var b strings.Builder

	header := fmt.Sprintf("   %-30s %s", "Name", "Status")
	b.WriteString(listHeaderStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(listSepStyle.Render(strings.Repeat("─", max(t.size.Width-4, 40))))
	b.WriteString("\n")

	visible := t.visibleRows()
	end := min(t.offset+visible, len(t.items))

	for i := t.offset; i < end; i++ {
		initiative := t.items[i]

		var statusDot string
		if initiative.Status == "Active" {
			statusDot = lipgloss.NewStyle().Foreground(greenColor).Render("●")
		} else {
			statusDot = lipgloss.NewStyle().Foreground(cyanColor).Render("●")
		}

		stateCol := fmt.Sprintf("%s %s", statusDot, initiative.Status)

		if i == t.cursor {
			name := fmt.Sprintf("%-30s", truncate(initiative.Name, 28))
			line := fmt.Sprintf(" %s %s %s",
				accentBarStyle.Render("▎"),
				selectedTitleStyle.Render(name),
				stateCol,
			)
			b.WriteString(selectedRowStyle.Width(t.size.Width).Render(line))
		} else {
			line := fmt.Sprintf("   %-30s %s",
				truncate(initiative.Name, 28),
				stateCol,
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

func (t *initiativeTab) renderDetail() string {
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

	fmt.Fprintf(&b, "  Status:  %s\n", d.Status)
	if d.TargetDate != "" {
		fmt.Fprintf(&b, "  Target:  %s\n", d.TargetDate)
	}
	if d.Owner != nil {
		fmt.Fprintf(&b, "  Owner:   %s\n", d.Owner.DisplayName)
	}
	if d.URL != "" {
		fmt.Fprintf(&b, "  URL:     %s\n", d.URL)
	}

	if len(t.detailProjects) > 0 {
		b.WriteString("\n")
		fmt.Fprintf(&b, "Linked Projects (%d)\n", len(t.detailProjects))
		for _, p := range t.detailProjects {
			fmt.Fprintf(&b, "  - %s (%s)\n", p.Name, p.Status.Name)
		}
	}

	return b.String()
}

func (t *initiativeTab) ShortHelp() []key.Binding {
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
