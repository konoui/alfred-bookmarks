package bookmarker

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/mitchellh/go-homedir"
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

type chromeBookmarkEntries struct {
	Checksum string `json:"checksum"`
	Roots    struct {
		BookmarkBar struct {
			BookmarkEntries []chromeBookmarkEntry `json:"children"`
		} `json:"bookmark_bar"`
		Other struct {
			BookmarkEntries []chromeBookmarkEntry `json:"children"`
		} `json:"other"`
		Synced struct {
			BookmarkEntries []chromeBookmarkEntry `json:"children"`
		} `json:"synced"`
	} `json:"roots"`
	Version int `json:"version"`
}

type chromeBookmark struct {
	chromeBookmarkEntries chromeBookmarkEntries
	bookmarkPath          string
}

// NewChrome return a new chrome instance to get bookmarks
func NewChrome(path string) Bookmarker {
	return &chromeBookmark{
		bookmarkPath: path,
	}
}

// Bookmarks load chrome bookmark entries and return general bookmark structure
func (b *chromeBookmark) Bookmarks() (Bookmarks, error) {
	if err := b.unmarshal(); err != nil {
		return Bookmarks{}, err
	}

	return b.chromeBookmarkEntries.convertToBookmarks(), nil
}

func (b *chromeBookmark) unmarshal() error {
	f, err := os.Open(b.bookmarkPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&b.chromeBookmarkEntries)
}

// convertToBookmarks parse a top of root entries which is array
func (entries *chromeBookmarkEntries) convertToBookmarks() Bookmarks {
	bookmarks := Bookmarks{}
	for _, entry := range entries.Roots.BookmarkBar.BookmarkEntries {
		bookmarks = append(bookmarks, entry.convertToBookmarks("")...)
	}

	return bookmarks
}

// convertToBookmarks parse a entry and children of the entry
func (entry *chromeBookmarkEntry) convertToBookmarks(folder string) Bookmarks {
	if entry.Type == "folder" && entry.Children == nil {
		return Bookmarks{}
	}

	bookmarks := Bookmarks{}
	if entry.Type == "url" {
		u, err := url.Parse(entry.URL)
		if err != nil {
			return Bookmarks{}
		}
		if u.Host == "" {
			return Bookmarks{}
		}

		b := &Bookmark{
			Browser: Chrome,
			Folder:  folder,
			Title:   entry.Name,
			URI:     entry.URL,
			Domain:  u.Host,
		}
		bookmarks = append(bookmarks, b)
		return bookmarks
	}

	if entry.Type == "folder" {
		// we append folder name to parent folder name
		folder = fmt.Sprintf("%s/%s", folder, entry.Name)
	}

	// loop folder type wihch has children
	for _, e := range entry.Children {
		bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
	}

	return bookmarks
}

// GetChromeBookmarkFile return a chrome bookmark filepath
func GetChromeBookmarkFile(profile string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	chromeDir := fmt.Sprintf("%s/Library/Application Support/Google/Chrome", home)
	profileDirName, err := searchSuffixDir(chromeDir, profile)
	if err != nil {
		return "", err
	}

	bookmarkFile := fmt.Sprintf("%s/%s/Bookmarks", chromeDir, profileDirName)
	return bookmarkFile, nil
}
