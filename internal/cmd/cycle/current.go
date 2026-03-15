package cycle

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newCurrentCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the active cycle",
		RunE: func(cmd *cobra.Command, args []string) error {
			teamID, err := cmdutil.ResolveTeamID(cmd.Context(), f)
			if err != nil {
				return err
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			cycle, err := client.GetActiveCycle(cmd.Context(), teamID)
			if err != nil {
				if errors.Is(err, api.ErrNoActiveCycle) {
					_, printErr := fmt.Fprintln(f.IO.Out, "No active cycle for this team")
					return printErr
				}
				return fmt.Errorf("get active cycle: %w", err)
			}

			return ui.PrintCycleDetail(f.IO, cycle)
		},
	}
}
