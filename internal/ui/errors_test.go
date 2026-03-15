package ui

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintError_Plain(t *testing.T) {
	ios := NewTestIOStreams()
	ios.PrintError(errors.New("something went wrong"))

	buf, ok := ios.ErrOut.(*bytes.Buffer)
	require.True(t, ok)
	assert.Contains(t, buf.String(), "error:")
	assert.Contains(t, buf.String(), "something went wrong")
}

func TestPrintError_Styled(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	ios.PrintError(errors.New("something went wrong"))

	buf, ok := ios.ErrOut.(*bytes.Buffer)
	require.True(t, ok)
	assert.Contains(t, buf.String(), "something went wrong")
}
