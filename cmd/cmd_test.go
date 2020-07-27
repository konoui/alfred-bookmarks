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
			description: "list all bookmarks. config file exists in current directory",
			expectErr:   false,
			command:     "",
			config:      testConfig,
			filepath:    filepath.Join(testdataPath, "test-rm-duplicate-firefox-chrome.json"),
		},
		{
			description: "flag format argument. no error occurs",
			expectErr:   false,
			command:     "--pass-no-match-query-as-flag-format",
			config:      testConfig,
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
			awf.EmptyWarning(emptyTitle, emptySsubtitle)

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
