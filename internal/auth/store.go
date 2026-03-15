package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

// ErrNotAuthenticated is returned when no token is stored.
var ErrNotAuthenticated = errors.New("not authenticated")

// TokenStore abstracts token persistence.
type TokenStore interface {
	SaveToken(token *oauth2.Token) error
	LoadToken() (*oauth2.Token, error)
	ClearToken() error
	HasToken() bool
}

const (
	keyringService = "lnr-cli"
	keyringKey     = "token"
)

// KeyringTokenStore stores tokens in the OS keyring.
type KeyringTokenStore struct {
	configDir string // needed for migration from file-based storage
}

// NewKeyringTokenStore creates a keyring-backed token store.
func NewKeyringTokenStore(configDir string) *KeyringTokenStore {
	return &KeyringTokenStore{configDir: configDir}
}

func (s *KeyringTokenStore) SaveToken(token *oauth2.Token) error {
	data, err := json.Marshal(token) //nolint:gosec // intentional: persisting OAuth token to keyring
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}
	if err := keyring.Set(keyringService, keyringKey, string(data)); err != nil {
		return fmt.Errorf("save token to keyring: %w", err)
	}
	return nil
}

func (s *KeyringTokenStore) LoadToken() (*oauth2.Token, error) {
	data, err := keyring.Get(keyringService, keyringKey)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return s.migrateFromFile()
		}
		return nil, fmt.Errorf("load token from keyring: %w", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("unmarshal token: %w", err)
	}
	return &token, nil
}

func (s *KeyringTokenStore) migrateFromFile() (*oauth2.Token, error) {
	filePath := filepath.Join(s.configDir, "token.json")
	data, err := os.ReadFile(filePath) //nolint:gosec // path constructed from known configDir
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotAuthenticated
		}
		return nil, fmt.Errorf("read token file for migration: %w", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("unmarshal token file for migration: %w", err)
	}

	if err := s.SaveToken(&token); err != nil {
		return nil, fmt.Errorf("migrate token to keyring: %w", err)
	}

	_ = os.Remove(filePath)

	return &token, nil
}

func (s *KeyringTokenStore) ClearToken() error {
	err := keyring.Delete(keyringService, keyringKey)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("clear token from keyring: %w", err)
	}
	return nil
}

func (s *KeyringTokenStore) HasToken() bool {
	_, err := keyring.Get(keyringService, keyringKey)
	return err == nil
}
