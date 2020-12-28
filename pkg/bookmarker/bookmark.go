package bookmarker

import (
	"sort"
	"strings"
)

// bookmarkerName is a type of supported browser name
type bookmarkerName string

const (
	// Firefox is supported
	Firefox bookmarkerName = "firefox"
	// Chrome is supported
	Chrome bookmarkerName = "chrome"
	// Safari is supported
	Safari bookmarkerName = "safari"
)

func getSupportedBookmarkerNames() []bookmarkerName {
	names := []bookmarkerName{
		Firefox,
		Chrome,
		Safari,
	}
	// sort by name asc for making idempotency result
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	return names
}

// Bookmark abstract each browser bookmark
type Bookmark struct {
	BookmarkerName bookmarkerName
	Folder         string
	Title          string
	Domain         string
	URI            string
}

// Bookmarker is a interface to load each bookmark file
type Bookmarker interface {
	Bookmarks() (Bookmarks, error)
}

// Bookmarks a slice of Bookmark struct
type Bookmarks []*Bookmark

func (b Bookmarks) uniqByURI() (uniq Bookmarks) {
	m := make(map[string]bool)

	for _, e := range b {
		if !m[e.URI] {
			m[e.URI] = true
			uniq = append(uniq, e)
		}
	}

	return
}

func (b Bookmarks) filterByFolderPrefix(query string) (fb Bookmarks) {
	if query == "" {
		return b
	}

	for _, e := range b {
		if hasFolderPrefix(e.Folder, query) {
			fb = append(fb, e)
		}
	}

	return
}

func hasFolderPrefix(folder, prefix string) bool {
	folder = strings.ToLower(folder)
	folder = strings.ReplaceAll(folder, " ", "")
	prefix = strings.ToLower(prefix)
	prefix = strings.ReplaceAll(prefix, " ", "")
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	if strings.HasPrefix(folder, prefix) {
		return true
	}

	return strings.HasPrefix(folder+"/", prefix)
}
