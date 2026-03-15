package auth

import (
	"os/exec"
	"sync"
)

// BrowserOpener abstracts browser launching for testability.
type BrowserOpener interface {
	Open(url string) error
}

// SystemBrowserOpener opens URLs in the system browser (macOS).
type SystemBrowserOpener struct{}

func (o *SystemBrowserOpener) Open(url string) error {
	return exec.Command("open", url).Start() //nolint:gosec // intentional: opens URL in browser
}

// FakeBrowserOpener records the URL for testing.
type FakeBrowserOpener struct {
	mu      sync.Mutex
	LastURL string
}

func (o *FakeBrowserOpener) Open(url string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.LastURL = url
	return nil
}

// URL returns the last opened URL (thread-safe).
func (o *FakeBrowserOpener) URL() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.LastURL
}
