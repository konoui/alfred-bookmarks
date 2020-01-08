package bookmarker

import (
	"sort"

	"github.com/google/go-cmp/cmp"
)

// DiffBookmark is helper function that compare unsorted Bookmarks
// return "" if got is equal to want regardless of sorted or unsorted
func DiffBookmark(want, got Bookmarks) string {
	sort.Slice(want, func(i, j int) bool {
		return want[i].URI < want[j].URI
	})
	sort.Slice(got, func(i, j int) bool {
		return got[i].URI < got[j].URI
	})
	diff := cmp.Diff(got, want)

	return diff
}
