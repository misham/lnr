package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageInfo(t *testing.T) {
	pi := PageInfo{
		HasNextPage: true,
		EndCursor:   "cursor-abc",
	}
	assert.True(t, pi.HasNextPage)
	assert.Equal(t, "cursor-abc", pi.EndCursor)
}
