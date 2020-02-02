package bookmarker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/frioux/leatherman/pkg/mozlz4"
	"github.com/mitchellh/go-homedir"
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

type firefoxBookmark struct {
	firefoxBookmarkEntry firefoxBookmarkEntry
	bookmarkPath         string
}

// firefoxBookmarkEntry.TypeCode
const (
	typeURI = iota + 1
	typeFolder
)

// NewFirefox return a new firefox instance to get bookmarks
func NewFirefox(path string) Bookmarker {
	return &firefoxBookmark{
		bookmarkPath: path,
	}
}

// Bookmarks load firefox bookmark entry and return general bookmark structure
func (b *firefoxBookmark) Bookmarks() (Bookmarks, error) {
	if err := b.load(); err != nil {
		return Bookmarks{}, err
	}

	return b.firefoxBookmarkEntry.convertToBookmarks(""), nil
}

// load a compressed as .jsonlz4 file
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

	return json.NewDecoder(r).Decode(&b.firefoxBookmarkEntry)
}

// convertToBookmarks parse a top of root entry
// we assume firefox has one root entry that has many children
func (entry *firefoxBookmarkEntry) convertToBookmarks(folder string) Bookmarks {
	if entry.Children == nil {
		return Bookmarks{}
	}

	// if entry type is folder, append folder name to current folder
	if entry.TypeCode == typeFolder {
		if entry.Title != "" {
			folder = fmt.Sprintf("%s/%s", folder, entry.Title)
		}
	}

	bookmarks := Bookmarks{}
	for _, e := range entry.Children {
		switch e.TypeCode {
		case typeFolder:
		case typeURI:
			u, err := url.Parse(e.URI)
			// Ignore invalid URLs
			if err != nil {
				continue
			}
			if u.Host == "" {
				continue
			}

			b := &Bookmark{
				Bookmarker: Firefox,
				Folder:     folder,
				Title:      e.Title,
				URI:        e.URI,
				Domain:     u.Host,
			}
			bookmarks = append(bookmarks, b)
		}
		// tell the folder name to children bookmark entry
		bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
	}

	return bookmarks
}

// GetFirefoxBookmarkFile return a firefox bookmark filepath in bookmarkbackups direcotory
func GetFirefoxBookmarkFile(profile string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	profileDir := fmt.Sprintf("%s/Library/Application Support/Firefox/Profiles", home)
	profileDirName, err := searchSuffixDir(profileDir, profile)
	if err != nil {
		return "", err
	}
	bookmarkDir := fmt.Sprintf("%s/%s/bookmarkbackups/", profileDir, profileDirName)
	bookmarkFile, err := getLatestFile(bookmarkDir)
	if err != nil {
		return "", err
	}

	return bookmarkFile, nil
}

// getLatestFile return a path to latest files in dir
func getLatestFile(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	latestIndex := 0
	for i, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if time.Since(file.ModTime()) <= time.Since(files[latestIndex].ModTime()) {
			latestIndex = i
		}
	}

	return filepath.Join(dir, files[latestIndex].Name()), nil
}

// searchSuffixDir return a directory name of suffix ignoring case-sensitive
func searchSuffixDir(dir, suffux string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if name := file.Name(); file.IsDir() &&
			strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffux)) {
			return name, nil
		}
	}

	return "", fmt.Errorf("not found a directory of suffix (%s) in %s directory", suffux, dir)
}
