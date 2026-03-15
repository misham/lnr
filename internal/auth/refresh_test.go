package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func TestTokenRefresher_ValidToken(t *testing.T) {
	keyring.MockInit()
	store := NewKeyringTokenStore(t.TempDir())
	token := &oauth2.Token{
		AccessToken:  "valid-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}
	require.NoError(t, store.SaveToken(token))

	refresher := NewTokenRefresher(store, &oauth2.Config{})
	got, err := refresher.ValidToken()
	require.NoError(t, err)
	assert.Equal(t, "valid-token", got.AccessToken)
}

func TestTokenRefresher_ExpiredToken(t *testing.T) {
	// Fake token endpoint that returns refreshed tokens
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"access_token":  "new-access",
			"refresh_token": "new-refresh",
			"token_type":    "Bearer",
			"expires_in":    3600,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	keyring.MockInit()
	store := NewKeyringTokenStore(t.TempDir())
	token := &oauth2.Token{
		AccessToken:  "expired-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(-time.Hour),
	}
	require.NoError(t, store.SaveToken(token))

	cfg := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: srv.URL,
		},
	}
	refresher := NewTokenRefresher(store, cfg)
	got, err := refresher.ValidToken()
	require.NoError(t, err)
	assert.Equal(t, "new-access", got.AccessToken)

	// Verify the new token was persisted
	persisted, err := store.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "new-access", persisted.AccessToken)
	assert.Equal(t, "new-refresh", persisted.RefreshToken)
}

func TestTokenRefresher_NoToken(t *testing.T) {
	keyring.MockInit()
	store := NewKeyringTokenStore(t.TempDir())
	refresher := NewTokenRefresher(store, &oauth2.Config{})
	_, err := refresher.ValidToken()
	require.ErrorIs(t, err, ErrNotAuthenticated)
}
