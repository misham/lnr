package ui

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestIOStreams(t *testing.T) {
	ios := NewTestIOStreams()
	assert.NotNil(t, ios.Out)
	assert.NotNil(t, ios.ErrOut)
	assert.True(t, ios.IsPlain(), "test IOStreams should default to plain")
}

func TestIOStreams_SetPlain(t *testing.T) {
	ios := NewTestIOStreams()
	ios.SetPlain(false)
	assert.False(t, ios.IsPlain())
	ios.SetPlain(true)
	assert.True(t, ios.IsPlain())
}

func TestIOStreams_OutBuffer(t *testing.T) {
	ios := NewTestIOStreams()
	_, _ = ios.Out.Write([]byte("hello"))
	buf, ok := ios.Out.(*bytes.Buffer)
	require.True(t, ok)
	assert.Equal(t, "hello", buf.String())
}
