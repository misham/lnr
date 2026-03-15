package cmdutil

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/auth"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

// Factory holds shared dependencies for all commands.
type Factory struct {
	Config    func() (config.Store, error)
	APIClient func() (api.Client, error)
	Auth      func() (auth.TokenStore, error)
	IO        *ui.IOStreams
	TeamKey   string // set from --team flag in PersistentPreRunE
}

// NewFactory creates a production Factory with lazy initialization.
func NewFactory() *Factory {
	ios := ui.NewIOStreams()

	configDir := filepath.Join(userConfigDir(), "lnr")

	f := &Factory{
		IO: ios,
	}

	f.Config = func() (config.Store, error) {
		store := config.NewViperStore(configDir)
		if err := store.Load(); err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
		return store, nil
	}

	f.Auth = func() (auth.TokenStore, error) {
		return auth.NewKeyringTokenStore(configDir), nil
	}

	f.APIClient = func() (api.Client, error) {
		tokenStore, err := f.Auth()
		if err != nil {
			return nil, err
		}

		clientID := clientID()
		oauthCfg := &oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{ //nolint:gosec // OAuth endpoint URL, not credentials
				TokenURL:  "https://api.linear.app/oauth/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		}

		refresher := auth.NewTokenRefresher(tokenStore, oauthCfg)
		token, err := refresher.ValidToken()
		if err != nil {
			return nil, err
		}

		return api.NewGraphQLClient("https://api.linear.app/graphql", token.AccessToken), nil
	}

	return f
}

func userConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

func clientID() string {
	if id := os.Getenv("LNR_CLIENT_ID"); id != "" {
		return id
	}
	return auth.DefaultClientID()
}
