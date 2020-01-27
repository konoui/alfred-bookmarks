package bookmarker

import (
	"encoding/json"
	"sort"

	"github.com/sahilm/fuzzy"
)

// Bookmark abstract each browser bookmark as the structure
type Bookmark struct {
	Browser Browser
	Folder  string
	Title   string
	Domain  string
	URI     string
}

// Bookmarker is interface to load each bookmark file
// TODO add Marshal/Unmarshal
type Bookmarker interface {
	Bookmarks() (Bookmarks, error)
}

// Bookmarks a slice of Bookmark struct
type Bookmarks []*Bookmark

func (b Bookmarks) uniqByURI() Bookmarks {
	m := make(map[string]bool)
	uniq := Bookmarks{}

	sort.Slice(b, func(i, j int) bool {
		return b[i].Browser < b[j].Browser
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

// Marshal is used to serialize the type to json
func (b Bookmarks) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

// Unmarshal is used to deserialize json types into Conditional
func (b Bookmarks) Unmarshal(jsonData []byte) error {
	return json.Unmarshal(jsonData, &b)
}

// String retrun a bookmark title of index for fuzzy interface
func (b Bookmarks) String(i int) string {
	return b[i].Title
}

// Len return length of Bookmarks for fuzzy interface
func (b Bookmarks) Len() int {
	return len(b)
}

// Filter fuzzy search Bookmarks using query
func (b Bookmarks) Filter(query string) Bookmarks {
	bookmarks := Bookmarks{}
	results := fuzzy.FindFrom(query, b)
	for _, r := range results {
		bookmarks = append(bookmarks, b[r.Index])
	}

	return bookmarks
}
