package bookmarker

import (
	"os"
	"time"

	"github.com/konoui/alfred-bookmarks/pkg/cacher"
	"github.com/pkg/errors"
)

var cacheDir string

// Browser　 is a type of supported engine
type Browser string

const (
	// Firefox is supported
	Firefox Browser = "firefox"
	// Chrome is supported
	Chrome    Browser = "chrome"
	cacheFile string  = "alfred-bookmarks.cache"
)

// Browsers determine which bookmark read from
type Browsers struct {
	bookmarkers     map[Browser]Bookmarker
	removeDuplicate bool
	cacher          cacher.Cacher
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

		b.bookmarkers[Firefox] = NewFirefox(path)
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

		b.bookmarkers[Chrome] = NewChrome(path)
		return nil
	}
}

// OptionRemoveDuplicate remove same url bookmark e.g) search from multi browsers
func OptionRemoveDuplicate() Option {
	return func(b *Browsers) error {
		b.removeDuplicate = true
		return nil
	}
}

// OptionCacheMaxAge is bookmark cache time. unit indicate hours
// if passed arg is zero, set 24 hours. if passed arg is minus, disable cache
func OptionCacheMaxAge(hour int) Option {
	return func(b *Browsers) error {
		if hour == 0 {
			hour = 24
		} else if hour < 0 {
			hour = 0
		}

		c, err := cacher.New(
			cacheDir, cacheFile,
			time.Duration(hour)*time.Hour)
		if err != nil {
			return err
		}

		b.cacher = c
		return nil
	}
}

// OptionNone noop
func OptionNone() Option {
	return func(b *Browsers) error {
		return nil
	}
}

// NewBrowsers is instance to get Bookmarks of multi browser
func NewBrowsers(opts ...Option) Bookmarker {
	b := &Browsers{
		bookmarkers: make(map[Browser]Bookmarker),
		cacher:      cacher.NewNilCache(),
	}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			panic(err)
		}
	}

	return b
}

// Bookmarks return Bookmarks struct by loading cache file
func (browsers *Browsers) Bookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	if !browsers.cacher.Expired() {
		if err := browsers.cacher.Load(&bookmarks); err != nil {
			return Bookmarks{}, errors.Wrap(err, "failed to load cache data")
		}
		return bookmarks, nil
	}

	bookmarks, err := browsers.bookmarks()
	if err != nil {
		return Bookmarks{}, err
	}
	if err := browsers.cacher.Store(&bookmarks); err != nil {
		return Bookmarks{}, errors.Wrap(err, "failed to save data into cache")
	}

	return bookmarks, nil
}

// bookmarks return Bookmarks struct by loading each browser
func (browsers *Browsers) bookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	for browser, bookmarker := range browsers.bookmarkers {
		b, err := bookmarker.Bookmarks()
		if err != nil {
			// Note： not continue but return err if error occurs
			return Bookmarks{}, errors.Wrapf(err, "failed to load bookmarks in %s", browser)
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
