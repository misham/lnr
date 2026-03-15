package ui

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsITerm2 returns true if running in iTerm2, including inside tmux.
func IsITerm2() bool {
	if os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		return true
	}
	// iTerm2 sets these even inside tmux/screen.
	return os.Getenv("ITERM_SESSION_ID") != "" || os.Getenv("LC_TERMINAL") == "iTerm2"
}

// RenderInlineImage returns the iTerm2 escape sequence for inline image display.
// Returns empty string if data is empty.
func RenderInlineImage(data []byte, filename string) string {
	if len(data) == 0 {
		return ""
	}
	b64Data := base64.StdEncoding.EncodeToString(data)
	b64Name := base64.StdEncoding.EncodeToString([]byte(filename))
	return fmt.Sprintf("\033]1337;File=name=%s;size=%d;inline=1:%s\a",
		b64Name, len(data), b64Data)
}

// IsImageFile returns true if the filename has an image extension (png, jpg, jpeg, gif).
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif":
		return true
	default:
		return false
	}
}
