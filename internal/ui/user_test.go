package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestPrintUser_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	user := &api.User{
		ID:          "user-1",
		Name:        "John Doe",
		DisplayName: "johnd",
		Email:       "john@example.com",
		Active:      true,
	}

	err := PrintUser(ios, user)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "user-1")
	assert.Contains(t, out, "John Doe")
	assert.Contains(t, out, "johnd")
	assert.Contains(t, out, "john@example.com")
	assert.Contains(t, out, "true")
}

func TestPrintUser_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	user := &api.User{
		ID:          "user-1",
		Name:        "John Doe",
		DisplayName: "johnd",
		Email:       "john@example.com",
		Active:      true,
	}

	err := PrintUser(ios, user)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "John Doe")
	assert.Contains(t, out, "john@example.com")
}
