package project

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestViewCmd_ShowsProject(t *testing.T) {
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		project: &api.Project{
			Name:          "Auth Rewrite",
			Status:        api.ProjectStatus{Name: "In Progress", Type: "started"},
			Lead:          &api.User{DisplayName: "Alice"},
			PriorityLabel: "High",
			Progress:      0.45,
			URL:           "https://linear.app/team/project/auth-rewrite",
			Description:   "Rewrite the auth system",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"view", "auth-rewrite"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Auth Rewrite")
	assert.Contains(t, out, "In Progress")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "Rewrite the auth system")
}

func TestViewCmd_WithMilestones(t *testing.T) {
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	fc := &fakeClient{
		project: &api.Project{
			Name:     "Auth Rewrite",
			Status:   api.ProjectStatus{Name: "Started", Type: "started"},
			Progress: 0.5,
			Milestones: []api.ProjectMilestone{
				{Name: "Design Phase", TargetDate: "2026-03-15"},
				{Name: "Implementation", TargetDate: "2026-05-01"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"view", "auth-rewrite"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Design Phase")
	assert.Contains(t, out, "2026-03-15")
	assert.Contains(t, out, "Implementation")
	assert.Contains(t, out, "2026-05-01")
}
