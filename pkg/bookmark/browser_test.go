package bookmark

import (
	"testing"
	"time"
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
			description: "enable firefox, chrome, remove dupulication. return chrome bookmark",
			options: []Option{
				OptionFirefox(testFirefoxBookmarkJsonlz4File),
				OptionChrome(testChromeBookmarkJSONFile),
				OptionRemoveDuplicate(),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
		{
			description: "enable firefox, chrome, remove dupulication, cacheOption. return chrome bookmark",
			options: []Option{
				OptionFirefox(testFirefoxBookmarkJsonlz4File),
				OptionChrome(testChromeBookmarkJSONFile),
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

func TestOptionCacheMaxAge(t *testing.T) {
	tests := []struct {
		description string
		options     []Option
		want        time.Duration
	}{
		{
			description: "age eq 0 then cache is 24 hours",
			options: []Option{
				OptionCacheMaxAge(0),
			},
			want: 24 * time.Hour,
		},
		{
			description: "age eq -1 cache is 0 hours",
			options: []Option{
				OptionCacheMaxAge(-1),
			},
			want: 0 * time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			b := NewBrowsers(tt.options...)
			if b.cacheMaxAge != tt.want {
				t.Errorf("unexpected response \nwant: %+v\ngot: %+v", tt.want, b.cacheMaxAge)
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
			jsonData, err := browsers.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if err := browsers.UnmarshalJSON(jsonData); err != nil {
				t.Fatal(err)
			}
		})
	}
}
