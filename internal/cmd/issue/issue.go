package issue

import (
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

// NewIssueCmd creates the issue command group.
func NewIssueCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "Manage issues",
		Aliases: []string{"i"},
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newViewCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newCloseCmd(f))
	cmd.AddCommand(newArchiveCmd(f))
	cmd.AddCommand(newSearchCmd(f))
	cmd.AddCommand(newCommentCmd(f))
	cmd.AddCommand(newLabelCmd(f))
	cmd.AddCommand(newFileCmd(f))

	return cmd
}
