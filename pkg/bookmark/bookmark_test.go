package bookmark

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBrowsersBookmarks(t *testing.T) {
}

func TestBookmarksMarshaUnmarshalJson(t *testing.T) {
	tests := []struct {
		description string
		expectErr   bool
	}{
		{
			description: "",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt)
	}

}
func TestBrowsersMarshaUnmarshalJson(t *testing.T) {
	tests := []struct {
		description string
		options     []Option
		expectErr   bool
	}{
		{
			description: "enable chrome bookmark",
			options: []Option{
				OptionChrome("test-chrome-bookmarks.json"),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			browsers := NewBrowsers(tt.options...)
			bookmarks, err := browsers.Bookmarks()
			if err != nil {
				t.Fatal(err)
			}

			jsonData, err := browsers.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if err := browsers.UnmarshalJSON(jsonData); err != nil {
				t.Fatal(err)
			}

			diff := DiffBookmark(bookmarks, browsers.bookmarks)
			if diff != "" {
				t.Errorf("unexpected response: (+want -got)\n%+v", diff)
			}
		})
	}
}

func DiffBookmark(got, want Bookmarks) string {
	sort.Slice(got, func(i, j int) bool {
		return got[i].Domain < got[j].Domain
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].Domain < want[j].Domain
	})
	diff := cmp.Diff(got, want)

	return diff
}
