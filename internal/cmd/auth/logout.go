package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

func newLogoutCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear stored authentication tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenStore, err := f.Auth()
			if err != nil {
				return err
			}

			if err := tokenStore.ClearToken(); err != nil {
				return fmt.Errorf("clear token: %w", err)
			}

			_, _ = fmt.Fprintln(f.IO.Out, "Logged out successfully.")
			return nil
		},
	}
}
