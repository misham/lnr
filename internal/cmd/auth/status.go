package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

func newStatusCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return fmt.Errorf("not authenticated: %w", err)
			}

			user, err := client.Viewer(cmd.Context())
			if err != nil {
				return fmt.Errorf("get viewer: %w", err)
			}

			_, _ = fmt.Fprintf(f.IO.Out, "Logged in as %s (%s)\n", user.Name, user.Email)
			return nil
		},
	}
}
