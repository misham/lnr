package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
	"go.uber.org/goleak"
	"golang.org/x/oauth2"
)

func TestMain(m *testing.M) {
	keyring.MockInit()
	goleak.VerifyTestMain(m)
}

func TestLogin_Success(t *testing.T) {
	keyring.MockInit()

	// Fake OAuth token endpoint
	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"access_token":  "access-abc",
			"refresh_token": "refresh-def",
			"token_type":    "Bearer",
			"expires_in":    3600,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer tokenSrv.Close()

	store := NewKeyringTokenStore(t.TempDir())
	browser := &FakeBrowserOpener{}

	flow := &LoginFlow{
		OAuthConfig: &oauth2.Config{
			ClientID: "test-client",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://example.com/authorize",
				TokenURL: tokenSrv.URL,
			},
			RedirectURL: "http://localhost:0/callback",
		},
		TokenStore: store,
		Browser:    browser,
		ListenAddr: "localhost:0",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run login in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- flow.Run(ctx)
	}()

	// Wait for the flow to open the browser
	require.Eventually(t, func() bool {
		return browser.URL() != ""
	}, 5*time.Second, 10*time.Millisecond)

	// Parse the auth URL to get state and the callback port
	authURL, err := url.Parse(browser.URL())
	require.NoError(t, err)
	state := authURL.Query().Get("state")
	require.NotEmpty(t, state)

	// Simulate OAuth callback — retry briefly since the server goroutine
	// may not be serving yet when browser URL is set.
	callbackURL := flow.CallbackURL() + "?code=auth-code-123&state=" + state
	var resp *http.Response
	require.Eventually(t, func() bool {
		var getErr error
		resp, getErr = http.Get(callbackURL) //nolint:gosec // test URL
		return getErr == nil
	}, 2*time.Second, 10*time.Millisecond, "callback server never became ready")
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Wait for flow to complete
	err = <-errCh
	require.NoError(t, err)

	// Verify token was stored
	token, err := store.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "access-abc", token.AccessToken)
	assert.Equal(t, "refresh-def", token.RefreshToken)
}

func TestLogin_StateMismatch(t *testing.T) {
	keyring.MockInit()

	store := NewKeyringTokenStore(t.TempDir())
	browser := &FakeBrowserOpener{}

	flow := &LoginFlow{
		OAuthConfig: &oauth2.Config{
			ClientID: "test-client",
			Endpoint: oauth2.Endpoint{ //nolint:gosec // test values
				AuthURL:  "https://example.com/authorize",
				TokenURL: "https://example.com/token",
			},
			RedirectURL: "http://localhost:0/callback",
		},
		TokenStore: store,
		Browser:    browser,
		ListenAddr: "localhost:0",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- flow.Run(ctx)
	}()

	require.Eventually(t, func() bool {
		return browser.URL() != ""
	}, 5*time.Second, 10*time.Millisecond)

	// Send callback with wrong state — retry briefly for server readiness
	callbackURL := flow.CallbackURL() + "?code=auth-code-123&state=wrong-state"
	var resp *http.Response
	require.Eventually(t, func() bool {
		var getErr error
		resp, getErr = http.Get(callbackURL) //nolint:gosec // test URL
		return getErr == nil
	}, 2*time.Second, 10*time.Millisecond, "callback server never became ready")
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Flow should still be running (waiting for correct callback)
	// Cancel to end it
	cancel()
	err := <-errCh
	require.Error(t, err)
}

func TestLogin_NoClientID(t *testing.T) {
	flow := &LoginFlow{
		OAuthConfig: &oauth2.Config{
			ClientID: "",
		},
		TokenStore: NewKeyringTokenStore(t.TempDir()),
		Browser:    &FakeBrowserOpener{},
	}

	err := flow.Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client ID")
}
