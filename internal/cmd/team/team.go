package team

import (
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

// NewTeamCmd creates the team command group.
func NewTeamCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "Manage teams",
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newSetCmd(f))

	return cmd
}
