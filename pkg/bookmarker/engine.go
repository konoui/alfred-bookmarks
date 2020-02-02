package bookmarker

import (
	"os"
	"time"

	"github.com/konoui/alfred-bookmarks/pkg/cacher"
	"github.com/pkg/errors"
)

const cacheFile string = "alfred-bookmarks.cache"

var cacheDir string = os.TempDir()

// engine determine which bookmark read from
type engine struct {
	bookmarkers     map[name]Bookmarker
	removeDuplicate bool
	cacher          cacher.Cacher
}

// Option is the type to replace default parameters.
type Option func(e *engine) error

// OptionFirefox if called, search firefox bookmark
func OptionFirefox(profile string) Option {
	return func(e *engine) error {
		path, err := GetFirefoxBookmarkFile(profile)
		if err != nil {
			return err
		}

		e.bookmarkers[Firefox] = NewFirefox(path)
		return nil
	}
}

// OptionChrome if called, search chrome bookmark
func OptionChrome(profile string) Option {
	return func(e *engine) error {
		path, err := GetChromeBookmarkFile(profile)
		if err != nil {
			return err
		}

		e.bookmarkers[Chrome] = NewChrome(path)
		return nil
	}
}

// OptionRemoveDuplicate remove same url bookmark e.g) search from each bookmarker
func OptionRemoveDuplicate() Option {
	return func(e *engine) error {
		e.removeDuplicate = true
		return nil
	}
}

// OptionCacheMaxAge is bookmark cache time. unit indicate hours
// if passed arg is zero, set 24 hours. if passed arg is minus, disable cache
func OptionCacheMaxAge(hour int) Option {
	return func(e *engine) error {
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

		e.cacher = c
		return nil
	}
}

// OptionNone noop
func OptionNone() Option {
	return func(e *engine) error {
		return nil
	}
}

// New is instance to get Bookmarks of each bookmarker
func New(opts ...Option) Bookmarker {
	e := &engine{
		bookmarkers: make(map[name]Bookmarker),
		cacher:      cacher.NewNilCache(),
	}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			panic(err)
		}
	}

	return e
}

// Bookmarks return Bookmarks struct by loading cache file
func (e *engine) Bookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	if !e.cacher.Expired() {
		if err := e.cacher.Load(&bookmarks); err != nil {
			return Bookmarks{}, errors.Wrap(err, "failed to load cache data")
		}
		return bookmarks, nil
	}

	bookmarks, err := e.bookmarks()
	if err != nil {
		return Bookmarks{}, err
	}
	if err := e.cacher.Store(&bookmarks); err != nil {
		return Bookmarks{}, errors.Wrap(err, "failed to save data into cache")
	}

	return bookmarks, nil
}

// bookmarks return Bookmarks struct by loading each bookmarker
func (e *engine) bookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	for name, bookmarker := range e.bookmarkers {
		b, err := bookmarker.Bookmarks()
		if err != nil {
			// Note： not continue but return err if error occurs
			return Bookmarks{}, errors.Wrapf(err, "failed to load bookmarks in %s", name)
		}
		bookmarks = append(bookmarks, b...)
	}

	if e.removeDuplicate {
		bookmarks = bookmarks.uniqByURI()
	}

	return bookmarks, nil
}

// Marshal is used to serialize the type to json
func (e *engine) Marshal() ([]byte, error) {
	b, err := e.Bookmarks()
	if err != nil {
		return []byte{}, err
	}
	return b.Marshal()
}

// Unmarshal is used to deserialize json types into Conditional
func (e *engine) Unmarshal(jsonData []byte) error {
	b, err := e.Bookmarks()
	if err != nil {
		return err
	}
	return b.Unmarshal(jsonData)
}
