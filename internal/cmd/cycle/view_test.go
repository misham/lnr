package cycle

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestViewCmd_ByNumber(t *testing.T) {
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		cycle: &api.Cycle{
			Number:    5,
			Name:      "Sprint 5",
			IsActive:  true,
			StartsAt:  start,
			EndsAt:    end,
			Progress:  0.75,
			Team:      api.Team{Name: "Engineering"},
			CreatedAt: start,
			UpdatedAt: end,
		},
	}

	f := newTestFactory(t, fc)
	cmd := newViewCmd(f)
	cmd.SetArgs([]string{"5"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Sprint 5")
	assert.Contains(t, buf.String(), "#5")
}

func TestViewCmd_ByUUID(t *testing.T) {
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		cycle: &api.Cycle{
			ID:        "cycle-uuid-123",
			Number:    5,
			Name:      "Sprint 5",
			StartsAt:  start,
			EndsAt:    end,
			Team:      api.Team{Name: "Engineering"},
			CreatedAt: start,
			UpdatedAt: end,
		},
	}

	f := newTestFactory(t, fc)
	cmd := newViewCmd(f)
	cmd.SetArgs([]string{"cycle-uuid-123"})

	err := cmd.Execute()
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	assert.Contains(t, buf.String(), "Sprint 5")
}

func TestViewCmd_NoArg(t *testing.T) {
	fc := &fakeClient{}

	f := newTestFactory(t, fc)
	cmd := newViewCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
}
