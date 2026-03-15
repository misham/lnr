package cycle

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var showAll bool

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List cycles",
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

			result, err := client.ListCycles(cmd.Context(), teamID, showAll, 50, "")
			if err != nil {
				return fmt.Errorf("list cycles: %w", err)
			}

			return ui.PrintCycles(f.IO, result.Cycles)
		},
	}

	cmd.Flags().BoolVar(&showAll, "all", false, "Include past cycles")

	return cmd
}
