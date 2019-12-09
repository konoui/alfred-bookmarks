package bookmark

import (
	"encoding/json"
	"os"
	"time"

	"github.com/konoui/alfred-bookmarks/pkg/cache"
	"github.com/sahilm/fuzzy"
)

// Bookmark abstract each browser bookmark as the structure
type Bookmark struct {
	Folder string // Folder of Bookmarks
	Title  string // Bookmark title
	Domain string // Domain of URI
	URI    string // Bookmark URI
}

// Bookmarker is interface to load each bookmark file
// TODO add Marshal/Unmarshal
type Bookmarker interface {
	Bookmarks() (Bookmarks, error)
}

// Bookmarks a slice of Bookmark struct
type Bookmarks []*Bookmark

type browser string

const (
	firefox browser = "firefox"
	chrome  browser = "chrome"
)

// Browsers determine which bookmark read from
type Browsers struct {
	bookmarkers map[browser]Bookmarker
	bookmarks   Bookmarks
}

// Option is the type to replace default parameters.
type Option func(browsers *Browsers) error

// OptionFirefox if called, search firefox bookmark
func OptionFirefox(path string) Option {
	return func(b *Browsers) error {
		if path == "" {
			var err error
			if path, err = GetFirefoxBookmarkFile(); err != nil {
				return err
			}
		}

		b.bookmarkers[firefox] = NewFirefoxBookmark(path)
		return nil
	}
}

// OptionChrome if called, search chrome bookmark
func OptionChrome(path string) Option {
	return func(b *Browsers) error {
		if path == "" {
			var err error
			if path, err = GetChromeBookmarkFile(); err != nil {
				return err
			}
		}

		b.bookmarkers[chrome] = NewChromeBookmark(path)
		return nil
	}
}

// OptionNone noop
func OptionNone() Option {
	return func(b *Browsers) error {
		return nil
	}
}

// NewBrowsers return Browsers
// TODO return bookmarker
func NewBrowsers(opts ...Option) *Browsers {
	b := &Browsers{
		bookmarkers: make(map[browser]Bookmarker),
	}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			panic(err)
		}
	}

	return b
}

// BookmarksFromCache return Bookmarks struct, loading cache file.
func (browsers *Browsers) BookmarksFromCache() (Bookmarks, error) {
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

	bookmarks, err = browsers.Bookmarks()
	if err != nil {
		return Bookmarks{}, err
	}
	if err = c.Store(&bookmarks); err != nil {
		return Bookmarks{}, err
	}

	return bookmarks, nil
}

// Bookmarks return Bookmarks struct, loading each browser bookmarks and parse them.
func (browsers *Browsers) Bookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	for _, bookmarker := range browsers.bookmarkers {
		b, err := bookmarker.Bookmarks()
		if err != nil {
			// Note： not continue but return err if error occurs
			return Bookmarks{}, err
		}
		bookmarks = append(bookmarks, b...)
	}

	browsers.bookmarks = bookmarks
	return bookmarks, nil
}

// MarshalJSON is used to serialize the type to json
func (browsers *Browsers) MarshalJSON() ([]byte, error) {
	return browsers.bookmarks.Marshal()
}

// UnmarshalJSON is used to deserialize json types into Conditional
func (browsers *Browsers) UnmarshalJSON(jsonData []byte) error {
	return browsers.bookmarks.Unmarshal(jsonData)
}

// Marshal is used to serialize the type to json
func (b Bookmarks) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

// Unmarshal is used to deserialize json types into Conditional
func (b Bookmarks) Unmarshal(jsonData []byte) error {
	return json.Unmarshal(jsonData, &b)
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
