package bookmark

import (
	"log"
	"net/url"

	"github.com/sahilm/fuzzy"
)

// Bookmark a instance of bookmark
type Bookmark struct {
	Title  string // Bookmark title
	Domain string // Domain of URI
	URI    string // Bookmark URI
}

// Bookmarks a slice of Bookmark struct
type Bookmarks []*Bookmark

func (b Bookmarks) String(i int) string {
	return b[i].Title
}

// Len return length of Bookmarks for fuzzy interface
func (b Bookmarks) Len() int {
	return len(b)
}

// Filter fuzzy search bookmarks using query
func (b Bookmarks) Filter(query string) Bookmarks {
	bookmarks := Bookmarks{}
	results := fuzzy.FindFrom(query, b)
	for _, r := range results {
		bookmarks = append(bookmarks, b[r.Index])
	}

	return bookmarks
}

// LoadBookmarksFromCache return Bookmarks struct, loading cache file.
func LoadBookmarksFromCache() (Bookmarks, error) {
	return LoadBookmarks()
}

// LoadBookmarks return Bookmarks struct, loading firefox bookmarks and parse them.
func LoadBookmarks() (Bookmarks, error) {
	entry := firefoxBookmarkEntry{}
	if err := LoadFirefoxBookmarkEntries(&entry); err != nil {
		return Bookmarks{}, err
	}

	bookmarks := Bookmarks{}
	entryToBookmarks(&entry, &bookmarks)
	return bookmarks, nil
}

func entryToBookmarks(entry *firefoxBookmarkEntry, bookmarks *Bookmarks) {
	if entry == nil {
		return
	}

	for _, e := range entry.Children {
		switch e.TypeCode {
		case typeFolder:
			//fmt.Printf("Folder: %s\n", e.Title)
		case typeURI:
			u, err := url.Parse(e.URI)
			// Ignore invalid URLs
			if err != nil {
				log.Printf("could not parse URL \"%s\" (%s): %v", e.URI, e.Title, err)
				continue
			}

			if u.Host == "" {
				log.Printf("Domain is empty \"%s\" (%s)", e.URI, e.Title)
				continue
			}

			bookmark := &Bookmark{
				Title:  e.Title,
				URI:    e.URI,
				Domain: u.Host,
			}
			*bookmarks = append(*bookmarks, bookmark)
		}
		entryToBookmarks(e, bookmarks)
	}
}
