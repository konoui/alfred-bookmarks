package bookmarker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type chromeBookmarkEntry struct {
	DateAdded    string                 `json:"date_added"`
	GUID         string                 `json:"guid"`
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	URL          string                 `json:"url,omitempty"`
	Children     []*chromeBookmarkEntry `json:"children,omitempty"`
	DateModified string                 `json:"date_modified,omitempty"`
}

type chromeBookmarkRoot struct {
	Checksum string `json:"checksum"`
	Roots    struct {
		BookmarkBar *chromeBookmarkEntry `json:"bookmark_bar"`
		Other       *chromeBookmarkEntry `json:"other"`
		Synced      *chromeBookmarkEntry `json:"synced"`
	} `json:"roots"`
	Version int `json:"version"`
}

type chromeBookmark struct {
	bookmarkRoot chromeBookmarkRoot
	bookmarkPath string
}

// NewChrome returns a new chrome instance to get bookmarks
func NewChrome(path string) Bookmarker {
	return &chromeBookmark{
		bookmarkPath: path,
	}
}

// Bookmarks load chrome bookmark entries and return general bookmark structure
func (b *chromeBookmark) Bookmarks() (bookmarks Bookmarks, err error) {
	if err = b.load(); err != nil {
		return
	}

	barBookmarks := b.bookmarkRoot.Roots.BookmarkBar.convertToBookmarks("/")
	syncedBookmarks := b.bookmarkRoot.Roots.Synced.convertToBookmarks("/")
	othersBookmarks := b.bookmarkRoot.Roots.Other.convertToBookmarks("/")
	bookmarks = append(bookmarks, barBookmarks...)
	bookmarks = append(bookmarks, syncedBookmarks...)
	bookmarks = append(bookmarks, othersBookmarks...)

	return bookmarks, nil
}

// load a chrome bookmark file
func (b *chromeBookmark) load() error {
	f, err := os.Open(b.bookmarkPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&b.bookmarkRoot)
}

// convertToBookmarks parse a entry and children of the entry
func (entry *chromeBookmarkEntry) convertToBookmarks(folder string) (bookmarks Bookmarks) {
	switch entry.Type {
	case "folder":
		if entry.Children == nil {
			// if node has no entry, stop recursive
			return
		}
		if entry.Name != "" {
			// we append folder name to parent folder name
			folder = filepath.Join(folder, entry.Name)
		}
		for _, e := range entry.Children {
			bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
		}
	case "url":
		u, err := parseURL(entry.URL)
		if err != nil {
			return
		}

		b := &Bookmark{
			BookmarkerName: Chrome,
			Folder:         folder,
			Title:          entry.Name,
			URI:            entry.URL,
			Domain:         u.Host,
		}
		bookmarks = append(bookmarks, b)
	}

	return
}

// GetChromeBookmarkFile returns a chrome bookmark filepath
func GetChromeBookmarkFile(profilePath, profileName string) (string, error) {
	profileDirName, err := searchSuffixDir(profilePath, profileName)
	if err != nil {
		return "", err
	}

	bookmarkFile := filepath.Join(profilePath, profileDirName, "Bookmarks")
	if err := hasReadCapability(bookmarkFile); err != nil {
		return "", fmt.Errorf("chrome error: %w", err)
	}

	return bookmarkFile, nil
}
