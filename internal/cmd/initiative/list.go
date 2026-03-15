package initiative

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

var validInitiativeStatuses = []string{"Planned", "Active", "Completed"}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var statusFilter string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List initiatives",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if statusFilter != "" {
				valid := false
				for _, s := range validInitiativeStatuses {
					if strings.EqualFold(s, statusFilter) {
						statusFilter = s
						valid = true
						break
					}
				}
				if !valid {
					return fmt.Errorf("invalid status %q: must be one of %v (case-insensitive)", statusFilter, validInitiativeStatuses)
				}
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			var all []api.Initiative
			cursor := ""
			for {
				result, err := client.ListInitiatives(cmd.Context(), statusFilter, 50, cursor)
				if err != nil {
					return fmt.Errorf("list initiatives: %w", err)
				}
				all = append(all, result.Initiatives...)
				if !result.PageInfo.HasNextPage {
					break
				}
				cursor = result.PageInfo.EndCursor
			}

			return ui.PrintInitiatives(f.IO, all)
		},
	}

	cmd.Flags().StringVar(&statusFilter, "status", "", "Filter by status (Planned, Active, Completed; case-insensitive)")

	return cmd
}
