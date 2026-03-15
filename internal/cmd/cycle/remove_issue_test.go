package cycle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveIssueCmd_CallsRemoveIssueCycle(t *testing.T) {
	fc := &fakeClient{}

	f := newTestFactory(t, fc)
	cmd := newRemoveIssueCmd(f)
	cmd.SetArgs([]string{"issue-1"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "issue-1", fc.removeCycleIssueID)
}

func TestRemoveIssueCmd_RequiresOneArg(t *testing.T) {
	fc := &fakeClient{}

	f := newTestFactory(t, fc)
	cmd := newRemoveIssueCmd(f)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
}
