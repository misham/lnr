package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestCompletionCmd_Bash(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}
	buf, ok := ios.Out.(*bytes.Buffer)
	require.True(t, ok)

	root := newRootCmd(f)
	root.SetArgs([]string{"completion", "bash"})
	root.SetOut(buf)

	err := root.Execute()
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "bash completion")
	assert.Contains(t, buf.String(), "lnr")
}

func TestCompletionCmd_Zsh(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}
	buf, ok := ios.Out.(*bytes.Buffer)
	require.True(t, ok)

	root := newRootCmd(f)
	root.SetArgs([]string{"completion", "zsh"})
	root.SetOut(buf)

	err := root.Execute()
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
}

func TestCompletionCmd_NoArgs(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	root := newRootCmd(f)
	root.SetArgs([]string{"completion"})

	err := root.Execute()
	require.NoError(t, err)
}

func TestCompletionCmd_Exists(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	cmd := newCompletionCmd(f)
	assert.Equal(t, "completion", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.Len(t, cmd.Commands(), 2, "should have bash and zsh subcommands")
}
