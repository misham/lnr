package issue

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newCommentCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Manage issue comments",
	}
	cmd.AddCommand(newCommentListCmd(f))
	cmd.AddCommand(newCommentAddCmd(f))
	return cmd
}

func newCommentListCmd(f *cmdutil.Factory) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list <identifier>",
		Short: "List comments on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			comments, _, err := client.ListComments(cmd.Context(), args[0], limit, "")
			if err != nil {
				return fmt.Errorf("list comments: %w", err)
			}

			return ui.PrintComments(f.IO, comments)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of comments")

	return cmd
}

func newCommentAddCmd(f *cmdutil.Factory) *cobra.Command {
	var body string

	cmd := &cobra.Command{
		Use:   "add <identifier>",
		Short: "Add a comment to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if body == "" {
				return fmt.Errorf("body is required (use --body)")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.CreateComment(cmd.Context(), args[0], body)
			if err != nil {
				return fmt.Errorf("add comment: %w", err)
			}

			_, err = fmt.Fprintln(f.IO.Out, "Comment added")
			return err
		},
	}

	cmd.Flags().StringVar(&body, "body", "", "Comment body (required)")

	return cmd
}
