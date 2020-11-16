package cmd

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/konoui/go-alfred"
)

const testdataPath = "testdata"

func TestRun(t *testing.T) {
	tests := []struct {
		description string
		expectErr   bool
		filepath    string
		command     string
		config      *Config
		errMsg      string
	}{
		{
			description: "enbale only firefox",
			config: &Config{
				MaxCacheAge: -1,
				Firefox: Firefox{
					Enable:      true,
					ProfileName: firefoxDefaultProfileName,
					ProfilePath: firefoxDefaultProfilePath,
				},
			},
			filepath: filepath.Join(testdataPath, "test-firefox.json"),
		},
		{
			description: "enbale only chrome",
			config: &Config{
				MaxCacheAge: -1,
				Chrome: Chrome{
					Enable:      true,
					ProfileName: chromeDefaultProfileName,
					ProfilePath: chromeDefaultProfilePath,
				},
			},
			filepath: filepath.Join(testdataPath, "test-chrome.json"),
		},
		{
			description: "enbale only safari",
			config: &Config{
				MaxCacheAge: -1,
				Safari: Safari{
					Enable: true,
				},
			},
			filepath: filepath.Join(testdataPath, "test-safari.json"),
		},
		{
			description: "enable firefox, chrome, safari. duplicate bookmarks should be removed ",
			config: &Config{
				RemoveDuplicate: true,
				MaxCacheAge:     -1,
				Firefox: Firefox{
					Enable:      true,
					ProfileName: firefoxDefaultProfileName,
					ProfilePath: firefoxDefaultProfilePath,
				},
				Chrome: Chrome{
					Enable:      true,
					ProfileName: chromeDefaultProfileName,
					ProfilePath: chromeDefaultProfilePath,
				},
				Safari: Safari{
					Enable: true,
				},
			},
			filepath: filepath.Join(testdataPath, "test-rm-duplicate-firefox-chrome-safari.json"),
		},
		{
			description: "pass flag format argument. no errors should occur",
			command:     "--pass-no-match-query-as-flag-format",
			config:      &Config{},
			filepath:    filepath.Join(testdataPath, "empty-results.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			want, err := ioutil.ReadFile(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
			awf = alfred.NewWorkflow()
			awf.SetOut(outBuf)
			awf.SetErr(errBuf)
			awf.SetEmptyWarning(emptyTitle, emptySubtitle)

			err = tt.config.run(tt.command)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error happens %+v", err.Error())
			}

			got := outBuf.Bytes()
			if diff := alfred.DiffScriptFilter(want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
