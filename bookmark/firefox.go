package bookmark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

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

type firefoxBookmarkEntries []*firefoxBookmarkEntry

// firefoxBookmarkEntry.TypeCode
const (
	typeURI = iota + 1
	typeFolder
)

// LoadFirefoxBookmarkEntries Read data int entry from decompressed .jsonlz4 file of bookmarkbackups direcotory
func LoadFirefoxBookmarkEntries(entry *firefoxBookmarkEntry) error {
	bookmarkMozlz4File, err := GetBookmarkFile()
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

// GetBookmarkFile return firefox bookmarkbackups direcotory
func GetBookmarkFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	profileDir := fmt.Sprintf("%s/Library/Application Support/Firefox/Profiles", usr.HomeDir)
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
