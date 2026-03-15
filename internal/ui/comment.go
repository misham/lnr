package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

// PrintComments formats and prints a list of comments.
func PrintComments(ios *IOStreams, comments []api.Comment) error {
	if len(comments) == 0 {
		_, err := fmt.Fprintln(ios.Out, "No comments")
		return err
	}

	if ios.IsPlain() {
		return printCommentsPlain(ios, comments)
	}
	return printCommentsStyled(ios, comments)
}

func commentAuthor(c api.Comment) string {
	if c.User == nil {
		return "unknown"
	}
	return c.User.Name
}

func printCommentsPlain(ios *IOStreams, comments []api.Comment) error {
	for _, c := range comments {
		if _, err := fmt.Fprintf(ios.Out, "%s\t%s\t%s\n",
			commentAuthor(c),
			c.CreatedAt.Format("2006-01-02 15:04"),
			c.Body,
		); err != nil {
			return err
		}
	}
	return nil
}

func printCommentsStyled(ios *IOStreams, comments []api.Comment) error {
	authorStyle := lipgloss.NewStyle().Bold(true)
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))

	ew := &errWriter{w: ios.Out}
	for i, c := range comments {
		if i > 0 {
			ew.printf("\n")
		}
		ew.printf("%s  %s\n", authorStyle.Render(commentAuthor(c)), timeStyle.Render(c.CreatedAt.Format("2006-01-02 15:04")))
		ew.printf("  %s\n", c.Body)
	}
	return ew.err
}
