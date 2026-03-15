package initiative

import (
	"context"
	"fmt"
	"strings"

	"github.com/misham/linear-cli/internal/api"
)

// resolveInitiative finds an initiative by ID or name.
// It first tries the arg as an ID via GetInitiative.
// If that fails, it lists all initiatives and matches by name (case-insensitive).
func resolveInitiative(ctx context.Context, client api.Client, arg string) (*api.Initiative, error) {
	ini, err := client.GetInitiative(ctx, arg)
	if err == nil {
		return ini, nil
	}

	cursor := ""
	for {
		result, listErr := client.ListInitiatives(ctx, "", 50, cursor)
		if listErr != nil {
			return nil, fmt.Errorf("get initiative: %w", err)
		}

		for i := range result.Initiatives {
			if strings.EqualFold(result.Initiatives[i].Name, arg) {
				return &result.Initiatives[i], nil
			}
		}

		if !result.PageInfo.HasNextPage {
			break
		}
		cursor = result.PageInfo.EndCursor
	}

	return nil, fmt.Errorf("initiative %q not found", arg)
}
