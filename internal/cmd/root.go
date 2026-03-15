package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	authcmd "github.com/misham/linear-cli/internal/cmd/auth"
	cyclecmd "github.com/misham/linear-cli/internal/cmd/cycle"
	initiativecmd "github.com/misham/linear-cli/internal/cmd/initiative"
	issuecmd "github.com/misham/linear-cli/internal/cmd/issue"
	projectcmd "github.com/misham/linear-cli/internal/cmd/project"
	statecmd "github.com/misham/linear-cli/internal/cmd/state"
	teamcmd "github.com/misham/linear-cli/internal/cmd/team"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/version"
)

func newRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "lnr",
		Short:         "Linear CLI — manage issues, cycles, and projects",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			plain, _ := cmd.Flags().GetBool("plain")
			if plain {
				f.IO.SetPlain(true)
			}
			teamFlag, _ := cmd.Flags().GetString("team")
			if teamFlag != "" {
				f.TeamKey = teamFlag
			}
			return nil
		},
	}

	cmd.PersistentFlags().Bool("plain", false, "Disable styled output")
	cmd.PersistentFlags().String("team", "", "Team key or ID (overrides default)")

	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		_ = cmd.Usage()
		return err
	})

	cmd.AddCommand(authcmd.NewAuthCmd(f))
	cmd.AddCommand(cyclecmd.NewCycleCmd(f))
	cmd.AddCommand(issuecmd.NewIssueCmd(f))
	cmd.AddCommand(initiativecmd.NewInitiativeCmd(f))
	cmd.AddCommand(projectcmd.NewProjectCmd(f))
	cmd.AddCommand(statecmd.NewStateCmd(f))
	cmd.AddCommand(teamcmd.NewTeamCmd(f))
	cmd.AddCommand(newMeCmd(f))
	cmd.AddCommand(newTuiCmd(f))
	cmd.AddCommand(newCompletionCmd(f))

	return cmd
}

// Execute runs the root command.
func Execute() error {
	f := cmdutil.NewFactory()
	rootCmd := newRootCmd(f)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	return nil
}
