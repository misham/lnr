package auth

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	internalAuth "github.com/misham/linear-cli/internal/auth"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestLogoutCmd(t *testing.T) {
	keyring.MockInit()
	dir := t.TempDir()
	ios := ui.NewTestIOStreams()

	f := &cmdutil.Factory{
		IO: ios,
		Auth: func() (internalAuth.TokenStore, error) {
			return internalAuth.NewKeyringTokenStore(dir), nil
		},
	}

	cmd := newLogoutCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, ok := ios.Out.(*bytes.Buffer)
	require.True(t, ok)
	assert.Contains(t, buf.String(), "Logged out")
}
