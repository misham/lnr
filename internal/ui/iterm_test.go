package ui

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsITerm2_True(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "iTerm.app")
	assert.True(t, IsITerm2())
}

func TestIsITerm2_False(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "Apple_Terminal")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("LC_TERMINAL", "")
	assert.False(t, IsITerm2())
}

func TestIsITerm2_Unset(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("LC_TERMINAL", "")
	assert.False(t, IsITerm2())
}

func TestIsITerm2_Tmux(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "tmux")
	t.Setenv("ITERM_SESSION_ID", "w0t6p0:some-id")
	t.Setenv("LC_TERMINAL", "iTerm2")
	assert.True(t, IsITerm2())
}

func TestRenderInlineImage(t *testing.T) {
	data := []byte("fake-image-data")
	result := RenderInlineImage(data, "test.png")

	// Starts with escape sequence
	assert.True(t, strings.HasPrefix(result, "\033]1337;File="))
	// Contains required params
	assert.Contains(t, result, "inline=1")
	assert.Contains(t, result, "size=15")
	// Ends with BEL
	assert.True(t, strings.HasSuffix(result, "\a"))

	// Extract and verify base64 data
	parts := strings.SplitN(result, ":", 2)
	require.Len(t, parts, 2)
	b64Data := strings.TrimSuffix(parts[1], "\a")
	decoded, err := base64.StdEncoding.DecodeString(b64Data)
	require.NoError(t, err)
	assert.Equal(t, data, decoded)
}

func TestRenderInlineImage_Empty(t *testing.T) {
	assert.Equal(t, "", RenderInlineImage(nil, "test.png"))
	assert.Equal(t, "", RenderInlineImage([]byte{}, "test.png"))
}

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"screenshot.png", true},
		{"photo.jpg", true},
		{"photo.jpeg", true},
		{"animation.gif", true},
		{"SCREENSHOT.PNG", true}, // case-insensitive
		{"Photo.JPG", true},      // case-insensitive
		{"document.pdf", false},
		{"readme.txt", false},
		{"archive.zip", false},
		{"noextension", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.want, IsImageFile(tt.filename))
		})
	}
}
