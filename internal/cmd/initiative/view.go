package initiative

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view <name-or-id>",
		Short: "View initiative details",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				_ = cmd.Usage()
				return fmt.Errorf("requires exactly 1 arg: initiative name or id")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			initiative, err := resolveInitiative(cmd.Context(), client, args[0])
			if err != nil {
				return err
			}
			return ui.PrintInitiativeDetail(f.IO, initiative)
		},
	}
}
