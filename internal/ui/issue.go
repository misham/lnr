package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

var (
	identifierStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	stateStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	issueLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	commentStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
)

// errWriter wraps an io.Writer and captures the first error from any write.
type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) printf(format string, a ...any) {
	if ew.err != nil {
		return
	}
	_, ew.err = fmt.Fprintf(ew.w, format, a...)
}

// PrintIssues formats and prints a list of issues.
func PrintIssues(ios *IOStreams, issues []api.Issue) error {
	if len(issues) == 0 {
		_, err := fmt.Fprintln(ios.Out, "No issues found")
		return err
	}

	if ios.IsPlain() {
		return printIssuesPlain(ios, issues)
	}
	return printIssuesStyled(ios, issues)
}

func assigneeName(u *api.User) string {
	if u == nil {
		return ""
	}
	return u.Name
}

func printIssuesPlain(ios *IOStreams, issues []api.Issue) error {
	for _, issue := range issues {
		if _, err := fmt.Fprintf(ios.Out, "%s\t%s\t%s\t%s\t%s\n",
			issue.Identifier,
			issue.Title,
			issue.State.Name,
			issue.PriorityLabel,
			assigneeName(issue.Assignee),
		); err != nil {
			return err
		}
	}
	return nil
}

// padRight pads a styled string to the given visible width.
// ANSI escape codes are not counted toward width.
func padRight(styled string, visibleLen, width int) string {
	if pad := width - visibleLen; pad > 0 {
		return styled + strings.Repeat(" ", pad)
	}
	return styled
}

func printIssuesStyled(ios *IOStreams, issues []api.Issue) error {
	// Calculate column widths from plain text to avoid ANSI escape codes
	// throwing off alignment.
	idW, titleW, stateW, prioW := len("ID"), len("TITLE"), len("STATE"), len("PRIORITY")
	for _, issue := range issues {
		if l := len(issue.Identifier); l > idW {
			idW = l
		}
		if l := len(issue.Title); l > titleW {
			titleW = l
		}
		if l := len(issue.State.Name); l > stateW {
			stateW = l
		}
		if l := len(issue.PriorityLabel); l > prioW {
			prioW = l
		}
	}

	const gap = "  "
	ew := &errWriter{w: ios.Out}
	ew.printf("%s%s%s%s%s%s%s%s%s\n",
		padRight(headerStyle.Render("ID"), len("ID"), idW), gap,
		padRight(headerStyle.Render("TITLE"), len("TITLE"), titleW), gap,
		padRight(headerStyle.Render("STATE"), len("STATE"), stateW), gap,
		padRight(headerStyle.Render("PRIORITY"), len("PRIORITY"), prioW), gap,
		headerStyle.Render("ASSIGNEE"),
	)
	for _, issue := range issues {
		ew.printf("%s%s%s%s%s%s%s%s%s\n",
			padRight(identifierStyle.Render(issue.Identifier), len(issue.Identifier), idW), gap,
			padRight(issue.Title, len(issue.Title), titleW), gap,
			padRight(stateStyle.Render(issue.State.Name), len(issue.State.Name), stateW), gap,
			padRight(issue.PriorityLabel, len(issue.PriorityLabel), prioW), gap,
			assigneeName(issue.Assignee),
		)
	}
	return ew.err
}

// PrintIssueDetail formats and prints a single issue with full detail.
func PrintIssueDetail(ios *IOStreams, issue *api.Issue) error {
	if ios.IsPlain() {
		return printIssueDetailPlain(ios, issue)
	}
	return printIssueDetailStyled(ios, issue)
}

func printIssueDetailPlain(ios *IOStreams, issue *api.Issue) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("%s  %s\n", issue.Identifier, issue.Title)
	ew.printf("State:     %s\n", issue.State.Name)
	ew.printf("Priority:  %s\n", issue.PriorityLabel)
	ew.printf("Team:      %s\n", issue.Team.Name)

	if issue.Assignee != nil {
		ew.printf("Assignee:  %s\n", issue.Assignee.Name)
	} else {
		ew.printf("Assignee:  Unassigned\n")
	}

	if issue.Estimate > 0 {
		ew.printf("Estimate:  %d\n", issue.Estimate)
	}
	if issue.DueDate != "" {
		ew.printf("Due:       %s\n", issue.DueDate)
	}

	if len(issue.Labels) > 0 {
		names := make([]string, len(issue.Labels))
		for i, l := range issue.Labels {
			names[i] = l.Name
		}
		ew.printf("Labels:    %s\n", strings.Join(names, ", "))
	}

	ew.printf("URL:       %s\n", issue.URL)
	ew.printf("Created:   %s\n", issue.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("Updated:   %s\n", issue.UpdatedAt.Format("2006-01-02 15:04"))

	if issue.CompletedAt != nil {
		ew.printf("Completed: %s\n", issue.CompletedAt.Format("2006-01-02 15:04"))
	}

	if issue.Description != "" {
		ew.printf("\n%s\n", issue.Description)
	}

	if len(issue.Comments) > 0 {
		ew.printf("\nComments (%d):\n", len(issue.Comments))
		for _, c := range issue.Comments {
			author := "unknown"
			if c.User != nil {
				author = c.User.Name
			}
			ew.printf("  %s (%s):\n", author, c.CreatedAt.Format("2006-01-02 15:04"))
			ew.printf("    %s\n", c.Body)
		}
	}

	files := issueFiles(issue)
	if len(files) > 0 {
		ew.printf("\n")
		for _, f := range files {
			ew.printf("%s\t%s\n", f.Name, f.URL)
		}
	}

	return ew.err
}

func printIssueDetailStyled(ios *IOStreams, issue *api.Issue) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("%s  %s\n", identifierStyle.Render(issue.Identifier), lipgloss.NewStyle().Bold(true).Render(issue.Title))
	ew.printf("%s  %s\n", headerStyle.Render("State:"), stateStyle.Render(issue.State.Name))
	ew.printf("%s  %s\n", headerStyle.Render("Priority:"), issue.PriorityLabel)
	ew.printf("%s  %s\n", headerStyle.Render("Team:"), issue.Team.Name)

	if issue.Assignee != nil {
		ew.printf("%s  %s\n", headerStyle.Render("Assignee:"), issue.Assignee.Name)
	} else {
		ew.printf("%s  Unassigned\n", headerStyle.Render("Assignee:"))
	}

	if issue.Estimate > 0 {
		ew.printf("%s  %d\n", headerStyle.Render("Estimate:"), issue.Estimate)
	}
	if issue.DueDate != "" {
		ew.printf("%s  %s\n", headerStyle.Render("Due:"), issue.DueDate)
	}

	if len(issue.Labels) > 0 {
		names := make([]string, len(issue.Labels))
		for i, l := range issue.Labels {
			names[i] = issueLabelStyle.Render(l.Name)
		}
		ew.printf("%s  %s\n", headerStyle.Render("Labels:"), strings.Join(names, ", "))
	}

	ew.printf("%s  %s\n", headerStyle.Render("URL:"), issue.URL)
	ew.printf("%s  %s\n", headerStyle.Render("Created:"), issue.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("%s  %s\n", headerStyle.Render("Updated:"), issue.UpdatedAt.Format("2006-01-02 15:04"))

	if issue.CompletedAt != nil {
		ew.printf("%s  %s\n", headerStyle.Render("Completed:"), issue.CompletedAt.Format("2006-01-02 15:04"))
	}

	if issue.Description != "" {
		ew.printf("\n%s\n", issue.Description)
	}

	if len(issue.Comments) > 0 {
		ew.printf("\n%s\n", headerStyle.Render(fmt.Sprintf("Comments (%d):", len(issue.Comments))))
		for _, c := range issue.Comments {
			author := "unknown"
			if c.User != nil {
				author = c.User.Name
			}
			ew.printf("  %s\n", commentStyle.Render(fmt.Sprintf("%s (%s):", author, c.CreatedAt.Format("2006-01-02 15:04"))))
			ew.printf("    %s\n", c.Body)
		}
	}

	files := issueFiles(issue)
	if len(files) > 0 {
		ew.printf("\n%s\n", headerStyle.Render(fmt.Sprintf("Files (%d):", len(files))))
		for _, f := range files {
			ew.printf("  %s\n", HyperlinkOSC8(f.URL, f.Name))
		}
	}

	return ew.err
}

func issueFiles(issue *api.Issue) []FileRef {
	allText := []string{issue.Description}
	for _, c := range issue.Comments {
		allText = append(allText, c.Body)
	}
	return ExtractFiles(allText...)
}

// PrintIssueCreated prints a confirmation after creating an issue.
func PrintIssueCreated(ios *IOStreams, issue *api.Issue) error {
	if ios.IsPlain() {
		_, err := fmt.Fprintf(ios.Out, "Created %s %s\n%s\n", issue.Identifier, issue.Title, issue.URL)
		return err
	}
	_, err := fmt.Fprintf(ios.Out, "Created %s  %s\n%s\n",
		identifierStyle.Render(issue.Identifier),
		issue.Title,
		issue.URL,
	)
	return err
}

// PrintIssueUpdated prints a confirmation after updating an issue.
func PrintIssueUpdated(ios *IOStreams, issue *api.Issue) error {
	if ios.IsPlain() {
		_, err := fmt.Fprintf(ios.Out, "Updated %s %s\n%s\n", issue.Identifier, issue.Title, issue.URL)
		return err
	}
	_, err := fmt.Fprintf(ios.Out, "Updated %s  %s\n%s\n",
		identifierStyle.Render(issue.Identifier),
		issue.Title,
		issue.URL,
	)
	return err
}
