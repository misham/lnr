package cmd

import (
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
)

func newCompletionCmd(_ *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for bash or zsh.

To load completions:

Bash:
  source <(lnr completion bash)

  # To install permanently:
  lnr completion bash > /etc/bash_completion.d/lnr

Zsh:
  source <(lnr completion zsh)

  # To install permanently (ensure $fpath contains the directory):
  lnr completion zsh > "${fpath[1]}/_lnr"
  compinit`,
	}

	cmd.AddCommand(newCompletionBashCmd())
	cmd.AddCommand(newCompletionZshCmd())

	return cmd
}

func newCompletionBashCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: "Generate bash completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Root().GenBashCompletionV2(cmd.OutOrStdout(), true)
		},
	}
}

func newCompletionZshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Generate zsh completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		},
	}
}
