package team

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
)

func newSetCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "set [team-key]",
		Short: "Set default team",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			teams, err := client.ListTeams(cmd.Context())
			if err != nil {
				return fmt.Errorf("list teams: %w", err)
			}

			var selected api.Team

			if len(args) == 1 {
				key := strings.ToUpper(args[0])
				found := false
				for _, t := range teams {
					if strings.EqualFold(t.Key, key) {
						selected = t
						found = true
						break
					}
				}
				if !found {
					keys := make([]string, len(teams))
					for i, t := range teams {
						keys[i] = t.Key
					}
					return fmt.Errorf("team %q not found, valid keys: %s", args[0], strings.Join(keys, ", "))
				}
			} else {
				if f.IO.IsPlain() {
					return fmt.Errorf("team key required in non-interactive mode (interactive picker requires a TTY)")
				}

				options := make([]huh.Option[string], len(teams))
				for i, t := range teams {
					options[i] = huh.NewOption(fmt.Sprintf("%s (%s)", t.Key, t.Name), t.ID)
				}

				var selectedID string
				err := huh.NewSelect[string]().
					Title("Select a team").
					Options(options...).
					Value(&selectedID).
					Run()
				if err != nil {
					return fmt.Errorf("select team: %w", err)
				}

				for _, t := range teams {
					if t.ID == selectedID {
						selected = t
						break
					}
				}
			}

			cfg, err := f.Config()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			cfg.SetTeamID(selected.ID)
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			_, _ = fmt.Fprintf(f.IO.Out, "Default team set to %s (%s)\n", selected.Key, selected.Name)
			return nil
		},
	}
}
