package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestPrintInitiatives_Empty(t *testing.T) {
	ios := NewTestIOStreams()
	err := PrintInitiatives(ios, nil)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "No initiatives found")
}

func TestPrintInitiatives_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	initiatives := []api.Initiative{
		{
			Name:       "Platform Reliability",
			Status:     "Active",
			Owner:      &api.User{DisplayName: "Alice"},
			TargetDate: "2026-06-01",
		},
		{
			Name:       "Cost Reduction",
			Status:     "Planned",
			TargetDate: "2026-09-01",
		},
	}

	err := PrintInitiatives(ios, initiatives)
	require.NoError(t, err)

	out := outString(t, ios)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 2)
	assert.Contains(t, lines[0], "Platform Reliability")
	assert.Contains(t, lines[0], "Active")
	assert.Contains(t, lines[0], "Alice")
	assert.Contains(t, lines[0], "2026-06-01")
	assert.Contains(t, lines[1], "Cost Reduction")
	assert.Contains(t, lines[1], "Planned")
	assert.Contains(t, lines[1], "-")
}

func TestPrintInitiatives_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	initiatives := []api.Initiative{
		{
			Name:       "Platform Reliability",
			Status:     "Active",
			Owner:      &api.User{DisplayName: "Alice"},
			TargetDate: "2026-06-01",
		},
	}

	err := PrintInitiatives(ios, initiatives)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Platform Reliability")
	assert.Contains(t, out, "Active")
}

func TestPrintInitiativeDetail_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	initiative := &api.Initiative{
		Name:        "Platform Reliability",
		Status:      "Active",
		Health:      "onTrack",
		Owner:       &api.User{DisplayName: "Alice"},
		TargetDate:  "2026-06-01",
		URL:         "https://linear.app/team/initiative/platform-reliability",
		Description: "Improve platform reliability to 99.9%",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := PrintInitiativeDetail(ios, initiative)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Platform Reliability")
	assert.Contains(t, out, "Active")
	assert.Contains(t, out, "onTrack")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "2026-06-01")
	assert.Contains(t, out, "https://linear.app/team/initiative/platform-reliability")
	assert.Contains(t, out, "Improve platform reliability to 99.9%")
}

func TestPrintInitiativeDetail_NoOwner(t *testing.T) {
	ios := NewTestIOStreams()
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	initiative := &api.Initiative{
		Name:      "Solo Initiative",
		Status:    "Planned",
		Owner:     nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := PrintInitiativeDetail(ios, initiative)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Solo Initiative")
	assert.Contains(t, out, "—")
}
