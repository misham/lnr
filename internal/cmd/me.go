package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newMeCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show authenticated user info",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return fmt.Errorf("not authenticated: %w", err)
			}

			user, err := client.Viewer(cmd.Context())
			if err != nil {
				return fmt.Errorf("get viewer: %w", err)
			}

			return ui.PrintUser(f.IO, user)
		},
	}
}
