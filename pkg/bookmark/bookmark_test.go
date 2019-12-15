package bookmark

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBrowsersBookmarks(t *testing.T) {
	tests := []struct {
		description string
		options     []Option
		want        Bookmarks
		expectErr   bool
	}{
		{
			description: "enable firefox bookmark",
			options: []Option{
				OptionFirefox(testFirefoxBookmarkJsonlz4File),
			},
			want:      testFirefoxBookmarks,
			expectErr: false,
		},
		{
			description: "enable chrome bookmark",
			options: []Option{
				OptionChrome(testChromeBookmarkJSONFile),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
		{
			description: "enable firefox, chrome remove dupulication. return return firefox bookmark",
			options: []Option{
				OptionFirefox(testFirefoxBookmarkJsonlz4File),
				OptionChrome(testChromeBookmarkJSONFile),
				OptionRemoveDuplicate(),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			browsers := NewBrowsers(tt.options...)
			bookmarks, err := browsers.Bookmarks()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err.Error())
			}

			diff := DiffBookmark(bookmarks, tt.want)
			if diff != "" {
				t.Errorf("unexpected response: (+want -got)\n%+v", diff)
			}
		})
	}
}

func TestBrowsersMarshaUnmarshalJson(t *testing.T) {
	tests := []struct {
		description string
		options     []Option
		expectErr   bool
	}{
		{
			description: "enable firefox bookmark",
			options: []Option{
				OptionFirefox(testFirefoxBookmarkJsonlz4File),
			},
			expectErr: false,
		},
		{
			description: "enable chrome bookmark",
			options: []Option{
				OptionChrome(testChromeBookmarkJSONFile),
			},
			expectErr: false,
		},
		{
			description: "enable firefox and chrome bookmark",
			options: []Option{
				OptionFirefox(testFirefoxBookmarkJsonlz4File),
				OptionChrome(testChromeBookmarkJSONFile),
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
		return got[i].URI < got[j].URI
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].URI < want[j].URI
	})
	diff := cmp.Diff(got, want)

	return diff
}
