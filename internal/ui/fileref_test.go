package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractFiles_SingleLink(t *testing.T) {
	text := "Check this [report.pdf](https://uploads.linear.app/org/file.pdf)"
	files := ExtractFiles(text)
	require.Len(t, files, 1)
	assert.Equal(t, "report.pdf", files[0].Name)
	assert.Equal(t, "https://uploads.linear.app/org/file.pdf", files[0].URL)
	assert.False(t, files[0].IsImage) // PDF is not an image
}

func TestExtractFiles_ImageLink(t *testing.T) {
	text := "![screenshot](https://uploads.linear.app/org/screenshot.png)"
	files := ExtractFiles(text)
	require.Len(t, files, 1)
	assert.Equal(t, "screenshot.png", files[0].Name) // "screenshot" has no ext, falls back to URL basename
	assert.True(t, files[0].IsImage)                 // determined from URL path
}

func TestExtractFiles_ImageLinkNameFallback(t *testing.T) {
	// When alt text has no file extension, Name falls back to URL basename
	text := "![a screenshot](https://uploads.linear.app/org/image.png)"
	files := ExtractFiles(text)
	require.Len(t, files, 1)
	assert.Equal(t, "image.png", files[0].Name) // falls back to URL basename
	assert.True(t, files[0].IsImage)
}

func TestExtractFiles_AltTextWithExtension(t *testing.T) {
	// When alt text has a file extension, use it as-is
	text := "[chart-reviews.pdf](https://uploads.linear.app/org/abc123)"
	files := ExtractFiles(text)
	require.Len(t, files, 1)
	assert.Equal(t, "chart-reviews.pdf", files[0].Name) // alt text has extension
}

func TestExtractFiles_MultipleTexts(t *testing.T) {
	desc := "[file1.pdf](https://uploads.linear.app/org/f1)"
	comment := "[file2.txt](https://uploads.linear.app/org/f2)"
	files := ExtractFiles(desc, comment)
	require.Len(t, files, 2)
	assert.Equal(t, "file1.pdf", files[0].Name)
	assert.Equal(t, "file2.txt", files[1].Name)
}

func TestExtractFiles_Deduplication(t *testing.T) {
	text1 := "[file.pdf](https://uploads.linear.app/org/same)"
	text2 := "[file.pdf](https://uploads.linear.app/org/same)"
	files := ExtractFiles(text1, text2)
	require.Len(t, files, 1)
}

func TestExtractFiles_ExternalURLsIgnored(t *testing.T) {
	text := "[link](https://github.com/some/repo) and [doc](https://uploads.linear.app/org/f1)"
	files := ExtractFiles(text)
	require.Len(t, files, 1)
	assert.Contains(t, files[0].URL, "uploads.linear.app")
}

func TestExtractFiles_NoFiles(t *testing.T) {
	files := ExtractFiles("just some plain text with no links")
	assert.Empty(t, files)
}

func TestExtractFiles_MixedContent(t *testing.T) {
	text := `
Here is a [report](https://uploads.linear.app/org/report.pdf) and
some [GitHub link](https://github.com/foo) and
![img](https://uploads.linear.app/org/img.jpg)
`
	files := ExtractFiles(text)
	require.Len(t, files, 2)
	assert.Equal(t, "report.pdf", files[0].Name) // "report" has no ext, falls back to URL basename
	assert.False(t, files[0].IsImage)
	assert.True(t, files[1].IsImage)
}

func TestExtractFiles_URLWithQueryParams(t *testing.T) {
	text := "[report.pdf](https://uploads.linear.app/org/report.pdf?sig=abc123&exp=456)"
	files := ExtractFiles(text)
	require.Len(t, files, 1)
	assert.Equal(t, "report.pdf", files[0].Name) // alt text has extension, use as-is
	assert.False(t, files[0].IsImage)
}

func TestExtractFiles_ImageURLWithQueryParams(t *testing.T) {
	text := "![screenshot](https://uploads.linear.app/org/image.png?sig=abc123)"
	files := ExtractFiles(text)
	require.Len(t, files, 1)
	assert.Equal(t, "image.png", files[0].Name) // fallback to URL basename without query
	assert.True(t, files[0].IsImage)            // should detect .png even with query params
}

func TestHyperlinkOSC8(t *testing.T) {
	result := HyperlinkOSC8("https://example.com", "Click here")
	assert.Contains(t, result, "Click here")
	// OSC 8 sequence should be present (exact format depends on lipgloss)
	assert.True(t, strings.Contains(result, "\033]8;") || strings.Contains(result, "\x1b]8;"),
		"expected OSC 8 escape sequence")
}
