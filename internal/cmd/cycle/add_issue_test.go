package cycle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestAddIssueCmd_SetsCycleOnIssue(t *testing.T) {
	fc := &fakeClient{
		cycle: &api.Cycle{
			ID:        "cycle-1",
			Number:    5,
			Name:      "Sprint 5",
			IsActive:  true,
			StartsAt:  time.Now(),
			EndsAt:    time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		issue: &api.Issue{
			ID:         "issue-1",
			Identifier: "ENG-42",
			Title:      "Fix bug",
		},
	}

	f := newTestFactory(t, fc)
	cmd := newAddIssueCmd(f)
	cmd.SetArgs([]string{"5", "issue-1"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "issue-1", fc.updateID)
	require.NotNil(t, fc.updateInput.CycleID)
	assert.Equal(t, "cycle-1", *fc.updateInput.CycleID)
}

func TestAddIssueCmd_RequiresTwoArgs(t *testing.T) {
	fc := &fakeClient{}

	f := newTestFactory(t, fc)
	cmd := newAddIssueCmd(f)
	cmd.SetArgs([]string{"5"})

	err := cmd.Execute()
	require.Error(t, err)
}
