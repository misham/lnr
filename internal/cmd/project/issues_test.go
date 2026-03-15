package project

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/misham/linear-cli/internal/api"
)

func TestIssuesCmd_ShowsIssues(t *testing.T) {
	fc := &fakeClient{
		project: &api.Project{ID: "proj-1", Name: "Auth Rewrite"},
		issues: &api.IssueListResult{
			Issues: []api.Issue{
				{Identifier: "ENG-123", Title: "Fix login bug"},
				{Identifier: "ENG-456", Title: "Add dark mode"},
			},
		},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"issues", "auth-rewrite"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "ENG-123")
	assert.Contains(t, out, "ENG-456")
}

func TestIssuesCmd_NoIssues(t *testing.T) {
	fc := &fakeClient{
		project: &api.Project{ID: "proj-1", Name: "Auth Rewrite"},
		issues:  &api.IssueListResult{Issues: nil},
	}
	f := newTestFactory(fc)

	cmd := NewProjectCmd(f)
	cmd.SetArgs([]string{"issues", "auth-rewrite"})
	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	buf := f.IO.Out.(*bytes.Buffer)
	out := buf.String()
	assert.Contains(t, out, "No issues found")
}
