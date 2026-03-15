package auth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/misham/linear-cli/internal/auth"
	"github.com/misham/linear-cli/internal/cmdutil"
)

func newLoginCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Linear",
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenStore, err := f.Auth()
			if err != nil {
				return err
			}

			clientID := os.Getenv("LNR_CLIENT_ID")
			if clientID == "" {
				clientID = auth.DefaultClientID()
			}

			oauthCfg := &oauth2.Config{
				ClientID: clientID,
				Endpoint: oauth2.Endpoint{ //nolint:gosec // OAuth endpoint URLs, not credentials
					AuthURL:   "https://linear.app/oauth/authorize",
					TokenURL:  "https://api.linear.app/oauth/token",
					AuthStyle: oauth2.AuthStyleInParams,
				},
				RedirectURL: "http://localhost:8585/callback",
			}

			flow := &auth.LoginFlow{
				OAuthConfig: oauthCfg,
				TokenStore:  tokenStore,
				Browser:     &auth.SystemBrowserOpener{},
			}

			if err := flow.Run(cmd.Context()); err != nil {
				return err
			}

			_, _ = fmt.Fprintln(f.IO.Out, "Login successful!")
			return nil
		},
	}
}
