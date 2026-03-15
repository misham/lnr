package cycle

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view <id-or-number>",
		Short: "View cycle details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			arg := args[0]

			// If numeric, look up by cycle number within the team.
			if num, parseErr := strconv.Atoi(arg); parseErr == nil {
				teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
				if err != nil {
					return err
				}
				cycle, err := client.GetCycleByNumber(cmd.Context(), teamID, num)
				if err != nil {
					return fmt.Errorf("get cycle #%d: %w", num, err)
				}
				return ui.PrintCycleDetail(f.IO, cycle)
			}

			// Otherwise treat as UUID.
			cycle, err := client.GetCycle(cmd.Context(), arg)
			if err != nil {
				return fmt.Errorf("get cycle: %w", err)
			}
			return ui.PrintCycleDetail(f.IO, cycle)
		},
	}
}
