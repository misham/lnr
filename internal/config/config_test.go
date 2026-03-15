package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViperStore_LoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	store := NewViperStore(dir)

	err := store.Load()
	require.NoError(t, err, "missing config file should not error")
	assert.Empty(t, store.TeamID(), "team ID should be empty when no config exists")
}

func TestViperStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := NewViperStore(dir)

	store.SetTeamID("team-123")
	err := store.Save()
	require.NoError(t, err)

	// Verify file permissions
	configPath := filepath.Join(dir, "config.yaml")
	info, err := os.Stat(configPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(), "config file should have 0600 permissions")

	// Verify directory permissions
	dirInfo, err := os.Stat(dir)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), dirInfo.Mode().Perm(), "config dir should have 0700 permissions")

	// Load in a fresh store
	store2 := NewViperStore(dir)
	err = store2.Load()
	require.NoError(t, err)
	assert.Equal(t, "team-123", store2.TeamID())
}

func TestViperStore_EnvOverride(t *testing.T) {
	dir := t.TempDir()
	store := NewViperStore(dir)

	t.Setenv("LNR_TEAM_ID", "env-team")
	err := store.Load()
	require.NoError(t, err)
	assert.Equal(t, "env-team", store.TeamID())
}

func TestViperStore_SetTeamID(t *testing.T) {
	dir := t.TempDir()
	store := NewViperStore(dir)

	assert.Empty(t, store.TeamID())
	store.SetTeamID("abc")
	assert.Equal(t, "abc", store.TeamID())
}
