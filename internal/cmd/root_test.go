package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestRootCommand_Version(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	rootCmd := newRootCmd(f)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if got == "" {
		t.Fatal("expected version output, got empty string")
	}
}

func TestRootCommand_FlagError(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	rootCmd := newRootCmd(f)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--nonexistent"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestRootCommand_HasSubcommands(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	rootCmd := newRootCmd(f)

	names := make(map[string]bool)
	for _, sub := range rootCmd.Commands() {
		names[sub.Name()] = true
	}

	if !names["auth"] {
		t.Error("expected auth subcommand")
	}
	if !names["team"] {
		t.Error("expected team subcommand")
	}
}

func TestExecute_ReturnsNilOnSuccess(t *testing.T) {
	// Execute creates a real factory but --version doesn't touch auth/config
	// so this just verifies the wiring works end-to-end
	origArgs := os.Args
	os.Args = []string{"lnr", "--version"}
	defer func() { os.Args = origArgs }()

	err := Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootCommand_PlainFlag(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{IO: ios}

	rootCmd := newRootCmd(f)
	// Add a no-op subcommand so PersistentPreRunE fires (Cobra skips it for --version/--help)
	noop := &cobra.Command{
		Use:  "noop",
		RunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
	rootCmd.AddCommand(noop)
	rootCmd.SetArgs([]string{"--plain", "noop"})
	rootCmd.SetOut(new(bytes.Buffer))

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ios.IsPlain() {
		t.Error("expected IOStreams to be plain after --plain flag")
	}
}
