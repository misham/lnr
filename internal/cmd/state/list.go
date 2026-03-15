package state

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List workflow states",
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

			states, err := client.ListWorkflowStates(cmd.Context(), teamID)
			if err != nil {
				return fmt.Errorf("list states: %w", err)
			}

			sort.Slice(states, func(i, j int) bool {
				return states[i].Position < states[j].Position
			})

			return ui.PrintStates(f.IO, states)
		},
	}

	return cmd
}
