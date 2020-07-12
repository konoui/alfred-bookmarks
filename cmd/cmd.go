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
	cacheDir            = os.TempDir()
)

const (
	cacheSuffix    = "alfred-bookmarks.cache"
	cacheKey       = "bookmarks"
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
	fatalError     = "Fatal errors occur"
	firefoxImage   = "firefox.png"
	chromeImage    = "chrome.png"
)

func init() {
	awf = alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	awf.SetCacheSuffix(cacheSuffix)
	if err := awf.SetCacheDir(cacheDir); err != nil {
		awf.Fatal(fatalError, err.Error())
	}
	awf.EmptyWarning(emptyTitle, emptySsubtitle)
}

// Execute runs cmd
func Execute(args ...string) {
	if err := run(strings.Join(args, " ")); err != nil {
		awf.Fatal(fatalError, err.Error())
	}
}

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

	ttl := convertDefaultTTL(c.MaxCacheAge)
	if awf.Cache(cacheKey).MaxAge(ttl).LoadItems().Err() == nil {
		awf.Filter(query).Output()
		return nil
	}

	engine, err := bookmarker.New(
		firefoxOption,
		chromeOption,
		duplicateOption,
	)
	if err != nil {
		return err
	}

	bookmarks, err := engine.Bookmarks()
	if err != nil {
		return err
	}

	for _, b := range bookmarks {
		var image string
		if b.BookmarkerName == bookmarker.Firefox {
			image = firefoxImage
		} else {
			image = chromeImage
		}
		awf.Append(
			alfred.NewItem().
				SetTitle(b.Title).
				SetSubtitle(fmt.Sprintf("[%s] %s", b.Folder, b.Domain)).
				SetAutocomplete(b.Title).
				SetArg(b.URI).
				SetIcon(
					alfred.NewIcon().
						SetPath(image),
				),
		)
	}

	awf.Cache(cacheKey).StoreItems().Workflow().Filter(query).Output()
	return nil
}
