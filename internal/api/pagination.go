package api

// PageInfo contains pagination cursor state from GraphQL connections.
type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}
