package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/konoui/go-alfred"
	"github.com/mattn/go-shellwords"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		description string
		expectErr   bool
		filepath    string
		command     string
		errMsg      string
	}{
		{
			description: "list all bookmarks. config file exists in current directory",
			expectErr:   false,
			command:     "",
			filepath:    "test-rm-duplicate-firefox-chrome.json",
			errMsg:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			want, err := ioutil.ReadFile(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
			outStream, errStream = outBuf, errBuf
			cmdArgs, err := shellwords.Parse(tt.command)
			if err != nil {
				t.Fatalf("args parse error: %+v", err)
			}
			rootCmd := NewRootCmd()
			rootCmd.SetOutput(outStream)
			rootCmd.SetArgs(cmdArgs)

			err = rootCmd.Execute()
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
