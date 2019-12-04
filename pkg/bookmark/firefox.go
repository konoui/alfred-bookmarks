package bookmark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

// firefoxBookmarkEntry.TypeCode
const (
	typeURI = iota + 1
	typeFolder
)

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
				log.Printf("could not parse URL \"%s\" (%s): %v", e.URI, e.Title, err)
				continue
			}

			if u.Host == "" {
				log.Printf("Domain is empty \"%s\" (%s)", e.URI, e.Title)
				continue
			}
			b := &Bookmark{
				Folder: folder,
				Title:  e.Title,
				URI:    e.URI,
				Domain: u.Host,
			}
			bookmarks = append(bookmarks, b)
		}
		// tell the folder name to children bookmark entry
		bookmarks = append(bookmarks, e.convertToBookmarks(folder)...)
	}

	return bookmarks
}

// LoadBookmarkEntries Read data into entry from decompressed .jsonlz4 file of bookmarkbackups direcotory
func (entry *firefoxBookmarkEntry) LoadBookmarkEntries() error {
	bookmarkMozlz4File, err := GetFirefoxBookmarkFile()
	if err != nil {
		return err
	}

	f, err := os.Open(bookmarkMozlz4File)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := mozlz4.NewReader(f)
	if err != nil {
		return err
	}

	return json.NewDecoder(r).Decode(entry)
}

// GetFirefoxBookmarkFile return firefox bookmarkbackups direcotory
func GetFirefoxBookmarkFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	profileDir := fmt.Sprintf("%s/Library/Application Support/Firefox/Profiles", home)
	defaultDirName, err := searchSuffixDir(profileDir, "default")
	if err != nil {
		return "", err
	}
	bookmarkDir := fmt.Sprintf("%s/%s/bookmarkbackups/", profileDir, defaultDirName)
	bookmarkFile, err := getLatestFile(bookmarkDir)
	if err != nil {
		return "", err
	}
	return bookmarkFile, nil
}

// getLatestFile return the path to latest files of dir.
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

// searchSuffixDir return directory name of suffix
func searchSuffixDir(dir, suffux string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if name := file.Name(); file.IsDir() && strings.HasSuffix(name, suffux) {
			return name, nil
		}
	}
	return "", fmt.Errorf("not found a directory of suffix (%s) in a directory (%s)", suffux, dir)
}
