package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestPrintCycles_Empty(t *testing.T) {
	ios := NewTestIOStreams()
	err := PrintCycles(ios, nil)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "No cycles found")
}

func TestPrintCycles_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	cycles := []api.Cycle{
		{Number: 5, Name: "Sprint 5", IsActive: true, StartsAt: start, EndsAt: end, Progress: 0.75},
		{Number: 6, Name: "Sprint 6", IsNext: true, StartsAt: end, EndsAt: end.AddDate(0, 0, 14), Progress: 0},
	}

	err := PrintCycles(ios, cycles)
	require.NoError(t, err)

	out := outString(t, ios)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 2)
	assert.Contains(t, lines[0], "5")
	assert.Contains(t, lines[0], "Sprint 5")
	assert.Contains(t, lines[0], "Active")
	assert.Contains(t, lines[0], "75%")
	assert.Contains(t, lines[1], "6")
	assert.Contains(t, lines[1], "Sprint 6")
	assert.Contains(t, lines[1], "Next")
}

func TestPrintCycles_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	cycles := []api.Cycle{
		{Number: 5, Name: "Sprint 5", IsActive: true, StartsAt: start, EndsAt: end, Progress: 0.5},
	}

	err := PrintCycles(ios, cycles)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Sprint 5")
	assert.Contains(t, out, "Active")
	assert.Contains(t, out, "50%")
}

func TestPrintCycleDetail_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	completed := time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC)
	cycle := &api.Cycle{
		Number:      5,
		Name:        "Sprint 5",
		Description: "Focus on auth improvements",
		IsActive:    false,
		IsPast:      true,
		StartsAt:    start,
		EndsAt:      end,
		CompletedAt: &completed,
		Progress:    1.0,
		Team:        api.Team{Name: "Engineering", Key: "ENG"},
		CreatedAt:   start,
		UpdatedAt:   end,
	}

	err := PrintCycleDetail(ios, cycle)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Sprint 5")
	assert.Contains(t, out, "#5")
	assert.Contains(t, out, "Past")
	assert.Contains(t, out, "Engineering")
	assert.Contains(t, out, "100%")
	assert.Contains(t, out, "Focus on auth improvements")
	assert.Contains(t, out, "Mar 01")
	assert.Contains(t, out, "Mar 14")
}

func TestPrintCycleDetail_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	cycle := &api.Cycle{
		Number:    5,
		Name:      "Sprint 5",
		IsActive:  true,
		StartsAt:  start,
		EndsAt:    end,
		Progress:  0.5,
		Team:      api.Team{Name: "Engineering"},
		CreatedAt: start,
		UpdatedAt: end,
	}

	err := PrintCycleDetail(ios, cycle)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Sprint 5")
	assert.Contains(t, out, "Active")
	assert.Contains(t, out, "Engineering")
	assert.Contains(t, out, "50%")
}

func TestCycleStatus(t *testing.T) {
	tests := []struct {
		cycle    api.Cycle
		expected string
	}{
		{api.Cycle{IsActive: true}, "Active"},
		{api.Cycle{IsNext: true}, "Next"},
		{api.Cycle{IsPast: true}, "Past"},
		{api.Cycle{}, "Future"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, cycleStatus(&tt.cycle))
	}
}
