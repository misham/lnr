package state

import (
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

// NewStateCmd creates the state command group.
func NewStateCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "state",
		Short:   "Manage workflow states",
		Aliases: []string{"st"},
	}

	cmd.AddCommand(newListCmd(f))

	return cmd
}
