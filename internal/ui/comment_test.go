package ui

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestPrintComments_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(true)

	comments := []api.Comment{
		{
			Body:      "Looks good to me",
			CreatedAt: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			User:      &api.User{Name: "Alice"},
		},
		{
			Body:      "Needs changes",
			CreatedAt: time.Date(2026, 1, 16, 12, 0, 0, 0, time.UTC),
			User:      &api.User{Name: "Bob"},
		},
	}

	err := PrintComments(ios, comments)
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "Looks good to me")
	assert.Contains(t, out, "Bob")
	assert.Contains(t, out, "Needs changes")
}

func TestPrintComments_Empty(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(true)

	err := PrintComments(ios, []api.Comment{})
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "No comments")
}

func TestPrintComments_NilUser(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(true)

	comments := []api.Comment{
		{
			Body:      "Anonymous comment",
			CreatedAt: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			User:      nil,
		},
	}

	err := PrintComments(ios, comments)
	require.NoError(t, err)

	buf, _ := ios.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "unknown")
	assert.Contains(t, buf.String(), "Anonymous comment")
}
