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

func (entries *chromeBookmarkEntries) convertToBookmarks(folder string) Bookmarks {
	bookmarks := Bookmarks{}
	for _, e := range entries.Roots.BookmarkBar.BookmarkEntries {
		bookmarks = append(bookmarks, e.convertToBookmarks("")...)
	}

	return bookmarks
}

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

	for _, e := range entry.Children {
		bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
	}

	return bookmarks
}

// LoadBookmarkEntries read data into entry
func (entry *chromeBookmarkEntries) LoadBookmarkEntries() error {
	filename, err := GetChromeBookmarkFile()
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(entry)
}

// GetChromeBookmarkFile return firefox bookmarkbackups direcotory
func GetChromeBookmarkFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	bookmarkFile := fmt.Sprintf("%s/Library/Application Support/Google/Chrome/Default/Bookmarks", home)

	return bookmarkFile, nil
}
