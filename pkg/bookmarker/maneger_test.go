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
				OptionFirefox(defaultFirefoxProfilePath, testProfile),
			},
			want: testFirefoxBookmarks,
		},
		{
			description: "enable chrome bookmark",
			options: []Option{
				OptionChrome(defaultChromeProfilePath, testProfile),
			},
			want: testChromeBookmarks,
		},
		{
			description: "enable safari bookmark",
			options: []Option{
				OptionSafari(),
			},
			want: testSafariBookmarks,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			m, err := New(tt.options...)
			if err != nil {
				t.Fatal(err)
			}

			bookmarks, err := m.Bookmarks()
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
				OptionFirefox(defaultFirefoxProfilePath, "default"),
				OptionChrome(defaultChromeProfilePath, "default"),
			},
		},
		{
			description: "Upper default profile name",
			options: []Option{
				OptionFirefox(defaultFirefoxProfilePath, "Default"),
				OptionChrome(defaultChromeProfilePath, "Default"),
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
