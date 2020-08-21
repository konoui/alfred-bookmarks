package bookmarker

import (
	"github.com/pkg/errors"
)

// engine determine which bookmark read from
type engine struct {
	bookmarkers     map[bookmarkerName]Bookmarker
	removeDuplicate bool
}

// Option is the type to replace default parameters.
type Option func(e *engine) error

// OptionFirefox if called, search firefox bookmark
func OptionFirefox(profilePath, profileName string) Option {
	return func(e *engine) error {
		path, err := GetFirefoxBookmarkFile(profilePath, profileName)
		if err != nil {
			return err
		}

		e.bookmarkers[Firefox] = NewFirefox(path)
		return nil
	}
}

// OptionChrome if called, search chrome bookmark
func OptionChrome(profilePath, profileName string) Option {
	return func(e *engine) error {
		path, err := GetChromeBookmarkFile(profilePath, profileName)
		if err != nil {
			return err
		}

		e.bookmarkers[Chrome] = NewChrome(path)
		return nil
	}
}

// OptionSafari if called, search safari bookmark
func OptionSafari() Option {
	return func(e *engine) error {
		path, err := GetSafariBookmarkFile()
		if err != nil {
			return err
		}

		e.bookmarkers[Safari] = NewSafari(path)
		return nil
	}
}

// OptionRemoveDuplicate removes same bookmarks by urls
func OptionRemoveDuplicate() Option {
	return func(e *engine) error {
		e.removeDuplicate = true
		return nil
	}
}

// New is a managed bookmarker to get each bookmarks
func New(opts ...Option) (Bookmarker, error) {
	e := &engine{
		bookmarkers: make(map[bookmarkerName]Bookmarker),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(e); err != nil {
			return e, err
		}
	}

	return e, nil
}

// Bookmarks return Bookmarks struct by loading each bookmarker
func (e *engine) Bookmarks() (Bookmarks, error) {
	bookmarks := Bookmarks{}
	for _, name := range getSupportedBookmarkerNames() {
		bookmarker, ok := e.bookmarkers[name]
		if !ok {
			continue
		}

		b, err := bookmarker.Bookmarks()
		if err != nil {
			// Note： not continue but return err if error occurs
			return bookmarks, errors.Wrapf(err, "failed to load bookmarks in %s", name)
		}
		bookmarks = append(bookmarks, b...)
	}

	if e.removeDuplicate {
		bookmarks = bookmarks.uniqByURI()
	}

	return bookmarks, nil
}
