package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	internalAuth "github.com/misham/linear-cli/internal/auth"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestLoginCmd_NoClientID(t *testing.T) {
	keyring.MockInit()
	dir := t.TempDir()
	ios := ui.NewTestIOStreams()

	t.Setenv("LNR_CLIENT_ID", "")

	f := &cmdutil.Factory{
		IO: ios,
		Auth: func() (internalAuth.TokenStore, error) {
			return internalAuth.NewKeyringTokenStore(dir), nil
		},
	}

	cmd := newLoginCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "client ID")
}
