package ui

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
)

// FileRef represents a file link extracted from markdown body text.
type FileRef struct {
	Name    string
	URL     string
	IsImage bool
}

var uploadLinkRe = regexp.MustCompile(`!?\[([^\]]*)\]\((https://uploads\.linear\.app/[^)]+)\)`)

// ExtractFiles parses markdown text for [name](https://uploads.linear.app/...) links.
// Returns deduplicated list (by URL). Accepts multiple text inputs.
func ExtractFiles(texts ...string) []FileRef {
	seen := make(map[string]bool)
	var files []FileRef

	for _, text := range texts {
		matches := uploadLinkRe.FindAllStringSubmatch(text, -1)
		for _, m := range matches {
			altText := m[1]
			url := m[2]

			if seen[url] {
				continue
			}
			seen[url] = true

			basename := urlBasename(url)
			name := altText
			if filepath.Ext(name) == "" {
				name = basename
			}

			files = append(files, FileRef{
				Name:    name,
				URL:     url,
				IsImage: IsImageFile(basename),
			})
		}
	}

	return files
}

// urlBasename extracts the filename from a URL, stripping query parameters.
func urlBasename(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return path.Base(rawURL)
	}
	return path.Base(parsed.Path)
}

// HyperlinkOSC8 renders text as a clickable terminal hyperlink using OSC 8 escape sequences.
// Format: ESC ] 8 ; ; URL ST text ESC ] 8 ; ; ST
func HyperlinkOSC8(url, text string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}
