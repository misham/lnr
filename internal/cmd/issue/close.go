package issue

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
)

func newCloseCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "close <identifier>",
		Short: "Close an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			issue, err := client.GetIssue(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get issue: %w", err)
			}

			states, err := client.ListWorkflowStates(cmd.Context(), issue.Team.ID)
			if err != nil {
				return fmt.Errorf("list workflow states: %w", err)
			}

			completedState, err := findCompletedState(states)
			if err != nil {
				return err
			}

			input := api.IssueUpdateInput{
				StateID: &completedState.ID,
			}

			_, err = client.UpdateIssue(cmd.Context(), args[0], input)
			if err != nil {
				return fmt.Errorf("close issue: %w", err)
			}

			_, err = fmt.Fprintf(f.IO.Out, "Closed %s\n", args[0])
			return err
		},
	}
}

func findCompletedState(states []api.WorkflowState) (*api.WorkflowState, error) {
	var completed []api.WorkflowState
	for _, s := range states {
		if s.Type == "completed" {
			completed = append(completed, s)
		}
	}
	if len(completed) == 0 {
		return nil, fmt.Errorf("no completed workflow state found")
	}
	sort.Slice(completed, func(i, j int) bool {
		return completed[i].Position < completed[j].Position
	})
	return &completed[0], nil
}
