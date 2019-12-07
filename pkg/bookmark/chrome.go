package bookmark

import (
	"encoding/json"
	"fmt"
	"log"
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

// NewChromeBookmark return new instance
func NewChromeBookmark(path string) Bookmarker {
	return &chromeBookmark{
		bookmarkPath: path,
	}
}

// LoadBookmark load chrome bookmark to general bookmark structure
func (b *chromeBookmark) LoadBookmarks() (Bookmarks, error) {
	if err := b.unmarshal(); err != nil {
		return Bookmarks{}, err
	}

	return b.chromeBookmarkEntries.convertToBookmarks(""), nil
}

func (b *chromeBookmark) unmarshal() error {
	f, err := os.Open(b.bookmarkPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&b.chromeBookmarkEntries)
}

// convertToBookmarks parse top of root entries which is array
func (entries *chromeBookmarkEntries) convertToBookmarks(folder string) Bookmarks {
	bookmarks := Bookmarks{}
	for _, entry := range entries.Roots.BookmarkBar.BookmarkEntries {
		bookmarks = append(bookmarks, entry.convertToBookmarks("")...)
	}

	return bookmarks
}

// convertToBookmarks parse a entry and the entry children
func (entry *chromeBookmarkEntry) convertToBookmarks(folder string) Bookmarks {
	if entry.Type == "folder" && entry.Children == nil {
		return Bookmarks{}
	}

	bookmarks := Bookmarks{}
	if entry.Type == "url" {
		u, err := url.Parse(entry.URL)
		// Ignore invalid URLs
		if err != nil {
			log.Printf("could not parse URL \"%s\" (%s): %v", entry.URL, entry.Name, err)
			return Bookmarks{}
		}

		if u.Host == "" {
			log.Printf("Domain is empty \"%s\" (%s)", entry.URL, entry.Name)
			return Bookmarks{}
		}
		b := &Bookmark{
			Folder: folder,
			Title:  entry.Name,
			URI:    entry.URL,
			Domain: u.Host,
		}
		bookmarks = append(bookmarks, b)
		return bookmarks
	}

	// loop folder wihch has children
	for _, e := range entry.Children {
		// entry.Name should be folder name
		folder = fmt.Sprintf("%s/%s", folder, entry.Name)
		bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
	}

	return bookmarks
}

// GetChromeBookmarkFile return chrome bookmark direcotory file
func GetChromeBookmarkFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	bookmarkFile := fmt.Sprintf("%s/Library/Application Support/Google/Chrome/Default/Bookmarks", home)

	return bookmarkFile, nil
}
