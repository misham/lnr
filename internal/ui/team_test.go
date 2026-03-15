package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func outString(t *testing.T, ios *IOStreams) string {
	t.Helper()
	buf, ok := ios.Out.(*bytes.Buffer)
	require.True(t, ok)
	return buf.String()
}

func TestPrintTeams_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	teams := []api.Team{
		{Key: "ENG", Name: "Engineering", Description: "Core engineering team"},
		{Key: "DES", Name: "Design", Description: "Product design"},
	}

	err := PrintTeams(ios, teams)
	require.NoError(t, err)

	out := outString(t, ios)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 2)
	assert.Contains(t, lines[0], "ENG")
	assert.Contains(t, lines[0], "Engineering")
	assert.Contains(t, lines[0], "Core engineering team")
	assert.Contains(t, lines[1], "DES")
	assert.Contains(t, lines[1], "Design")
}

func TestPrintTeams_Empty(t *testing.T) {
	ios := NewTestIOStreams()
	err := PrintTeams(ios, nil)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "No teams found")
}

func TestPrintTeams_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	teams := []api.Team{
		{Key: "ENG", Name: "Engineering", Description: "Core engineering team"},
	}

	err := PrintTeams(ios, teams)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "ENG")
	assert.Contains(t, out, "Engineering")
}
