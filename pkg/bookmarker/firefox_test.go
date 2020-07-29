package bookmarker

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pierrec/lz4"
)

var testFirefoxBookmarkJsonlz4File = filepath.Join(testdataPath, "test-firefox-bookmarks.jsonlz4")
var testFirefoxBookmarkJSONFile = filepath.Join(testdataPath, "test-firefox-bookmarks.json")
var testFirefoxBookmarks = Bookmarks{
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu",
		Title:          "Google",
		Domain:         "www.google.com",
		URI:            "https://www.google.com/",
	},
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu/1-hierarchy-a",
		Title:          "GitHub",
		Domain:         "github.com",
		URI:            "https://github.com/",
	},
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu/1-hierarchy-a/2-hierarchy-a/3-hierarchy-a",
		Title:          "Stack Overflow",
		Domain:         "stackoverflow.com",
		URI:            "https://stackoverflow.com/",
	},
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu/1-hierarchy-a/2-hierarchy-a/3-hierarchy-a",
		Title:          "Amazon Web Services",
		Domain:         "aws.amazon.com",
		URI:            "https://aws.amazon.com/?nc1=h_ls",
	},
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu/1-hierarchy-b",
		Title:          "Yahoo",
		Domain:         "www.yahoo.com",
		URI:            "https://www.yahoo.com/",
	},
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu/1-hierarchy-b/2-hierarchy-a",
		Title:          "Facebook",
		Domain:         "www.facebook.com",
		URI:            "https://www.facebook.com/",
	},
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu/1-hierarchy-b/2-hierarchy-a",
		Title:          "Twitter",
		Domain:         "twitter.com",
		URI:            "https://twitter.com/login",
	},
	&Bookmark{
		BookmarkerName: Firefox,
		Folder:         "/Bookmark Menu/1-hierarchy-b/2-hierarchy-b",
		Title:          "Amazon.com",
		Domain:         "www.amazon.com",
		URI:            "https://www.amazon.com/",
	},
}

func TestFirefoxBookmarks(t *testing.T) {
	tests := []struct {
		description  string
		bookmarkPath string
		want         Bookmarks
		expectErr    bool
	}{
		{
			description:  "valid bookmark file",
			bookmarkPath: testFirefoxBookmarkJsonlz4File,
			want:         testFirefoxBookmarks,
		},
		{
			description:  "invalid bookmark file",
			bookmarkPath: "test",
			expectErr:    true,
		},
	}

	setupFirefox(t)
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			b := firefoxBookmark{
				bookmarkPath: tt.bookmarkPath,
			}

			bookmarks, err := b.Bookmarks()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}

			diff := DiffBookmark(bookmarks, tt.want)
			if !tt.expectErr && diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func setupFirefox(t *testing.T) {
	t.Helper()
	if err := createTestFirefoxJsonlz4(); err != nil {
		t.Fatal(err)
	}
}

// create Jsonlz4 from jsonfile
func createTestFirefoxJsonlz4() error {
	// switch readDefaultFirefoxBookmarksJSON or readTestFirefoxBookmarkJSON
	_, _ = readDefaultFirefoxBookmarksJSON()
	str, err := readTestFirefoxBookmarkJSON()
	if err != nil {
		return err
	}

	w, err := os.Create(testFirefoxBookmarkJsonlz4File)
	if err != nil {
		return err
	}
	defer w.Close()

	r := strings.NewReader(str)
	err = compress(r, w, len(str))
	if err != nil {
		return err
	}
	return nil
}

func readTestFirefoxBookmarkJSON() (string, error) {
	jsonData, err := ioutil.ReadFile(testFirefoxBookmarkJSONFile)

	return string(jsonData), err
}

// return json string of .jsonlz4 loading from local profile
func readDefaultFirefoxBookmarksJSON() (string, error) {
	path, err := GetFirefoxBookmarkFile("default")
	if err != nil {
		return "", err
	}

	b := firefoxBookmark{
		bookmarkPath: path,
		bookmarkRoot: firefoxBookmarkRoot{},
	}

	err = b.load()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(b.bookmarkRoot.root)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// encode json to .jsonlz4
func compress(src io.Reader, dst io.Writer, intendedSize int) error {
	const magicHeader = "mozLz40\x00"
	_, err := dst.Write([]byte(magicHeader))
	if err != nil {
		return fmt.Errorf("couldn't Write header: %w", err)
	}

	b, err := ioutil.ReadAll(src)
	if err != nil {
		return fmt.Errorf("couldn't ReadAll to Compress: %w", err)
	}

	err = binary.Write(dst, binary.LittleEndian, uint32(intendedSize))
	if err != nil {
		return fmt.Errorf("couldn't encode length: %w", err)
	}

	dstBytes := make([]byte, 10*len(b))
	sz, err := lz4.CompressBlockHC(b, dstBytes, -1)
	if err != nil {
		return fmt.Errorf("couldn't CompressBlock: %w", err)
	}
	if sz == 0 {
		return errors.New("data incompressible")
	}

	_, err = dst.Write(dstBytes[:sz])
	if err != nil {
		return fmt.Errorf("couldn't Write compressed data: %w", err)
	}

	return nil
}
