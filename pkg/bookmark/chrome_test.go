package bookmark

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestLoadBookmarkEntries(t *testing.T) {
	tests := []struct {
		description string
	}{
		{
			description: "load chrome bookmarks",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			entries := chromeBookmarkEntries{}
			entries.LoadBookmarkEntries()
			if len(entries.Roots.BookmarkBar.BookmarkEntries) == 0 {
				t.Errorf("bookmark has no entry %+v", entries)
			}
		})
	}
}

func TestConvertToBookmarks(t *testing.T) {
	tests := []struct {
		description string
	}{
		{
			description: "convert chrome bookmark to general bookmark struct",
		},
	}
	for _, tt := range tests {
		log.SetOutput(os.Stdout)
		t.Run(tt.description, func(t *testing.T) {
			entries := chromeBookmarkEntries{}
			entries.LoadBookmarkEntries()
			bookmarks := Bookmarks{}
			for _, e := range entries.Roots.BookmarkBar.BookmarkEntries {
				bookmarks = append(bookmarks, e.convertToBookmarks("")...)
			}

			for i, b := range bookmarks {
				fmt.Printf("%d: %+v\n", i, b.Title)
			}
		})
	}
}
