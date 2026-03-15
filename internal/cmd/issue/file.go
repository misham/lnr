package issue

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func newFileCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Manage issue files",
	}
	cmd.AddCommand(newFileUploadCmd(f))
	return cmd
}

func newFileUploadCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "upload <identifier> <file-path>",
		Short: "Upload a file and append it to the issue description",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			identifier := args[0]
			filePath := args[1]

			filePath = filepath.Clean(filePath)
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func() { _ = file.Close() }()

			info, err := file.Stat()
			if err != nil {
				return err
			}

			filename := filepath.Base(filePath)
			contentType := mime.TypeByExtension(filepath.Ext(filename))
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			ctx := cmd.Context()

			issue, err := client.GetIssue(ctx, identifier)
			if err != nil {
				return fmt.Errorf("get issue: %w", err)
			}

			result, err := client.FileUpload(ctx, contentType, filename, info.Size())
			if err != nil {
				return fmt.Errorf("file upload: %w", err)
			}

			if _, err := file.Seek(0, 0); err != nil {
				return err
			}
			headers := append(result.Headers, api.UploadHeader{Key: "Content-Type", Value: contentType})
			if err := client.UploadToURL(ctx, result.UploadURL, headers, file); err != nil {
				return fmt.Errorf("upload to URL: %w", err)
			}

			prefix := ""
			if ui.IsImageFile(filename) {
				prefix = "!"
			}
			link := fmt.Sprintf("%s[%s](%s)", prefix, filename, result.AssetURL)
			desc := issue.Description
			if desc != "" {
				desc += "\n\n"
			}
			desc += link
			if _, err := client.UpdateIssue(ctx, issue.ID, api.IssueUpdateInput{
				Description: &desc,
			}); err != nil {
				return fmt.Errorf("update issue: %w", err)
			}

			_, err = fmt.Fprintf(f.IO.Out, "Uploaded %s to %s\n", filename, identifier)
			return err
		},
	}
}
