package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/misham/linear-cli/internal/api"
)

// resolveProject finds a project by ID, slug, or name.
// It first tries the arg as an ID/slug via GetProject.
// If that fails, it lists projects (scoped to team) and matches by name (case-insensitive).
func resolveProject(ctx context.Context, client api.Client, teamID, arg string) (*api.Project, error) {
	p, err := client.GetProject(ctx, arg)
	if err == nil {
		return p, nil
	}

	cursor := ""
	for {
		result, listErr := client.ListProjects(ctx, teamID, "", 50, cursor)
		if listErr != nil {
			return nil, fmt.Errorf("get project: %w", err)
		}

		for i := range result.Projects {
			if strings.EqualFold(result.Projects[i].Name, arg) {
				return &result.Projects[i], nil
			}
		}

		if !result.PageInfo.HasNextPage {
			break
		}
		cursor = result.PageInfo.EndCursor
	}

	return nil, fmt.Errorf("project %q not found", arg)
}
