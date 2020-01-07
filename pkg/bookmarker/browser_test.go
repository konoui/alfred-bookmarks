package bookmarker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func testOptionChrome(path string) Option {
	return func(b *Browsers) error {
		b.bookmarkers[chrome] = NewChrome(path)
		return nil
	}
}

func testOptionFirefox(path string) Option {
	return func(b *Browsers) error {
		b.bookmarkers[firefox] = NewFirefox(path)
		return nil
	}
}

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
				testOptionFirefox(testFirefoxBookmarkJsonlz4File),
			},
			want:      testFirefoxBookmarks,
			expectErr: false,
		},
		{
			description: "enable chrome bookmark",
			options: []Option{
				testOptionChrome(testChromeBookmarkJSONFile),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
		{
			description: "enable firefox, chrome, remove dupulication. return chrome bookmark",
			options: []Option{
				testOptionFirefox(testFirefoxBookmarkJsonlz4File),
				testOptionChrome(testChromeBookmarkJSONFile),
				OptionRemoveDuplicate(),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
		{
			description: "enable firefox, chrome, remove dupulication, cacheOption. return chrome bookmark",
			options: []Option{
				testOptionFirefox(testFirefoxBookmarkJsonlz4File),
				testOptionChrome(testChromeBookmarkJSONFile),
				OptionRemoveDuplicate(),
				OptionCacheMaxAge(0),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			browsers := NewBrowsers(tt.options...)
			browsers.(*Browsers).cache.Clear()

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

func TestOptionFirefoxChrome(t *testing.T) {
	tests := []struct {
		description string
		options     []Option
		want        bool
	}{
		{
			description: "Lower default profile",
			options: []Option{
				OptionFirefox("default"),
				OptionChrome("default"),
			},
			want: false,
		},
		{
			description: "Upper dirname",
			options: []Option{
				OptionFirefox("Default"),
				OptionChrome("Default"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			NewBrowsers(tt.options...)
		})
	}
}

func TestOptionChrome(t *testing.T) {

}

func TestOptionCacheMaxAge(t *testing.T) {
	tests := []struct {
		description string
		options     []Option
		want        bool
	}{
		{
			description: "age eq 0 then cache is 24 hours",
			options: []Option{
				OptionCacheMaxAge(0),
			},
			want: false,
		},
		{
			description: "age eq -1 cache is 0 hours",
			options: []Option{
				OptionCacheMaxAge(-1),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			bookmarker := NewBrowsers(tt.options...)
			browsers := bookmarker.(*Browsers)
			if got := browsers.cache.Expired(); got != tt.want {
				t.Errorf("unexpected response \nwant: %+v\ngot: %+v", tt.want, got)
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
				testOptionFirefox(testFirefoxBookmarkJsonlz4File),
			},
			expectErr: false,
		},
		{
			description: "enable chrome bookmark",
			options: []Option{
				testOptionChrome(testChromeBookmarkJSONFile),
			},
			expectErr: false,
		},
		{
			description: "enable firefox and chrome bookmark",
			options: []Option{
				testOptionFirefox(testFirefoxBookmarkJsonlz4File),
				testOptionChrome(testChromeBookmarkJSONFile),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			bookmarker := NewBrowsers(tt.options...)
			browsers := bookmarker.(*Browsers)
			jsonData, err := browsers.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			if err := browsers.Unmarshal(jsonData); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func makeBrowserProfileDirectory(profile string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	bs := map[browser]map[string]string{
		chrome: {
			"dir":  fmt.Sprintf("%s/Library/Application Support/Google/Chrome/%s/", home, profile),
			"file": "Bookmarks",
		},
		firefox: {
			"dir":  fmt.Sprintf("%s/Library/Application Support/Firefox/Profiles/%s/bookmarkbackups", home, profile),
			"file": "test..jsonlz4",
		},
	}

	for _, b := range bs {
		dir := b["dir"]
		if !pathExists(dir) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		}
		if path := filepath.Join(dir, b["file"]); !pathExists(path) {

		}
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
