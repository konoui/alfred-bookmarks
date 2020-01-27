package bookmarker

import (
	"testing"
)

var testProfile = "default"

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
				OptionFirefox(testProfile),
			},
			want:      testFirefoxBookmarks,
			expectErr: false,
		},
		{
			description: "enable chrome bookmark",
			options: []Option{
				OptionChrome(testProfile),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
		{
			description: "enable firefox, chrome, remove dupulication. return chrome bookmark",
			options: []Option{
				OptionFirefox(testProfile),
				OptionChrome(testProfile),
				OptionRemoveDuplicate(),
			},
			want:      testChromeBookmarks,
			expectErr: false,
		},
		{
			description: "enable firefox, chrome, remove dupulication, cacheOption. return chrome bookmark",
			options: []Option{
				OptionFirefox(testProfile),
				OptionChrome(testProfile),
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
			if err := browsers.(*Browsers).cacher.Clear(); err != nil {
				t.Fatal(err)
			}

			bookmarks, err := browsers.Bookmarks()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}

			diff := DiffBookmark(bookmarks, tt.want)
			if diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func TestOptionFirefoxChrome(t *testing.T) {
	tests := []struct {
		description string
		options     []Option
	}{
		{
			description: "Lower default profile name",
			options: []Option{
				OptionFirefox("default"),
				OptionChrome("default"),
			},
		},
		{
			description: "Upper default profile name",
			options: []Option{
				OptionFirefox("Default"),
				OptionChrome("Default"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// panic if error occurs
			NewBrowsers(tt.options...)
		})
	}
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
			if got := browsers.cacher.Expired(); got != tt.want {
				t.Errorf("want: %+v\n, got: %+v\n", tt.want, got)
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
				OptionFirefox(testProfile),
			},
			expectErr: false,
		},
		{
			description: "enable chrome bookmark",
			options: []Option{
				OptionChrome(testProfile),
			},
			expectErr: false,
		},
		{
			description: "enable firefox and chrome bookmark",
			options: []Option{
				OptionFirefox(testProfile),
				OptionChrome(testProfile),
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
