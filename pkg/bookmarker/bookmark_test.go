package bookmarker

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const testdataPath = "testdata"

// DiffBookmark is helper function that compare unsorted Bookmarks
// return "" if got is equal to want regardless of sorted or unsorted.
// format is "+want -got"
func DiffBookmark(want, got Bookmarks) string {
	sort.Slice(want, func(i, j int) bool {
		return want[i].URI < want[j].URI
	})
	sort.Slice(got, func(i, j int) bool {
		return got[i].URI < got[j].URI
	})

	return cmp.Diff(want, got)
}

func TestBookmarks_UniqByURI(t *testing.T) {
	tests := []struct {
		name string
		want Bookmarks
	}{
		{
			name: "enable firefox, chrome, safari and remove dupulication. return chrome bookmark",
			want: testChromeBookmarks,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getTestAllBookmarks(t)
			got := b.uniqByURI()
			if diff := DiffBookmark(got, tt.want); diff != "" {
				t.Errorf("+want/-got: %s", diff)
			}
		})
	}
}

func getTestBookmarks(t *testing.T, opts ...Option) Bookmarks {
	bookmarer, err := New(opts...)
	if err != nil {
		t.Fatal(err)
	}

	b, err := bookmarer.Bookmarks()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func getTestAllBookmarks(t *testing.T) Bookmarks {
	options := []Option{
		WithFirefox(defaultFirefoxProfilePath, testProfile),
		WithChrome(defaultChromeProfilePath, testProfile),
		WithSafari(),
	}
	return getTestBookmarks(t, options...)
}
