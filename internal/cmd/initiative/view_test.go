package initiative

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestViewCmd_ShowsInitiative(t *testing.T) {
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		initiative: &api.Initiative{
			Name:        "Platform Reliability",
			Status:      "Active",
			Health:      "onTrack",
			Owner:       &api.User{DisplayName: "Alice"},
			TargetDate:  "2026-06-01",
			URL:         "https://linear.app/team/initiative/platform-reliability",
			Description: "Improve platform reliability",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	f := newTestFactory(fc)

	cmd := NewInitiativeCmd(f)
	cmd.SetArgs([]string{"view", "some-id"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Platform Reliability")
	assert.Contains(t, out, "Active")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "Improve platform reliability")
}
