package bookmarker

import "fmt"

// Manager determine which bookmark read from
type Manager struct {
	bookmarkers      map[bookmarkerName]Bookmarker
	removeDuplicates bool
}

// Option is the type to replace default parameters.
type Option func(m *Manager) error

// OptionFirefox if called, search firefox bookmark
func OptionFirefox(profilePath, profileName string) Option {
	return func(m *Manager) error {
		path, err := GetFirefoxBookmarkFile(profilePath, profileName)
		if err != nil {
			return err
		}

		m.bookmarkers[Firefox] = NewFirefox(path)
		return nil
	}
}

// OptionChrome if called, search chrome bookmark
func OptionChrome(profilePath, profileName string) Option {
	return func(m *Manager) error {
		path, err := GetChromeBookmarkFile(profilePath, profileName)
		if err != nil {
			return err
		}

		m.bookmarkers[Chrome] = NewChrome(path)
		return nil
	}
}

// OptionSafari if called, search safari bookmark
func OptionSafari() Option {
	return func(m *Manager) error {
		path, err := GetSafariBookmarkFile()
		if err != nil {
			return err
		}

		m.bookmarkers[Safari] = NewSafari(path)
		return nil
	}
}

// OptionRemoveDuplicates removes same bookmarks by urls
func OptionRemoveDuplicates() Option {
	return func(m *Manager) error {
		m.removeDuplicates = true
		return nil
	}
}

// New is a managed bookmarker to get each bookmarks
func New(opts ...Option) (Bookmarker, error) {
	m := &Manager{
		bookmarkers: make(map[bookmarkerName]Bookmarker),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(m); err != nil {
			return m, err
		}
	}

	return m, nil
}

// Bookmarks return Bookmarks struct by loading each bookmarker
func (m *Manager) Bookmarks() (bookmarks Bookmarks, err error) {
	for _, name := range getSupportedBookmarkerNames() {
		bookmarker, ok := m.bookmarkers[name]
		if !ok {
			continue
		}

		b, err := bookmarker.Bookmarks()
		if err != nil {
			// Note： not continue but return err if error occurs
			return bookmarks, fmt.Errorf("failed to load bookmarks in %s: %w", name, err)
		}
		bookmarks = append(bookmarks, b...)
	}

	// Note: execute uniq after folder filter
	if m.removeDuplicates {
		bookmarks = bookmarks.uniqByURI()
	}

	return bookmarks, nil
}
