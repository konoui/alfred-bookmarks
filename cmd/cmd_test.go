package cmd

import (
	"bytes"
	"encoding/json"
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
		name      string
		filepath  string
		args      args
		config    *Config
		expectErr bool
		update    bool
	}{
		{
			name: "enbale only firefox",
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
			name: "enbale only chrome",
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
			name: "enbale only safari",
			config: &Config{
				MaxCacheAge: -1,
				Safari: Safari{
					Enable: true,
				},
			},
			filepath: filepath.Join(testdataPath, "test-safari.json"),
		},
		{
			name: "enable firefox, chrome, safari. duplicate bookmarks should be removed ",
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
			name: "pass flag format argument. no errors should occur",
			args: args{
				query: "--pass-no-match-query-as-flag-format",
			},
			config:   &Config{},
			filepath: filepath.Join(testdataPath, "empty-results.json"),
		},
		{
			name: "filter by folder name from all bookmarks. return firefox",
			args: args{
				query:  "",
				folder: "Bookmark Menu",
			},
			config: &Config{
				RemoveDuplicates: false,
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
		t.Run(tt.name, func(t *testing.T) {
			want, err := ioutil.ReadFile(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
			awf = alfred.NewWorkflow()
			awf.SetOut(outBuf)
			awf.SetLog(errBuf)
			awf.SetEmptyWarning(emptyTitle, emptySubtitle)

			r := &runtime{
				cfg:           tt.config,
				query:         tt.args.query,
				folderPrefixF: filterBySubtitle(tt.args.folder),
			}
			if err != nil {
				t.Fatal(err)
			}
			err = r.run()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error happens %+v", err.Error())
			}

			got := outBuf.Bytes()
			if diff := alfred.DiffOutput(want, got); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}

			// automatically update test data
			if tt.update {
				if err := writeFile(tt.filepath, got); err != nil {
					t.Fatal(err)
				}
			}

		})
	}
}

func Test_parse(t *testing.T) {
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
			_, err := parse(nil, tt.args.args...)
			if tt.expectedErr && err == nil {
				t.Errorf("unexpected error\n")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error %v\n", err)
			}
		})
	}
}

func writeFile(filename string, data []byte) error {
	pretty := new(bytes.Buffer)
	if err := json.Indent(pretty, data, "", "  "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, pretty.Bytes(), 0600); err != nil {
		return err
	}
	return nil
}
