package issue

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
)

func newLabelCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage issue labels",
	}
	cmd.AddCommand(newLabelListCmd(f))
	cmd.AddCommand(newLabelAddCmd(f))
	cmd.AddCommand(newLabelRemoveCmd(f))
	return cmd
}

func newLabelListCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list [identifier]",
		Short: "List available labels, or labels on a specific issue",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if len(args) == 1 {
				issue, err := client.GetIssue(cmd.Context(), args[0])
				if err != nil {
					return fmt.Errorf("get issue: %w", err)
				}
				if len(issue.Labels) == 0 {
					_, err := fmt.Fprintln(f.IO.Out, "No labels")
					return err
				}
				for _, l := range issue.Labels {
					if _, err := fmt.Fprintln(f.IO.Out, l.Name); err != nil {
						return err
					}
				}
				return nil
			}

			labels, err := client.ListLabels(cmd.Context(), "")
			if err != nil {
				return fmt.Errorf("list labels: %w", err)
			}

			for _, l := range labels {
				if _, err := fmt.Fprintln(f.IO.Out, l.Name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newLabelAddCmd(f *cmdutil.Factory) *cobra.Command {
	var labelName string

	cmd := &cobra.Command{
		Use:   "add <identifier>",
		Short: "Add a label to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			labelID, err := resolveLabel(cmd, client, args[0], labelName)
			if err != nil {
				return err
			}

			if err := client.AddIssueLabel(cmd.Context(), args[0], labelID); err != nil {
				return fmt.Errorf("add label: %w", err)
			}

			_, err = fmt.Fprintf(f.IO.Out, "Added label %q to %s\n", labelName, args[0])
			return err
		},
	}

	cmd.Flags().StringVar(&labelName, "label", "", "Label name (required)")

	return cmd
}

func newLabelRemoveCmd(f *cmdutil.Factory) *cobra.Command {
	var labelName string

	cmd := &cobra.Command{
		Use:   "remove <identifier>",
		Short: "Remove a label from an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			labelID, err := resolveLabel(cmd, client, args[0], labelName)
			if err != nil {
				return err
			}

			if err := client.RemoveIssueLabel(cmd.Context(), args[0], labelID); err != nil {
				return fmt.Errorf("remove label: %w", err)
			}

			_, err = fmt.Fprintf(f.IO.Out, "Removed label %q from %s\n", labelName, args[0])
			return err
		},
	}

	cmd.Flags().StringVar(&labelName, "label", "", "Label name (required)")

	return cmd
}

func resolveLabel(cmd *cobra.Command, client api.Client, issueID, labelName string) (string, error) {
	issue, err := client.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return "", fmt.Errorf("get issue: %w", err)
	}

	labels, err := client.ListLabels(cmd.Context(), issue.Team.ID)
	if err != nil {
		return "", fmt.Errorf("list labels: %w", err)
	}

	for _, l := range labels {
		if strings.EqualFold(l.Name, labelName) {
			return l.ID, nil
		}
	}

	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return "", fmt.Errorf("label %q not found; available: %s", labelName, strings.Join(names, ", "))
}
