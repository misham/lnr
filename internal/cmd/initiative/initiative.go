package initiative

import (
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

// NewInitiativeCmd creates the initiative command group.
func NewInitiativeCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "initiative",
		Short:   "View initiatives",
		Aliases: []string{"ini"},
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newViewCmd(f))
	cmd.AddCommand(newProjectsCmd(f))

	return cmd
}
