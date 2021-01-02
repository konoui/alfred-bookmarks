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
	type args struct {
		query  string
		folder string
	}
	tests := []struct {
		description string
		expectErr   bool
		filepath    string
		args        args
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
				RemoveDuplicates: true,
				MaxCacheAge:      -1,
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
			args: args{
				query: "--pass-no-match-query-as-flag-format",
			},
			config:   &Config{},
			filepath: filepath.Join(testdataPath, "empty-results.json"),
		},
		{
			description: "fildter by folder name. return firefox when RemoveDuplicate is true",
			args: args{
				query:  "",
				folder: "Bookmark Menu",
			},
			config: &Config{
				RemoveDuplicates: true,
				MaxCacheAge:      -1,
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
			filepath: filepath.Join(testdataPath, "test-firefox.json"),
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
			awf.SetLog(errBuf)
			awf.SetEmptyWarning(emptyTitle, emptySubtitle)

			err = tt.config.run(tt.args.query, tt.args.folder)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error happens %+v", err.Error())
			}

			got := outBuf.Bytes()
			if diff := alfred.DiffOutput(want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func Test_parseQuery(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name: "flag parse",
			args: args{
				[]string{
					"-f",
					"Bookmark menu",
				},
			},
		},
		{
			name: "command aflter flag parse",
			args: args{
				[]string{
					"github",
					"-f",
					"Bookmark menu",
				},
			},
		},
		{
			name: "flag parse",
			args: args{
				[]string{
					"-f",
				},
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseQuery(tt.args.args...)
			if tt.expectedErr && err == nil {
				t.Errorf("unexpected error\n")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error %v\n", err)
			}
		})
	}
}
