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
	cacheSuffix   = "alfred-bookmarks.cache"
	cacheKey      = "bookmarks"
	emptyTitle    = "No matching"
	emptySubtitle = ""
	firefoxImage  = "firefox.png"
	chromeImage   = "chrome.png"
	safariImage   = "safari.png"
)

func init() {
	awf = alfred.NewWorkflow(
		alfred.WithMaxResults(40),
		alfred.WithLogStream(outStream),
		alfred.WithLogStream(errStream),
	)
	awf.SetCacheSuffix(cacheSuffix)
	awf.SetEmptyWarning(emptyTitle, emptySubtitle)
}

// Execute runs cmd
func Execute(args ...string) {
	c, err := newConfig()
	if err != nil {
		fatal(err)
	}

	if err := c.run(strings.Join(args, " ")); err != nil {
		fatal(err)
	}
}

func (c *Config) run(query string) error {
	ttl := convertDefaultTTL(c.MaxCacheAge)
	if awf.Cache(cacheKey).LoadItems(ttl).Err() == nil {
		awf.Logger().Infoln("loading from cache file")
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
				Title(b.Title).
				Subtitle(fmt.Sprintf("[%s] %s", b.Folder, b.Domain)).
				Autocomplete(b.Title).
				Arg(b.URI).
				Icon(
					alfred.NewIcon().
						Path(image),
				),
		).Variable("nextAction", "open")
	}

	awf.Cache(cacheKey).StoreItems().Workflow().Filter(query).Output()
	return nil
}

func fatal(err error) {
	awf.Fatal("Fatal errors occur", err.Error())
}
