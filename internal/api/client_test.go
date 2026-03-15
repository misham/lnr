package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphQLClient_ListTeams(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"teams": map[string]any{
				"nodes": []map[string]any{
					{
						"id":          "team-1",
						"name":        "Engineering",
						"key":         "ENG",
						"description": "Core engineering",
						"private":     false,
						"icon":        "",
						"color":       "#0000FF",
						"timezone":    "America/Los_Angeles",
					},
					{
						"id":          "team-2",
						"name":        "Design",
						"key":         "DES",
						"description": "Product design",
						"private":     true,
						"icon":        "",
						"color":       "#FF0000",
						"timezone":    "Europe/London",
					},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	teams, err := client.ListTeams(context.Background())
	require.NoError(t, err)
	require.Len(t, teams, 2)

	assert.Equal(t, "team-1", teams[0].ID)
	assert.Equal(t, "Engineering", teams[0].Name)
	assert.Equal(t, "ENG", teams[0].Key)
	assert.Equal(t, "Core engineering", teams[0].Description)
	assert.False(t, teams[0].Private)
	assert.Equal(t, "#0000FF", teams[0].Color)

	assert.Equal(t, "team-2", teams[1].ID)
	assert.Equal(t, "Design", teams[1].Name)
	assert.True(t, teams[1].Private)
}

func TestGraphQLClient_Viewer(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"viewer": map[string]any{
				"id":          "user-1",
				"name":        "John Doe",
				"displayName": "johnd",
				"email":       "john@example.com",
				"active":      true,
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	user, err := client.Viewer(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "user-1", user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "johnd", user.DisplayName)
	assert.Equal(t, "john@example.com", user.Email)
	assert.True(t, user.Active)
}

func newTestServer(resp map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestGraphQLClient_ListIssues(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	resp := map[string]any{
		"data": map[string]any{
			"issues": map[string]any{
				"nodes": []map[string]any{
					{
						"id": "issue-1", "identifier": "ENG-1", "title": "Fix bug",
						"description": "A bug", "priority": 2.0, "priorityLabel": "High",
						"estimate": 3.0, "dueDate": "2026-02-01", "url": "https://linear.app/issue/ENG-1",
						"createdAt": now.Format(time.RFC3339), "updatedAt": now.Format(time.RFC3339),
						"completedAt": "0001-01-01T00:00:00Z", "archivedAt": "0001-01-01T00:00:00Z",
						"state":    map[string]any{"id": "state-1", "name": "In Progress", "type": "started", "color": "#00FF00", "position": 1.0},
						"team":     map[string]any{"id": "team-1", "name": "Engineering", "key": "ENG"},
						"assignee": map[string]any{"id": "user-1", "name": "John", "displayName": "johnd", "email": "john@test.com", "active": true},
						"labels":   map[string]any{"nodes": []map[string]any{{"id": "label-1", "name": "bug", "color": "#FF0000"}}},
					},
				},
				"pageInfo": map[string]any{"hasNextPage": true, "endCursor": "cursor-1"},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	result, err := client.ListIssues(context.Background(), "team-1", 50, "")
	require.NoError(t, err)
	require.Len(t, result.Issues, 1)

	issue := result.Issues[0]
	assert.Equal(t, "issue-1", issue.ID)
	assert.Equal(t, "ENG-1", issue.Identifier)
	assert.Equal(t, "Fix bug", issue.Title)
	assert.Equal(t, 2, issue.Priority)
	assert.Equal(t, 3, issue.Estimate)
	assert.Equal(t, "In Progress", issue.State.Name)
	assert.Equal(t, "ENG", issue.Team.Key)
	assert.Equal(t, "John", issue.Assignee.Name)
	require.Len(t, issue.Labels, 1)
	assert.Equal(t, "bug", issue.Labels[0].Name)
	assert.True(t, result.PageInfo.HasNextPage)
	assert.Equal(t, "cursor-1", result.PageInfo.EndCursor)
}

func TestGraphQLClient_GetIssue(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	resp := map[string]any{
		"data": map[string]any{
			"issue": map[string]any{
				"id": "issue-1", "identifier": "ENG-1", "title": "Fix bug",
				"description": "A bug", "priority": 1.0, "priorityLabel": "Urgent",
				"estimate": 0.0, "dueDate": "", "url": "https://linear.app/issue/ENG-1",
				"createdAt": now.Format(time.RFC3339), "updatedAt": now.Format(time.RFC3339),
				"completedAt": now.Format(time.RFC3339), "archivedAt": nil,
				"state":    map[string]any{"id": "state-1", "name": "Done", "type": "completed", "color": "#00FF00", "position": 5.0},
				"team":     map[string]any{"id": "team-1", "name": "Engineering", "key": "ENG"},
				"assignee": map[string]any{"id": "user-1", "name": "John", "displayName": "johnd", "email": "john@test.com", "active": true},
				"labels":   map[string]any{"nodes": []any{}},
				"comments": map[string]any{
					"nodes": []map[string]any{
						{
							"id": "comment-1", "body": "Looks good", "createdAt": now.Format(time.RFC3339), "updatedAt": now.Format(time.RFC3339),
							"user": map[string]any{"id": "user-2", "name": "Jane", "displayName": "janed", "email": "jane@test.com", "active": true},
						},
					},
				},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	issue, err := client.GetIssue(context.Background(), "issue-1")
	require.NoError(t, err)

	assert.Equal(t, "issue-1", issue.ID)
	assert.Equal(t, "ENG-1", issue.Identifier)
	assert.Equal(t, 1, issue.Priority)
	assert.NotNil(t, issue.CompletedAt)
	assert.Nil(t, issue.ArchivedAt)
	require.Len(t, issue.Comments, 1)
	assert.Equal(t, "Looks good", issue.Comments[0].Body)
	assert.Equal(t, "Jane", issue.Comments[0].User.Name)
}

func TestGraphQLClient_SearchIssues(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	resp := map[string]any{
		"data": map[string]any{
			"searchIssues": map[string]any{
				"nodes": []map[string]any{
					{
						"id": "issue-2", "identifier": "ENG-2", "title": "Search result",
						"description": "", "priority": 3.0, "priorityLabel": "Normal",
						"estimate": 0.0, "dueDate": "", "url": "https://linear.app/issue/ENG-2",
						"createdAt": now.Format(time.RFC3339), "updatedAt": now.Format(time.RFC3339),
						"completedAt": "0001-01-01T00:00:00Z", "archivedAt": "0001-01-01T00:00:00Z",
						"state":    map[string]any{"id": "state-1", "name": "Todo", "type": "unstarted", "color": "#888", "position": 0.0},
						"team":     map[string]any{"id": "team-1", "name": "Engineering", "key": "ENG"},
						"assignee": map[string]any{"id": "", "name": "", "displayName": "", "email": "", "active": false},
						"labels":   map[string]any{"nodes": []any{}},
					},
				},
				"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	result, err := client.SearchIssues(context.Background(), "search", "team-1", 50, "")
	require.NoError(t, err)
	require.Len(t, result.Issues, 1)
	assert.Equal(t, "ENG-2", result.Issues[0].Identifier)
	assert.False(t, result.PageInfo.HasNextPage)
}

func TestGraphQLClient_CreateIssue(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"issueCreate": map[string]any{
				"success": true,
				"issue": map[string]any{
					"id": "issue-new", "identifier": "ENG-99", "title": "New issue",
					"url":   "https://linear.app/issue/ENG-99",
					"state": map[string]any{"id": "state-1", "name": "Backlog", "type": "backlog"},
				},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	issue, err := client.CreateIssue(context.Background(), IssueCreateInput{
		Title:  "New issue",
		TeamID: "team-1",
	})
	require.NoError(t, err)
	assert.Equal(t, "issue-new", issue.ID)
	assert.Equal(t, "ENG-99", issue.Identifier)
	assert.Equal(t, "Backlog", issue.State.Name)
}

func TestGraphQLClient_UpdateIssue(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"issueUpdate": map[string]any{
				"success": true,
				"issue": map[string]any{
					"id": "issue-1", "identifier": "ENG-1", "title": "Updated title",
					"url":   "https://linear.app/issue/ENG-1",
					"state": map[string]any{"id": "state-2", "name": "In Progress", "type": "started"},
				},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	title := "Updated title"
	client := NewGraphQLClient(srv.URL, "test-token")
	issue, err := client.UpdateIssue(context.Background(), "issue-1", IssueUpdateInput{Title: &title})
	require.NoError(t, err)
	assert.Equal(t, "Updated title", issue.Title)
}

func TestGraphQLClient_ArchiveIssue(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"issueArchive": map[string]any{"success": true},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	err := client.ArchiveIssue(context.Background(), "issue-1")
	require.NoError(t, err)
}

func TestGraphQLClient_ListWorkflowStates(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"workflowStates": map[string]any{
				"nodes": []map[string]any{
					{"id": "ws-1", "name": "Backlog", "type": "backlog", "color": "#bbb", "position": 0.0},
					{"id": "ws-2", "name": "In Progress", "type": "started", "color": "#0f0", "position": 1.0},
				},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	states, err := client.ListWorkflowStates(context.Background(), "team-1")
	require.NoError(t, err)
	require.Len(t, states, 2)
	assert.Equal(t, "Backlog", states[0].Name)
	assert.Equal(t, "started", states[1].Type)
}

func TestGraphQLClient_ListComments(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	resp := map[string]any{
		"data": map[string]any{
			"issue": map[string]any{
				"comments": map[string]any{
					"nodes": []map[string]any{
						{
							"id": "c-1", "body": "Comment 1", "createdAt": now.Format(time.RFC3339), "updatedAt": now.Format(time.RFC3339),
							"user": map[string]any{"id": "u-1", "name": "Alice", "displayName": "alice", "email": "alice@test.com", "active": true},
						},
					},
					"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
				},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	comments, pageInfo, err := client.ListComments(context.Background(), "issue-1", 50, "")
	require.NoError(t, err)
	require.Len(t, comments, 1)
	assert.Equal(t, "Comment 1", comments[0].Body)
	assert.Equal(t, "Alice", comments[0].User.Name)
	assert.False(t, pageInfo.HasNextPage)
}

func TestGraphQLClient_CreateComment(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	resp := map[string]any{
		"data": map[string]any{
			"commentCreate": map[string]any{
				"success": true,
				"comment": map[string]any{
					"id": "c-new", "body": "New comment", "createdAt": now.Format(time.RFC3339),
					"user": map[string]any{"id": "u-1", "name": "Alice", "displayName": "alice", "email": "alice@test.com", "active": true},
				},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	comment, err := client.CreateComment(context.Background(), "issue-1", "New comment")
	require.NoError(t, err)
	assert.Equal(t, "c-new", comment.ID)
	assert.Equal(t, "New comment", comment.Body)
}

func TestGraphQLClient_AddIssueLabel(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"issueAddLabel": map[string]any{"success": true},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	err := client.AddIssueLabel(context.Background(), "issue-1", "label-1")
	require.NoError(t, err)
}

func TestGraphQLClient_RemoveIssueLabel(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"issueRemoveLabel": map[string]any{"success": true},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	err := client.RemoveIssueLabel(context.Background(), "issue-1", "label-1")
	require.NoError(t, err)
}

func TestGraphQLClient_ListLabels(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"issueLabels": map[string]any{
				"nodes": []map[string]any{
					{"id": "label-1", "name": "bug", "color": "#FF0000"},
					{"id": "label-2", "name": "feature", "color": "#00FF00"},
				},
			},
		},
	}
	srv := newTestServer(resp)
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "test-token")
	labels, err := client.ListLabels(context.Background(), "team-1")
	require.NoError(t, err)
	require.Len(t, labels, 2)
	assert.Equal(t, "bug", labels[0].Name)
	assert.Equal(t, "feature", labels[1].Name)
}

func TestGraphQLClient_ListTeams_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "bad-token")
	_, err := client.ListTeams(context.Background())
	require.Error(t, err)
}

func TestGraphQLClient_AuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		resp := map[string]any{
			"data": map[string]any{
				"teams": map[string]any{
					"nodes": []any{},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewGraphQLClient(srv.URL, "my-token")
	_, _ = client.ListTeams(context.Background())

	assert.Equal(t, "Bearer my-token", gotAuth)
}
