package bookmark

import (
	"os"
	"time"

	"github.com/konoui/alfred-bookmarks/pkg/cache"
	"github.com/sahilm/fuzzy"
)

// Browsers determine which bookmark read from
type Browsers struct {
	firefox bool
	chrome  bool
}

// Bookmark a instance of bookmark
type Bookmark struct {
	Folder string // Folder of Bookmarks
	Title  string // Bookmark title
	Domain string // Domain of URI
	URI    string // Bookmark URI
}

// Bookmarks a slice of Bookmark struct
type Bookmarks []*Bookmark

// NewBrowsers return Browser
func NewBrowsers(firefox, chrome bool) *Browsers {
	return &Browsers{
		chrome:  chrome,
		firefox: firefox,
	}
}

// LoadBookmarksFromCache return Bookmarks struct, loading cache file.
func (b *Browsers) LoadBookmarksFromCache() (Bookmarks, error) {
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

	bookmarks, err = b.LoadBookmarks()
	if err != nil {
		return Bookmarks{}, err
	}
	if err = c.Store(&bookmarks); err != nil {
		return Bookmarks{}, err
	}

	return bookmarks, nil
}

// LoadBookmarks return Bookmarks struct, loading browser bookmarks and parse them.
func (b *Browsers) LoadBookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	// FIXME implement interface and loop LoadBookmarkEntries() and convertToBookmarks("")
	if b.firefox {
		entry := firefoxBookmarkEntry{}
		if err := entry.LoadBookmarkEntries(); err != nil {
			return Bookmarks{}, err
		}

		bookmarks = append(bookmarks, entry.convertToBookmarks("")...)
	}

	if b.chrome {
		entries := chromeBookmarkEntries{}
		if err := entries.LoadBookmarkEntries(); err != nil {
			return Bookmarks{}, err
		}

		bookmarks = append(bookmarks, entries.convertToBookmarks("")...)
	}

	return bookmarks, nil
}

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
