package cycle

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
)

func newAddIssueCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "add-issue <cycle-id-or-number> <issue-id>",
		Short: "Add an issue to a cycle",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			cycleID, err := resolveCycleID(cmd, f, client, args[0])
			if err != nil {
				return err
			}

			issueID := args[1]
			_, err = client.UpdateIssue(cmd.Context(), issueID, api.IssueUpdateInput{
				CycleID: &cycleID,
			})
			if err != nil {
				return fmt.Errorf("add issue to cycle: %w", err)
			}

			_, err = fmt.Fprintf(f.IO.Out, "Added issue %s to cycle\n", issueID)
			return err
		},
	}
}

func resolveCycleID(cmd *cobra.Command, f *cmdutil.Factory, client api.Client, arg string) (string, error) {
	if num, parseErr := strconv.Atoi(arg); parseErr == nil {
		teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
		if err != nil {
			return "", err
		}
		cycle, err := client.GetCycleByNumber(cmd.Context(), teamID, num)
		if err != nil {
			return "", fmt.Errorf("get cycle #%d: %w", num, err)
		}
		return cycle.ID, nil
	}
	return arg, nil
}
