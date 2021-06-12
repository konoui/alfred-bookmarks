package bookmarker

import (
	"sort"
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

func (b Bookmarks) uniqByURI() Bookmarks {
	m := make(map[string]bool)
	uniq := make(Bookmarks, 0, len(b))
	for _, e := range b {
		if !m[e.URI] {
			m[e.URI] = true
			uniq = append(uniq, e)
		}
	}

	return uniq
}
