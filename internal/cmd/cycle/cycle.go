package cycle

import (
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

// NewCycleCmd creates the cycle command group.
func NewCycleCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cycle",
		Short:   "Manage cycles (sprints)",
		Aliases: []string{"c"},
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newCurrentCmd(f))
	cmd.AddCommand(newViewCmd(f))
	cmd.AddCommand(newAddIssueCmd(f))
	cmd.AddCommand(newRemoveIssueCmd(f))

	return cmd
}
