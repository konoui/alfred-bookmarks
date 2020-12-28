package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/konoui/alfred-bookmarks/pkg/bookmarker"
	"github.com/konoui/go-alfred"
)

var (
	awf *alfred.Workflow
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
	)
	awf.SetOut(os.Stdout)
	awf.SetLog(os.Stderr)
	awf.SetCacheSuffix(cacheSuffix)
	awf.SetEmptyWarning(emptyTitle, emptySubtitle)
}

// Execute runs cmd
func Execute(args ...string) {
	c, err := newConfig()
	if err != nil {
		fatal(err)
	}

	query, folderPrefix, err := parseQuery(args...)
	if err != nil {
		awf.Clear().
			SetEmptyWarning("-f option: filster by folder name", err.Error()).
			Output()
		return
	}
	if err := c.run(query, folderPrefix); err != nil {
		fatal(err)
	}
}

func parseQuery(args ...string) (query, folderPrefix string, err error) {
	fs := flag.NewFlagSet("bs", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	fs.StringVarP(&folderPrefix, "folder", "f", "", "filter by folder")
	if err := fs.Parse(args); err != nil {
		return "", "", err
	}
	query = fs.Arg(0)
	return alfred.Normalize(query), alfred.Normalize(folderPrefix), nil
}

func (c *Config) run(query, folderPrefix string) error {
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
		opts = append(opts, bookmarker.OptionRemoveDuplicates())
	}

	if folderPrefix != "" {
		opts = append(opts, bookmarker.OptionFilterByFolder(folderPrefix))
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

	return awf.Filter(query).Output().Cache(cacheKey).StoreItems().Err()
}

func fatal(err error) {
	awf.Fatal("Fatal errors occur", err.Error())
}
