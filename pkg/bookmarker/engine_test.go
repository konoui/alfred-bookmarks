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
			description: "enable safari bookmark",
			options: []Option{
				OptionSafari(),
			},
			want:      testSafariBookmarks,
			expectErr: false,
		},
		{
			description: "enable firefox, chrome, safari and remove dupulication. return chrome bookmark",
			options: []Option{
				OptionFirefox(testProfile),
				OptionChrome(testProfile),
				OptionSafari(),
				OptionRemoveDuplicate(),
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
