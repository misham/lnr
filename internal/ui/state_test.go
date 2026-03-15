package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestPrintStates_Empty(t *testing.T) {
	ios := NewTestIOStreams()
	err := PrintStates(ios, nil)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "No workflow states found")
}

func TestPrintStates_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	states := []api.WorkflowState{
		{ID: "s1", Name: "Backlog", Type: "backlog", Color: "#bec2c8", Position: 0},
		{ID: "s2", Name: "In Progress", Type: "started", Color: "#f2c94c", Position: 1},
	}

	err := PrintStates(ios, states)
	require.NoError(t, err)

	out := outString(t, ios)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 2)
	assert.Contains(t, lines[0], "Backlog")
	assert.Contains(t, lines[0], "backlog")
	assert.Contains(t, lines[0], "#bec2c8")
	assert.True(t, strings.HasSuffix(lines[0], "\t0"), "expected line to end with tab-separated 0")
	assert.Contains(t, lines[1], "In Progress")
	assert.Contains(t, lines[1], "started")
	assert.Contains(t, lines[1], "#f2c94c")
	assert.True(t, strings.HasSuffix(lines[1], "\t1"), "expected line to end with tab-separated 1")
}

func TestPrintStates_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	states := []api.WorkflowState{
		{ID: "s1", Name: "In Progress", Type: "started", Color: "#f2c94c", Position: 1},
	}

	err := PrintStates(ios, states)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "NAME")
	assert.Contains(t, out, "CATEGORY")
	assert.Contains(t, out, "COLOR")
	assert.Contains(t, out, "POSITION")
	assert.Contains(t, out, "In Progress")
	assert.Contains(t, out, "started")
	assert.Contains(t, out, "#f2c94c")
}
