package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/stretchr/testify/assert"
)

func TestDefaultKeyMap_AllBindingsEnabled(t *testing.T) {
	km := DefaultKeyMap()

	bindings := []struct {
		name    string
		binding key.Binding
	}{
		{"Up", km.Up},
		{"Down", km.Down},
		{"Enter", km.Enter},
		{"Back", km.Back},
		{"Tab", km.Tab},
		{"ShiftTab", km.ShiftTab},
		{"Tab1", km.Tab1},
		{"Tab2", km.Tab2},
		{"Tab3", km.Tab3},
		{"Tab4", km.Tab4},
		{"Search", km.Search},
		{"NextPage", km.NextPage},
		{"PrevPage", km.PrevPage},
		{"Top", km.Top},
		{"Bottom", km.Bottom},
		{"Retry", km.Retry},
		{"Help", km.Help},
		{"Quit", km.Quit},
	}

	for _, b := range bindings {
		assert.True(t, b.binding.Enabled(), "%s should be enabled", b.name)
		assert.NotEmpty(t, b.binding.Help().Key, "%s should have help key text", b.name)
		assert.NotEmpty(t, b.binding.Help().Desc, "%s should have help description", b.name)
	}
}

func TestDefaultKeyMap_NoDuplicateKeys(t *testing.T) {
	km := DefaultKeyMap()

	// Collect all key strings, excluding modifiers that are naturally shared.
	bindings := []key.Binding{
		km.Up, km.Down, km.Enter, km.Back,
		km.Tab1, km.Tab2, km.Tab3, km.Tab4,
		km.Search, km.NextPage, km.PrevPage,
		km.Top, km.Bottom, km.Retry, km.Help, km.Quit,
	}

	seen := make(map[string]string)
	for _, b := range bindings {
		for _, k := range b.Keys() {
			if prev, exists := seen[k]; exists {
				t.Errorf("key %q is used by both %q and %q", k, prev, b.Help().Desc)
			}
			seen[k] = b.Help().Desc
		}
	}
}
