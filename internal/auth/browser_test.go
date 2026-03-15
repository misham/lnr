package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeBrowserOpener(t *testing.T) {
	opener := &FakeBrowserOpener{}
	err := opener.Open("https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", opener.URL())
}
