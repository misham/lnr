package auth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

// TokenRefresher provides valid tokens, refreshing as needed.
type TokenRefresher struct {
	store TokenStore
	cfg   *oauth2.Config
}

// NewTokenRefresher creates a refresher backed by the given store and config.
func NewTokenRefresher(store TokenStore, cfg *oauth2.Config) *TokenRefresher {
	return &TokenRefresher{store: store, cfg: cfg}
}

// ValidToken returns a valid token, refreshing it if expired.
func (r *TokenRefresher) ValidToken() (*oauth2.Token, error) {
	token, err := r.store.LoadToken()
	if err != nil {
		return nil, err
	}

	if token.Valid() {
		return token, nil
	}

	// Token expired — use refresh token to get a new one
	src := r.cfg.TokenSource(context.Background(), token)
	newToken, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}

	if err := r.store.SaveToken(newToken); err != nil {
		return nil, fmt.Errorf("save refreshed token: %w", err)
	}

	return newToken, nil
}
