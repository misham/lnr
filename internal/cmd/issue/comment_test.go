package issue

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
	"github.com/misham/linear-cli/internal/cmdutil"
	"github.com/misham/linear-cli/internal/ui"
)

func TestCommentListCmd_ShowsComments(t *testing.T) {
	fc := &fakeClient{
		comments: []api.Comment{
			{
				ID:        "c1",
				Body:      "First comment",
				CreatedAt: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
				User:      &api.User{Name: "Alice"},
			},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newCommentCmd(f)
	cmd.SetArgs([]string{"list", "ENG-1"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Alice")
	assert.Contains(t, buf.String(), "First comment")
}

func TestCommentAddCmd_CreatesComment(t *testing.T) {
	fc := &fakeClient{
		comment: &api.Comment{
			ID:        "c2",
			Body:      "New comment",
			CreatedAt: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			User:      &api.User{Name: "Bob"},
		},
	}

	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return fc, nil
		},
	}

	cmd := newCommentCmd(f)
	cmd.SetArgs([]string{"add", "ENG-1", "--body", "New comment"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Comment added")
}

func TestCommentAddCmd_NoBody(t *testing.T) {
	ios := ui.NewTestIOStreams()
	f := &cmdutil.Factory{
		IO: ios,
		APIClient: func() (api.Client, error) {
			return &fakeClient{}, nil
		},
	}

	cmd := newCommentCmd(f)
	cmd.SetArgs([]string{"add", "ENG-1"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body is required")
}
