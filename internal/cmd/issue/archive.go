package issue

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

func newArchiveCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "archive <identifier>",
		Short: "Archive an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if err := client.ArchiveIssue(cmd.Context(), args[0]); err != nil {
				return fmt.Errorf("archive issue: %w", err)
			}

			_, err = fmt.Fprintf(f.IO.Out, "Archived %s\n", args[0])
			return err
		},
	}
}
