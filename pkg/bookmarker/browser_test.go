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
			if err := browsers.(*Browsers).cache.Clear(); err != nil {
				t.Fatal(err)
			}

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
