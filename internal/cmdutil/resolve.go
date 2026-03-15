package cmdutil

import (
	"context"
	"fmt"
	"strings"
)

// ResolveTeamID returns a team UUID from the --team flag or stored config.
// If the value looks like a UUID (contains hyphens, 36 chars), it is used directly.
// Otherwise it is treated as a team key and resolved via ListTeams.
func ResolveTeamID(ctx context.Context, f *Factory) (string, error) {
	if f.TeamKey != "" {
		if looksLikeUUID(f.TeamKey) {
			return f.TeamKey, nil
		}
		return resolveKeyToID(ctx, f, f.TeamKey)
	}

	cfg, err := f.Config()
	if err != nil {
		return "", fmt.Errorf("load config: %w", err)
	}

	id := cfg.TeamID()
	if id == "" {
		return "", fmt.Errorf("no team set — use `lnr team set` or `--team` flag")
	}
	return id, nil
}

func looksLikeUUID(s string) bool {
	return len(s) == 36 && strings.Count(s, "-") == 4
}

func resolveKeyToID(ctx context.Context, f *Factory, key string) (string, error) {
	client, err := f.APIClient()
	if err != nil {
		return "", err
	}

	teams, err := client.ListTeams(ctx)
	if err != nil {
		return "", fmt.Errorf("list teams: %w", err)
	}

	for _, t := range teams {
		if strings.EqualFold(t.Key, key) {
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("team %q not found", key)
}
