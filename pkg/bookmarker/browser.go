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
const cacheFile = "alfred-bookmarks.cache"

var cacheDir string

// Browsers determine which bookmark read from
type Browsers struct {
	bookmarkers     map[browser]Bookmarker
	removeDuplicate bool
	cache           cache.Cacher
}

// Option is the type to replace default parameters.
type Option func(browsers *Browsers) error

func init() {
	cacheDir = os.TempDir()
}

// OptionFirefox if called, search firefox bookmark
func OptionFirefox(profile string) Option {
	return func(b *Browsers) error {
		path, err := GetFirefoxBookmarkFile(profile)
		if err != nil {
			return err
		}

		b.bookmarkers[firefox] = NewFirefox(path)
		return nil
	}
}

// OptionChrome if called, search chrome bookmark
func OptionChrome(profile string) Option {
	return func(b *Browsers) error {
		path, err := GetChromeBookmarkFile(profile)
		if err != nil {
			return err
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
		if age == 0 {
			age = 24
		} else if age < 0 {
			age = 0
		}

		c, err := cache.New(
			cacheDir, cacheFile,
			time.Duration(age)*time.Hour)
		if err != nil {
			return err
		}

		b.cache = c
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
		cache:       cache.NewNilCache(),
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
	bookmarks := Bookmarks{}
	if !browsers.cache.Expired() {
		if err := browsers.cache.Load(&bookmarks); err != nil {
			return Bookmarks{}, err
		}
		return bookmarks, nil
	}

	bookmarks, err := browsers.bookmarks()
	if err != nil {
		return Bookmarks{}, err
	}
	if err := browsers.cache.Store(&bookmarks); err != nil {
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
