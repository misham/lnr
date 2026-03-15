package team

import (
	"testing"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestNewTeamCmd_HasSubcommands(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	cmd := NewTeamCmd(f)

	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}

	if !names["list"] {
		t.Error("expected list subcommand")
	}
}
