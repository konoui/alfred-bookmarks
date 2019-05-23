package bookmark

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/konoui/alfred-firefox-bookmarks/cache"
	"github.com/sahilm/fuzzy"
)

// Bookmark a instance of bookmark
type Bookmark struct {
	Folder string // Folder of Bookmarks
	Title  string // Bookmark title
	Domain string // Domain of URI
	URI    string // Bookmark URI
}

// Bookmarks a slice of Bookmark struct
type Bookmarks []*Bookmark

func (b Bookmarks) String(i int) string {
	return fmt.Sprintf("%s/%s", b[i].Folder, b[i].Title)
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
	cacheFile := "alfred-firefox-bookmarks.cache"
	expiredTime := 24 * time.Hour
	c, err := cache.NewCache(os.TempDir(), cacheFile, expiredTime)
	if err != nil {
		return Bookmarks{}, err
	}

	bookmarks := Bookmarks{}
	if c.Exists() && c.NotExpired() {
		if err = c.Load(&bookmarks); err != nil {
			return Bookmarks{}, err
		}
		return bookmarks, nil
	}

	bookmarks, err = LoadBookmarks()
	if err != nil {
		return Bookmarks{}, nil
	}
	if err = c.Store(&bookmarks); err != nil {
		return Bookmarks{}, err
	}

	return bookmarks, nil
}

// LoadBookmarks return Bookmarks struct, loading firefox bookmarks and parse them.
func LoadBookmarks() (Bookmarks, error) {
	entry := firefoxBookmarkEntry{}
	if err := LoadFirefoxBookmarkEntries(&entry); err != nil {
		return Bookmarks{}, err
	}

	bookmarks := convertEntriesToBookmarks(&entry, "")
	return bookmarks, nil
}

func convertEntriesToBookmarks(entry *firefoxBookmarkEntry, folder string) Bookmarks {
	if entry.Children == nil {
		return Bookmarks{}
	}

	bookmarks := Bookmarks{}
	for _, e := range entry.Children {
		switch e.TypeCode {
		case typeFolder:
			// If type of entry is folder, tell a entry of children the folder name of parent
			folder = e.Title
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
			b := &Bookmark{
				Folder: folder,
				Title:  e.Title,
				URI:    e.URI,
				Domain: u.Host,
			}
			bookmarks = append(bookmarks, b)
		}
		// tell the folder name to children bookmark entry
		bookmarks = append(bookmarks, convertEntriesToBookmarks(e, folder)...)
	}
	return bookmarks
}
