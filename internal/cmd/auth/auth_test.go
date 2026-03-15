package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestNewAuthCmd(t *testing.T) {
	f := &cmdutil.Factory{IO: ui.NewTestIOStreams()}
	cmd := NewAuthCmd(f)
	assert.Equal(t, "auth", cmd.Use)
	assert.True(t, cmd.HasSubCommands())
}
