package ui

import (
	"bytes"
	"io"
	"os"

	"golang.org/x/term"
)

// IOStreams provides access to standard I/O streams with TTY awareness.
type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
	plain  bool
}

// IsPlain returns true when output should be unstyled.
func (s *IOStreams) IsPlain() bool {
	return s.plain
}

// SetPlain overrides TTY auto-detection.
func (s *IOStreams) SetPlain(v bool) {
	s.plain = v
}

// NewIOStreams creates IOStreams with real stdio and TTY auto-detection.
func NewIOStreams() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
		plain:  !term.IsTerminal(int(os.Stdout.Fd())), //nolint:gosec // fd is always small
	}
}

// NewTestIOStreams creates IOStreams with buffers for testing.
func NewTestIOStreams() *IOStreams {
	return &IOStreams{
		In:     &bytes.Buffer{},
		Out:    &bytes.Buffer{},
		ErrOut: &bytes.Buffer{},
		plain:  true,
	}
}
