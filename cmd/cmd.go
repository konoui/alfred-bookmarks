package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/konoui/alfred-bookmarks/pkg/bookmarker"
	"github.com/konoui/go-alfred"
)

var (
	awf       *alfred.Workflow
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

func init() {
	awf = alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	awf.EmptyWarning(emptyTitle, emptySsubtitle)
}

// Execute runs cmd
func Execute(args ...string) {
	if err := run(strings.Join(args, " ")); err != nil {
		awf.Fatal("fatal error occurs", err.Error())
	}
}

const (
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
	firefoxImage   = "firefox.png"
	chromeImage    = "chrome.png"
)

func run(query string) error {
	c, err := newConfig()
	if err != nil {
		return err
	}

	firefoxOption, chromeOption := bookmarker.OptionNone(), bookmarker.OptionNone()
	duplicateOption := bookmarker.OptionNone()
	if c.Firefox.Enable {
		firefoxOption = bookmarker.OptionFirefox(c.Firefox.Profile)
	}
	if c.Chrome.Enable {
		chromeOption = bookmarker.OptionChrome(c.Chrome.Profile)
	}
	if c.RemoveDuplicate {
		duplicateOption = bookmarker.OptionRemoveDuplicate()
	}

	engine, err := bookmarker.New(
		firefoxOption,
		chromeOption,
		duplicateOption,
		bookmarker.OptionCacheMaxAge(c.MaxCacheAge),
	)
	if err != nil {
		return err
	}

	bookmarks, err := engine.Bookmarks()
	if err != nil {
		return err
	}

	if query != "" {
		bookmarks = bookmarks.Filter(query)
	}

	for _, b := range bookmarks {
		var image string
		if b.BookmarkerName == bookmarker.Firefox {
			image = firefoxImage
		} else {
			image = chromeImage
		}
		awf.Append(&alfred.Item{
			Title:        b.Title,
			Subtitle:     fmt.Sprintf("[%s] %s", b.Folder, b.Domain),
			Autocomplete: b.Title,
			Arg:          b.URI,
			Icon: &alfred.Icon{
				Path: image,
			},
		})
	}

	awf.Output()
	return nil
}
