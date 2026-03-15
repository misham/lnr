package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestPrintProjects_Empty(t *testing.T) {
	ios := NewTestIOStreams()
	err := PrintProjects(ios, nil)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "No projects found")
}

func TestPrintProjects_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	projects := []api.Project{
		{
			Name:       "Auth Rewrite",
			Status:     api.ProjectStatus{Name: "In Progress", Type: "started"},
			Lead:       &api.User{DisplayName: "Alice"},
			Progress:   0.45,
			StartDate:  "2026-03-01",
			TargetDate: "2026-06-01",
		},
		{
			Name:       "Dark Mode",
			Status:     api.ProjectStatus{Name: "Planned", Type: "planned"},
			Progress:   0,
			StartDate:  "2026-04-01",
			TargetDate: "2026-07-01",
		},
	}

	err := PrintProjects(ios, projects)
	require.NoError(t, err)

	out := outString(t, ios)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 2)
	assert.Contains(t, lines[0], "Auth Rewrite")
	assert.Contains(t, lines[0], "started")
	assert.Contains(t, lines[0], "Alice")
	assert.Contains(t, lines[0], "45%")
	assert.Contains(t, lines[1], "Dark Mode")
	assert.Contains(t, lines[1], "planned")
	assert.Contains(t, lines[1], "-")
}

func TestPrintProjects_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	projects := []api.Project{
		{
			Name:       "Auth Rewrite",
			Status:     api.ProjectStatus{Name: "In Progress", Type: "started"},
			Lead:       &api.User{DisplayName: "Alice"},
			Progress:   0.5,
			StartDate:  "2026-03-01",
			TargetDate: "2026-06-01",
		},
	}

	err := PrintProjects(ios, projects)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Auth Rewrite")
	assert.Contains(t, out, "started")
	assert.Contains(t, out, "50%")
}

func TestPrintProjectDetail_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	project := &api.Project{
		Name:          "Auth Rewrite",
		Status:        api.ProjectStatus{Name: "In Progress", Type: "started"},
		Lead:          &api.User{DisplayName: "Alice"},
		PriorityLabel: "High",
		Progress:      0.45,
		StartDate:     "2026-03-01",
		TargetDate:    "2026-06-01",
		URL:           "https://linear.app/team/project/auth-rewrite",
		Description:   "Rewrite the auth system",
		Milestones: []api.ProjectMilestone{
			{Name: "Design", TargetDate: "2026-03-15", Description: "Complete design docs"},
			{Name: "Implementation", TargetDate: "2026-05-01"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := PrintProjectDetail(ios, project)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Auth Rewrite")
	assert.Contains(t, out, "In Progress")
	assert.Contains(t, out, "started")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "High")
	assert.Contains(t, out, "45%")
	assert.Contains(t, out, "2026-03-01")
	assert.Contains(t, out, "2026-06-01")
	assert.Contains(t, out, "https://linear.app/team/project/auth-rewrite")
	assert.Contains(t, out, "Rewrite the auth system")
	assert.Contains(t, out, "Design")
	assert.Contains(t, out, "2026-03-15")
	assert.Contains(t, out, "Complete design docs")
	assert.Contains(t, out, "Implementation")
	assert.Contains(t, out, "2026-05-01")
}

func TestPrintProjectDetail_NoLead(t *testing.T) {
	ios := NewTestIOStreams()
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	project := &api.Project{
		Name:      "Solo Project",
		Status:    api.ProjectStatus{Name: "Planned", Type: "planned"},
		Lead:      nil,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := PrintProjectDetail(ios, project)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Solo Project")
	assert.Contains(t, out, "—")
}

func TestPrintProjectDetail_NoMilestones(t *testing.T) {
	ios := NewTestIOStreams()
	now := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	project := &api.Project{
		Name:       "Simple Project",
		Status:     api.ProjectStatus{Name: "Started", Type: "started"},
		Lead:       &api.User{DisplayName: "Bob"},
		Progress:   0.75,
		Milestones: nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := PrintProjectDetail(ios, project)
	require.NoError(t, err)

	out := outString(t, ios)
	assert.Contains(t, out, "Simple Project")
	assert.NotContains(t, out, "Milestones")
}
