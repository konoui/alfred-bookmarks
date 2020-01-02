package bookmarker

import (
	"os"
	"time"

	"github.com/konoui/alfred-bookmarks/pkg/cache"
)

type browser string

const (
	firefox browser = "firefox"
	chrome  browser = "chrome"
)

// Browsers determine which bookmark read from
type Browsers struct {
	bookmarkers     map[browser]Bookmarker
	removeDuplicate bool
	cacheMaxAge     time.Duration
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

		b.bookmarkers[firefox] = NewFirefox(path)
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

		b.bookmarkers[chrome] = NewChrome(path)
		return nil
	}
}

// OptionRemoveDuplicate remove same url bookmark e.g) search from multi browser
func OptionRemoveDuplicate() Option {
	return func(b *Browsers) error {
		b.removeDuplicate = true
		return nil
	}
}

// OptionCacheMaxAge bookmark cache time. unit indicate hours
// if passed arg is zero, set 24 hours. if passed arg is minus, set 0 hours
func OptionCacheMaxAge(age int) Option {
	return func(b *Browsers) error {
		if age < 0 {
			age = 0
		} else if age == 0 {
			age = 24
		}
		b.cacheMaxAge = time.Duration(age) * time.Hour
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
func NewBrowsers(opts ...Option) Bookmarker {
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

// Bookmarks return Bookmarks struct, loading cache file.
func (browsers *Browsers) Bookmarks() (Bookmarks, error) {
	cacheFile := "alfred-bookmarks.cache"
	bookmarks := Bookmarks{}
	c, err := cache.New(os.TempDir(), cacheFile, browsers.cacheMaxAge)
	if err != nil {
		return Bookmarks{}, err
	}

	if c.Exists() && c.NotExpired() {
		if err = c.Load(&bookmarks); err != nil {
			return Bookmarks{}, err
		}
		return bookmarks, nil
	}

	bookmarks, err = browsers.bookmarks()
	if err != nil {
		return Bookmarks{}, err
	}
	if err = c.Store(&bookmarks); err != nil {
		return Bookmarks{}, err
	}

	return bookmarks, nil
}

// bookmarks return Bookmarks struct, loading each browser bookmarks and parse them.
func (browsers *Browsers) bookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	for _, bookmarker := range browsers.bookmarkers {
		b, err := bookmarker.Bookmarks()
		if err != nil {
			// Note： not continue but return err if error occurs
			return Bookmarks{}, err
		}
		bookmarks = append(bookmarks, b...)
	}

	if browsers.removeDuplicate {
		bookmarks = bookmarks.uniqByURI()
	}

	return bookmarks, nil
}

// Marshal is used to serialize the type to json
func (browsers *Browsers) Marshal() ([]byte, error) {
	b, err := browsers.Bookmarks()
	if err != nil {
		return []byte{}, err
	}
	return b.Marshal()
}

// Unmarshal is used to deserialize json types into Conditional
func (browsers *Browsers) Unmarshal(jsonData []byte) error {
	b, err := browsers.Bookmarks()
	if err != nil {
		return err
	}
	return b.Unmarshal(jsonData)
}
