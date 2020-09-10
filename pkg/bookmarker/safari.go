package bookmarker

import (
	"os"
	"path/filepath"

	"howett.net/plist"
)

type safariBookmarkEntry struct {
	Title           string                 `plist:"Title"`
	WebBookmarkType string                 `plist:"WebBookmarkType"`
	URLString       string                 `plist:"URLString"`
	WebBookmarkUUID string                 `plist:"WebBookmarkUUID"`
	URIDictionary   map[string]string      `plist:"URIDictionary"`
	Children        []*safariBookmarkEntry `plist:"Children"`
}

type safariBookmarkRoot struct {
	root safariBookmarkEntry
}

type safariBookmark struct {
	bookmarkRoot safariBookmarkRoot
	bookmarkPath string
}

// NewSafari returns a new safari instance to get bookmarks
func NewSafari(path string) Bookmarker {
	return &safariBookmark{
		bookmarkPath: path,
	}
}

// Bookmarks load firefox bookmark entries and return general bookmark structure
func (b *safariBookmark) Bookmarks() (bookmarks Bookmarks, err error) {
	if err = b.load(); err != nil {
		return
	}

	return b.bookmarkRoot.root.convertToBookmarks("/"), nil
}

func (b *safariBookmark) load() error {
	file, err := os.Open(b.bookmarkPath)
	if err != nil {
		return err
	}

	err = plist.NewDecoder(file).Decode(&b.bookmarkRoot.root)
	if err != nil {
		return err
	}

	return nil
}

// convertToBookmarks parse a entry and children of the entry
func (entry *safariBookmarkEntry) convertToBookmarks(folder string) (bookmarks Bookmarks) {
	switch entry.WebBookmarkType {
	case "WebBookmarkTypeList":
		if entry.Children == nil {
			return
		}
		if entry.Title != "" {
			folder = filepath.Join(folder, entry.Title)
		}
		for _, e := range entry.Children {
			bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
		}
	case "WebBookmarkTypeLeaf":
		u, err := parseURL(entry.URLString)
		if err != nil {
			return
		}

		title, ok := entry.URIDictionary["title"]
		if !ok {
			title = "undefined"
		}

		b := &Bookmark{
			BookmarkerName: Safari,
			Folder:         folder,
			Title:          title,
			URI:            entry.URLString,
			Domain:         u.Host,
		}
		bookmarks = append(bookmarks, b)
	}

	return
}

// GetSafariBookmarkFile returns a safari bookmark filepath
func GetSafariBookmarkFile() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	bookmarkFile := filepath.Join(home, "Library", "Safari", "Bookmarks.plist")
	return bookmarkFile, nil
}
