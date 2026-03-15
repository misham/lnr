package cmdutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestNewFactory(t *testing.T) {
	keyring.MockInit()
	f := NewFactory()
	require.NotNil(t, f)
	assert.NotNil(t, f.IO)
	assert.NotNil(t, f.Config)
	assert.NotNil(t, f.Auth)
	assert.NotNil(t, f.APIClient)
}

func TestNewFactory_ConfigLoads(t *testing.T) {
	keyring.MockInit()
	f := NewFactory()
	store, err := f.Config()
	require.NoError(t, err)
	assert.NotNil(t, store)
}

func TestNewFactory_AuthReturnsTokenStore(t *testing.T) {
	keyring.MockInit()
	f := NewFactory()
	tokenStore, err := f.Auth()
	require.NoError(t, err)
	assert.NotNil(t, tokenStore)
}

func TestNewFactory_APIClientFailsWithoutToken(t *testing.T) {
	keyring.MockInit()
	f := NewFactory()
	_, err := f.APIClient()
	require.Error(t, err)
}

func TestClientID_EnvOverride(t *testing.T) {
	t.Setenv("LNR_CLIENT_ID", "test-client-id")
	got := clientID()
	assert.Equal(t, "test-client-id", got)
}

func TestClientID_FallsBackToDefault(t *testing.T) {
	t.Setenv("LNR_CLIENT_ID", "")
	got := clientID()
	// DefaultClientID() returns empty string when not set by ldflags
	assert.Equal(t, "", got)
}

func TestUserConfigDir(t *testing.T) {
	dir := userConfigDir()
	assert.NotEmpty(t, dir)

	// Should resolve to a real path
	_, err := os.Stat(dir)
	assert.NoError(t, err)
}
