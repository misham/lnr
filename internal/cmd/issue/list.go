package issue

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		limit    int
		state    string
		assignee string
		label    string
		priority int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List issues",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
			if err != nil {
				return err
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			result, err := client.ListIssues(cmd.Context(), teamID, limit, "")
			if err != nil {
				return fmt.Errorf("list issues: %w", err)
			}

			issues := filterIssues(result.Issues, state, assignee, label, priority)
			return ui.PrintIssues(f.IO, issues)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Max issues to fetch")
	cmd.Flags().StringVar(&state, "state", "", "Filter by state type (e.g. started, backlog)")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee name or email")
	cmd.Flags().StringVar(&label, "label", "", "Filter by label name")
	cmd.Flags().IntVar(&priority, "priority", 0, "Filter by priority level (1-4)")

	return cmd
}

func filterIssues(issues []api.Issue, state, assignee, label string, priority int) []api.Issue {
	if state == "" && assignee == "" && label == "" && priority == 0 {
		return issues
	}

	var filtered []api.Issue
	for _, issue := range issues {
		if state != "" && !strings.EqualFold(issue.State.Type, state) {
			continue
		}
		if assignee != "" && !matchAssignee(issue.Assignee, assignee) {
			continue
		}
		if label != "" && !hasLabel(issue.Labels, label) {
			continue
		}
		if priority != 0 && issue.Priority != priority {
			continue
		}
		filtered = append(filtered, issue)
	}
	return filtered
}

func matchAssignee(u *api.User, filter string) bool {
	if u == nil {
		return false
	}
	return strings.EqualFold(u.Name, filter) || strings.EqualFold(u.Email, filter)
}

func hasLabel(labels []api.IssueLabel, name string) bool {
	for _, l := range labels {
		if strings.EqualFold(l.Name, name) {
			return true
		}
	}
	return false
}
