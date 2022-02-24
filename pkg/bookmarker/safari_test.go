package bookmarker

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"howett.net/plist"
)

var testSafariBookmarkPlist = filepath.Join(testdataPath, "test-safari-bookmarks.plist")
var testSafariBookmarkJSONFile = filepath.Join(testdataPath, "test-safari-bookmarks.json")
var testSafariBookmarks = Bookmarks{
	&Bookmark{
		BookmarkerName: Safari,
		Folder:         "/1-hierarchy-a/2-hierarchy-a/3-hierarchy-a",
		Title:          "Stack Overflow",
		Domain:         "stackoverflow.com",
		URI:            "https://stackoverflow.com/",
	},
	&Bookmark{
		BookmarkerName: Safari,
		Folder:         "/1-hierarchy-a/2-hierarchy-a/3-hierarchy-a",
		Title:          "Amazon Web Services",
		Domain:         "aws.amazon.com",
		URI:            "https://aws.amazon.com/?nc1=h_ls",
	},
	&Bookmark{
		BookmarkerName: Safari,
		Folder:         "/1-hierarchy-b",
		Title:          "Yahoo",
		Domain:         "www.yahoo.com",
		URI:            "https://www.yahoo.com/",
	},
	&Bookmark{
		BookmarkerName: Safari,
		Folder:         "/1-hierarchy-b/2-hierarchy-a",
		Title:          "Facebook",
		Domain:         "www.facebook.com",
		URI:            "https://www.facebook.com/",
	},
	&Bookmark{
		BookmarkerName: Safari,
		Folder:         "/1-hierarchy-b/2-hierarchy-a",
		Title:          "Twitter",
		Domain:         "twitter.com",
		URI:            "https://twitter.com/login",
	},
	&Bookmark{
		BookmarkerName: Safari,
		Folder:         "/1-hierarchy-b/2-hierarchy-b",
		Title:          "Amazon.com",
		Domain:         "www.amazon.com",
		URI:            "https://www.amazon.com/",
	},
}

func TestSafariBookmarks(t *testing.T) {
	tests := []struct {
		name         string
		bookmarkPath string
		want         Bookmarks
		expectErr    bool
	}{
		{
			name:         "valid bookmark file",
			bookmarkPath: testSafariBookmarkPlist,
			want:         testSafariBookmarks,
		},
		{
			name:         "invalid bookmark file",
			bookmarkPath: "test",
			expectErr:    true,
		},
	}

	setupSafari(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewSafari(tt.bookmarkPath)
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

func setupSafari(t *testing.T) {
	t.Helper()
	if err := createTestSafariPlistFile(); err != nil {
		t.Fatal(err)
	}
}

func createTestSafariPlistFile() error {
	_, _ = readLocalBookmarkPlist()
	jsonData, err := os.ReadFile(testSafariBookmarkJSONFile)
	if err != nil {
		return err
	}

	b := safariBookmark{}
	err = json.Unmarshal(jsonData, &b.bookmarkRoot.root)
	if err != nil {
		return err
	}

	w, err := os.Create(testSafariBookmarkPlist)
	if err != nil {
		return err
	}
	defer w.Close()

	return generatePlist(&b.bookmarkRoot, w)
}

func readLocalBookmarkPlist() (string, error) {
	path, err := GetSafariBookmarkFile()
	if err != nil {
		return "", err
	}
	b := &safariBookmark{
		bookmarkPath: path,
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

func generatePlist(r *safariBookmarkRoot, out io.Writer) error {
	return plist.NewEncoder(out).Encode(r.root)
}
