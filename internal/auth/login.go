package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// defaultClientID is set at build time via ldflags.
var defaultClientID string //nolint:gochecknoglobals // set by ldflags

// LoginFlow manages the OAuth 2.0 PKCE login flow.
type LoginFlow struct {
	OAuthConfig *oauth2.Config
	TokenStore  TokenStore
	Browser     BrowserOpener
	ListenAddr  string

	mu          sync.Mutex
	callbackURL string
}

// CallbackURL returns the callback URL once the server is listening.
func (f *LoginFlow) CallbackURL() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.callbackURL
}

// DefaultClientID returns the build-time client ID.
func DefaultClientID() string {
	return defaultClientID
}

// Run executes the OAuth 2.0 PKCE login flow.
func (f *LoginFlow) Run(ctx context.Context) error {
	if f.OAuthConfig.ClientID == "" {
		return errors.New("client ID not set: set LNR_CLIENT_ID or build with -ldflags")
	}

	verifier, err := generateVerifier()
	if err != nil {
		return fmt.Errorf("generate PKCE verifier: %w", err)
	}

	state, err := generateState()
	if err != nil {
		return fmt.Errorf("generate state: %w", err)
	}

	listenAddr := f.ListenAddr
	if listenAddr == "" {
		listenAddr = "localhost:8585"
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	// Store the actual callback URL (with dynamic port if :0)
	addr := listener.Addr().String()
	f.mu.Lock()
	f.callbackURL = "http://" + addr + "/callback"
	f.mu.Unlock()

	// Build auth URL with PKCE and Linear's comma-separated scopes
	authURL := f.OAuthConfig.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("scope", "read,write"),
	)

	if err := f.Browser.Open(authURL); err != nil {
		_ = listener.Close()
		return fmt.Errorf("open browser: %w", err)
	}

	type result struct {
		token *oauth2.Token
		err   error
	}
	resultCh := make(chan result, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state mismatch", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}

		token, exchangeErr := f.OAuthConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if exchangeErr != nil {
			errMsg := fmt.Sprintf("token exchange failed: %v", exchangeErr)
			http.Error(w, errMsg, http.StatusInternalServerError)
			resultCh <- result{err: fmt.Errorf("exchange token: %w", exchangeErr)}
			return
		}

		_, _ = fmt.Fprint(w, "Login successful! You can close this window.")
		resultCh <- result{token: token}
	})

	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			resultCh <- result{err: fmt.Errorf("serve: %w", err)}
		}
	}()

	defer func() { _ = srv.Close() }()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-resultCh:
		if res.err != nil {
			return res.err
		}
		if err := f.TokenStore.SaveToken(res.token); err != nil {
			return fmt.Errorf("save token: %w", err)
		}
		return nil
	}
}

func generateVerifier() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("random: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func generateState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("random: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
