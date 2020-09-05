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

const (
	cacheSuffix    = "alfred-bookmarks.cache"
	cacheKey       = "bookmarks"
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
	fatalError     = "Fatal errors occur"
	firefoxImage   = "firefox.png"
	chromeImage    = "chrome.png"
	safariImage    = "safari.png"
)

func init() {
	awf = alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	awf.SetCacheSuffix(cacheSuffix)
	if err := awf.SetCacheDir(os.TempDir()); err != nil {
		awf.Fatal(fatalError, err.Error())
	}
	awf.EmptyWarning(emptyTitle, emptySsubtitle)
}

// Execute runs cmd
func Execute(args ...string) {
	c, err := newConfig()
	if err != nil {
		awf.Fatal(fatalError, err.Error())
	}

	if err := c.run(strings.Join(args, " ")); err != nil {
		awf.Fatal(fatalError, err.Error())
	}
}

func (c *Config) run(query string) error {
	ttl := convertDefaultTTL(c.MaxCacheAge)
	if awf.Cache(cacheKey).MaxAge(ttl).LoadItems().Err() == nil {
		awf.Filter(query).Output()
		return nil
	}

	var opts []bookmarker.Option
	if c.Firefox.Enable {
		opts = append(opts, bookmarker.OptionFirefox(c.Firefox.ProfilePath, c.Firefox.ProfileName))
	}
	if c.Chrome.Enable {
		opts = append(opts, bookmarker.OptionChrome(c.Chrome.ProfilePath, c.Chrome.ProfileName))
	}
	if c.Safari.Enable {
		opts = append(opts, bookmarker.OptionSafari())
	}
	if c.RemoveDuplicate {
		opts = append(opts, bookmarker.OptionRemoveDuplicate())
	}

	engine, err := bookmarker.New(opts...)
	if err != nil {
		return err
	}

	bookmarks, err := engine.Bookmarks()
	if err != nil {
		return err
	}

	for _, b := range bookmarks {
		var image string
		switch b.BookmarkerName {
		case bookmarker.Firefox:
			image = firefoxImage
		case bookmarker.Chrome:
			image = chromeImage
		case bookmarker.Safari:
			image = safariImage
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
