package bookmarker

import (
	"path/filepath"
	"testing"
)

var testChromeBookmarkJSONFile = filepath.Join(testdataPath, "test-chrome-bookmarks.json")

var testChromeBookmarks = Bookmarks{
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar",
		Title:          "Google",
		Domain:         "www.google.com",
		URI:            "https://www.google.com/",
	},
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar/1-hierarchy-a",
		Title:          "GitHub",
		Domain:         "github.com",
		URI:            "https://github.com/",
	},
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar/1-hierarchy-a/2-hierarchy-a/3-hierarchy-a",
		Title:          "Stack Overflow",
		Domain:         "stackoverflow.com",
		URI:            "https://stackoverflow.com/",
	},
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar/1-hierarchy-a/2-hierarchy-a/3-hierarchy-a",
		Title:          "Amazon Web Services",
		Domain:         "aws.amazon.com",
		URI:            "https://aws.amazon.com/?nc1=h_ls",
	},
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar/1-hierarchy-b",
		Title:          "Yahoo",
		Domain:         "www.yahoo.com",
		URI:            "https://www.yahoo.com/",
	},
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar/1-hierarchy-b/2-hierarchy-a",
		Title:          "Facebook",
		Domain:         "www.facebook.com",
		URI:            "https://www.facebook.com/",
	},
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar/1-hierarchy-b/2-hierarchy-a",
		Title:          "Twitter",
		Domain:         "twitter.com",
		URI:            "https://twitter.com/login",
	},
	&Bookmark{
		BookmarkerName: Chrome,
		Folder:         "/Bookmarks Bar/1-hierarchy-b/2-hierarchy-b",
		Title:          "Amazon.com",
		Domain:         "www.amazon.com",
		URI:            "https://www.amazon.com/",
	},
}

func TestChromeBookmarks(t *testing.T) {
	tests := []struct {
		description  string
		bookmarkPath string
		want         Bookmarks
		expectErr    bool
	}{
		{
			description:  "valid bookmark file",
			bookmarkPath: testChromeBookmarkJSONFile,
			want:         testChromeBookmarks,
			expectErr:    false,
		},
		{
			description:  "invalid bookmark file",
			bookmarkPath: "test",
			expectErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			b := chromeBookmark{
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
