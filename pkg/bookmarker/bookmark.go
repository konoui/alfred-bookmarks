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
)

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
	uniq := Bookmarks{}

	// Note: we sotrt by bookmarker name for making idempotency
	sort.Slice(b, func(i, j int) bool {
		return b[i].BookmarkerName < b[j].BookmarkerName
	})

	for _, e := range b {
		if !m[e.URI] {
			m[e.URI] = true
			uniq = append(uniq, e)
		}
	}

	b = uniq
	return uniq
}
