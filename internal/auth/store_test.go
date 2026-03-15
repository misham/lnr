package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func TestKeyringTokenStore_NoToken(t *testing.T) {
	keyring.MockInit()
	store := NewKeyringTokenStore(t.TempDir())

	assert.False(t, store.HasToken())

	_, err := store.LoadToken()
	require.ErrorIs(t, err, ErrNotAuthenticated)
}

func TestKeyringTokenStore_SaveAndLoad(t *testing.T) {
	keyring.MockInit()
	store := NewKeyringTokenStore(t.TempDir())

	token := &oauth2.Token{
		AccessToken:  "access-123",
		RefreshToken: "refresh-456",
		TokenType:    "Bearer",
		Expiry:       time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC),
	}

	err := store.SaveToken(token)
	require.NoError(t, err)
	assert.True(t, store.HasToken())

	loaded, err := store.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "access-123", loaded.AccessToken)
	assert.Equal(t, "refresh-456", loaded.RefreshToken)
	assert.Equal(t, "Bearer", loaded.TokenType)
	assert.True(t, loaded.Expiry.Equal(token.Expiry))
}

func TestKeyringTokenStore_Clear(t *testing.T) {
	keyring.MockInit()
	store := NewKeyringTokenStore(t.TempDir())

	token := &oauth2.Token{
		AccessToken:  "access-123",
		RefreshToken: "refresh-456",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	err := store.SaveToken(token)
	require.NoError(t, err)
	assert.True(t, store.HasToken())

	err = store.ClearToken()
	require.NoError(t, err)
	assert.False(t, store.HasToken())
}

func TestKeyringTokenStore_KeyringError(t *testing.T) {
	keyring.MockInitWithError(errors.New("keyring locked"))
	store := NewKeyringTokenStore(t.TempDir())

	token := &oauth2.Token{
		AccessToken: "access-123",
		TokenType:   "Bearer",
	}

	err := store.SaveToken(token)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "keyring locked")

	_, err = store.LoadToken()
	require.Error(t, err)

	err = store.ClearToken()
	require.Error(t, err)
}

func TestKeyringTokenStore_MigratesFromFile(t *testing.T) {
	keyring.MockInit()
	dir := t.TempDir()

	token := &oauth2.Token{
		AccessToken:  "migrated-access",
		RefreshToken: "migrated-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC),
	}
	data, err := json.Marshal(token) //nolint:gosec // test token
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "token.json"), data, 0o600))

	store := NewKeyringTokenStore(dir)

	// LoadToken should migrate from file to keyring
	loaded, err := store.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "migrated-access", loaded.AccessToken)
	assert.Equal(t, "migrated-refresh", loaded.RefreshToken)

	// File should be deleted after migration
	_, err = os.Stat(filepath.Join(dir, "token.json"))
	assert.True(t, os.IsNotExist(err))

	// Subsequent load should use keyring (no file needed)
	loaded2, err := store.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "migrated-access", loaded2.AccessToken)
}

func TestKeyringTokenStore_NoMigrationWhenKeyringHasToken(t *testing.T) {
	keyring.MockInit()
	dir := t.TempDir()

	// Pre-populate keyring
	store := NewKeyringTokenStore(dir)
	keyringToken := &oauth2.Token{
		AccessToken:  "keyring-access",
		RefreshToken: "keyring-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC),
	}
	require.NoError(t, store.SaveToken(keyringToken))

	// Also create a file token
	fileToken := &oauth2.Token{
		AccessToken:  "file-access",
		RefreshToken: "file-refresh",
		TokenType:    "Bearer",
	}
	data, err := json.Marshal(fileToken) //nolint:gosec // test token
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "token.json"), data, 0o600))

	// LoadToken should use keyring, NOT the file
	loaded, err := store.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "keyring-access", loaded.AccessToken)

	// File should NOT be deleted (keyring took precedence, no migration needed)
	_, err = os.Stat(filepath.Join(dir, "token.json"))
	assert.NoError(t, err)
}
