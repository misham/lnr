package issue

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/config"
	"github.com/misham/linear-cli/internal/ui"
)

func TestFileUploadCmd(t *testing.T) {
	// Create a temp file to upload.
	dir := t.TempDir()
	filePath := filepath.Join(dir, "report.pdf")
	require.NoError(t, os.WriteFile(filePath, []byte("fake-pdf-content"), 0o600))

	fc := &fakeClient{
		issue: &api.Issue{
			ID:          "issue-uuid-123",
			Identifier:  "ENG-123",
			Description: "Existing description",
		},
		uploadResult: &api.UploadResult{
			UploadURL: "https://s3.example.com/presigned",
			AssetURL:  "https://uploads.linear.app/org/report.pdf",
			Headers:   []api.UploadHeader{{Key: "x-amz-acl", Value: "public-read"}},
		},
	}

	ios := ui.NewTestIOStreams()
	store := config.NewViperStore(t.TempDir())
	require.NoError(t, store.Load())

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newFileCmd(f)
	cmd.SetArgs([]string{"upload", "ENG-123", filePath})

	err := cmd.Execute()
	require.NoError(t, err)

	// Verify upload was executed.
	assert.Equal(t, "https://s3.example.com/presigned", fc.uploadedURL)

	// Verify UpdateIssue was called to append the file link to the description.
	assert.Equal(t, "issue-uuid-123", fc.updateID)
	require.NotNil(t, fc.updateInput.Description)
	assert.Contains(t, *fc.updateInput.Description, "Existing description")
	assert.Contains(t, *fc.updateInput.Description, "[report.pdf](https://uploads.linear.app/org/report.pdf)")

	// Verify confirmation output.
	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Uploaded report.pdf")
}

func TestFileUploadCmd_FileNotFound(t *testing.T) {
	fc := &fakeClient{}

	ios := ui.NewTestIOStreams()
	store := config.NewViperStore(t.TempDir())
	require.NoError(t, store.Load())

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newFileCmd(f)
	cmd.SetArgs([]string{"upload", "ENG-123", "/nonexistent/file.pdf"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "open")
}

func TestFileUploadCmd_NoArgs(t *testing.T) {
	fc := &fakeClient{}

	ios := ui.NewTestIOStreams()
	store := config.NewViperStore(t.TempDir())
	require.NoError(t, store.Load())

	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
		Config: func() (config.Store, error) {
			return store, nil
		},
	}

	cmd := newFileCmd(f)
	cmd.SetArgs([]string{"upload"})

	err := cmd.Execute()
	require.Error(t, err)
}
