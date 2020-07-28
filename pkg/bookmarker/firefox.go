package bookmarker

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/frioux/leatherman/pkg/mozlz4"
)

// firefoxBookmarkEntry a bookmark structure of decompressed .jsonlz4
type firefoxBookmarkEntry struct {
	GUID         string `json:"guid"`
	Title        string `json:"title"`
	Index        int    `json:"index"`
	DateAdded    int64  `json:"dateAdded"`
	LastModified int64  `json:"lastModified"`
	ID           int    `json:"id"`
	TypeCode     int    `json:"typeCode"`
	Type         string `json:"type"`
	Root         string `json:"root"`
	Charset      string `json:"charset,omitempty"`
	Iconuri      string `json:"iconuri,omitempty"`
	Annos        []struct {
		Name    string `json:"name"`
		Value   string `json:"value"`
		Expires int    `json:"expires"`
		Flags   int    `json:"flags"`
	} `json:"annos,omitempty"`
	URI      string                  `json:"uri,omitempty"`
	Children []*firefoxBookmarkEntry `json:"children,omitempty"`
}

// firefoxBookmarkRoot has a single entry as root which has childrens
type firefoxBookmarkRoot struct {
	root firefoxBookmarkEntry
}

type firefoxBookmark struct {
	bookmarkRoot firefoxBookmarkRoot
	bookmarkPath string
}

// NewFirefox returns a new firefox instance to get bookmarks
func NewFirefox(path string) Bookmarker {
	return &firefoxBookmark{
		bookmarkPath: path,
	}
}

// Bookmarks load firefox bookmark entries and return general bookmark structure
func (b *firefoxBookmark) Bookmarks() (bookmarks Bookmarks, err error) {
	if err = b.load(); err != nil {
		return
	}

	return b.bookmarkRoot.root.convertToBookmarks("/"), nil
}

// load a compressed .jsonlz4 file
func (b *firefoxBookmark) load() error {
	bookmarkMozlz4File := b.bookmarkPath

	f, err := os.Open(bookmarkMozlz4File)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := mozlz4.NewReader(f)
	if err != nil {
		return err
	}

	return json.NewDecoder(r).Decode(&b.bookmarkRoot.root)
}

// convertToBookmarks parse a entry and children of the entry
func (entry *firefoxBookmarkEntry) convertToBookmarks(folder string) (bookmarks Bookmarks) {
	// firefoxBookmarkEntry.TypeCode
	const (
		typeURI = iota + 1
		typeFolder
	)

	switch entry.TypeCode {
	case typeFolder:
		if entry.Children == nil {
			// if node has no entry, stop recursive
			return
		}
		if entry.Title != "" {
			// if entry type is folder, append folder name to current folder
			folder = filepath.Join(folder, entry.Title)
		}
		for _, e := range entry.Children {
			// tell the folder name to children bookmark entry
			bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
		}
	case typeURI:
		u, err := parseURL(entry.URI)
		if err != nil {
			return
		}

		b := &Bookmark{
			BookmarkerName: Firefox,
			Folder:         folder,
			Title:          entry.Title,
			URI:            entry.URI,
			Domain:         u.Host,
		}
		bookmarks = append(bookmarks, b)
	}

	return
}

// GetFirefoxBookmarkFile returns a firefox bookmark filepath in bookmark-backups direcotory
func GetFirefoxBookmarkFile(profile string) (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	profileDir := filepath.Join(home, "Library/Application Support/Firefox/Profiles")
	profileDirName, err := searchSuffixDir(profileDir, profile)
	if err != nil {
		return "", err
	}
	bookmarkDir := filepath.Join(profileDir, profileDirName, "bookmarkbackups")
	bookmarkFile, err := getLatestFile(bookmarkDir)
	if err != nil {
		return "", err
	}

	return bookmarkFile, nil
}
