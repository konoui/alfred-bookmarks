package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

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
			description: "list all bookmarks. config file which enbale all setting exists in current directory",
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
				t.Errorf("unexpected error want: got: %+v", err.Error())
			}

			got := outBuf.String()
			if !EqualJSON(want, []byte(got)) {
				t.Errorf("unexpected response: want: \n%+v, got: \n%+v", string(want), string(got))
			}
		})
	}
}

func EqualJSON(a, b []byte) bool {
	var ao interface{}
	var bo interface{}

	if err := json.Unmarshal(a, &ao); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bo); err != nil {
		return false
	}

	return reflect.DeepEqual(ao, bo)
}
