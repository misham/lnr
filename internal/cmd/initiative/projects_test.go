package initiative

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestProjectsCmd_ShowsProjects(t *testing.T) {
	fc := &fakeClient{
		initiative: &api.Initiative{ID: "ini-1", Name: "Q2 Goals"},
		projects: &api.ProjectListResult{
			Projects: []api.Project{
				{Name: "Auth Rewrite", Status: api.ProjectStatus{Type: "started"}, Progress: 0.5},
				{Name: "Dark Mode", Status: api.ProjectStatus{Type: "planned"}, Progress: 0},
			},
		},
	}
	f := newTestFactory(fc)

	cmd := NewInitiativeCmd(f)
	cmd.SetArgs([]string{"projects", "some-id"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "Auth Rewrite")
	assert.Contains(t, out, "Dark Mode")
}

func TestProjectsCmd_NoProjects(t *testing.T) {
	fc := &fakeClient{
		initiative: &api.Initiative{ID: "ini-1", Name: "Q2 Goals"},
		projects:   &api.ProjectListResult{Projects: nil},
	}
	f := newTestFactory(fc)

	cmd := NewInitiativeCmd(f)
	cmd.SetArgs([]string{"projects", "some-id"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf, _ := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "No projects found")
}
