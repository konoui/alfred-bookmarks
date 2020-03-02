package bookmarker

import (
	"testing"
)

var testProfile = "default"

func TestEngineBookmarks(t *testing.T) {
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
			e, err := New(tt.options...)
			if err != nil {
				t.Fatal(err)
			}

			if err = e.(*engine).cacher.Clear(); err != nil {
				t.Fatal(err)
			}

			bookmarks, err := e.Bookmarks()
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
			if _, err := New(tt.options...); err != nil {
				t.Error(err)
			}
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
			bookmarker, err := New(tt.options...)
			if err != nil {
				t.Fatal(err)
			}

			e := bookmarker.(*engine)
			if got := e.cacher.Expired(); got != tt.want {
				t.Errorf("want: %+v\n, got: %+v\n", tt.want, got)
			}
		})
	}
}

func TestEngineMarshaUnmarshalJson(t *testing.T) {
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
			bookmarker, err := New(tt.options...)
			if err != nil {
				t.Fatal(err)
			}

			engine := bookmarker.(*engine)
			jsonData, err := engine.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			if err := engine.Unmarshal(jsonData); err != nil {
				t.Fatal(err)
			}
		})
	}
}
