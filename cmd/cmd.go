package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/konoui/alfred-bookmarks/pkg/bookmarker"
	"github.com/konoui/go-alfred"
)

var (
	awf     *alfred.Workflow
	version = "*"
)

const (
	cacheSuffix   = "alfred-bookmarks.cache"
	emptyTitle    = "No matching"
	emptySubtitle = ""
	firefoxImage  = "firefox.png"
	chromeImage   = "chrome.png"
	safariImage   = "safari.png"
)

func init() {
	awf = alfred.NewWorkflow(
		alfred.WithMaxResults(40),
		alfred.WithGitHubUpdater(
			"konoui", "alfred-bookmarks",
			version,
			14*24*time.Hour,
		),
	)
	awf.SetOut(os.Stdout)
	awf.SetLog(os.Stderr)
	awf.SetCacheSuffix(cacheSuffix)
	awf.SetEmptyWarning(emptyTitle, emptySubtitle)
}

type runtime struct {
	cfg          *Config
	query        string
	folderPrefix string
	clear        bool
}

// Execute runs cmd
func Execute(args ...string) {
	cfg, err := newConfig()
	if err != nil {
		fatal(err)
	}

	r, err := parse(cfg, args...)
	if err != nil {
		awf.Clear().Append(
			alfred.NewItem().
				Title("-f option: filster by folder name").
				Icon(alfred.IconAlertNote).
				Valid(false),
			alfred.NewItem().
				Title("--clear option: clear existing cache data").
				Icon(alfred.IconAlertNote).
				Valid(false),
		).Output()
		return
	}
	if err := r.run(); err != nil {
		fatal(err)
	}
}

func parse(cfg *Config, args ...string) (*runtime, error) {
	var folderPrefix string
	var clear bool
	fs := flag.NewFlagSet("bs", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	fs.StringVarP(&folderPrefix, "folder", "f", "", "filter by folder")
	fs.BoolVar(&clear, "clear", false, "clear cache")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	r := &runtime{
		cfg:          cfg,
		query:        alfred.Normalize(strings.Join(fs.Args(), " ")),
		folderPrefix: alfred.Normalize(folderPrefix),
		clear:        clear,
	}
	return r, nil
}

func (r *runtime) run() error {
	if err := awf.OnInitialize(); err != nil {
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if !alfred.HasUpdateArg() && awf.Updater().NewerVersionAvailable(c) {
		awf.SetSystemInfo(
			alfred.NewItem().
				Title("Newer wrokflow is available!").
				Subtitle("Please Enter!").
				Autocomplete(alfred.ArgWorkflowUpdate).
				Valid(false).
				Icon(alfred.IconAlertNote),
		)
	}

	cacheKey := "bookmarks"
	if r.clear {
		if err := awf.Cache(cacheKey).ClearItems().Err(); err != nil {
			awf.Logger().Warnln(err.Error())
		} else {
			awf.Logger().Infoln("cache cleared!")
		}
	}

	ttl := convertDefaultTTL(r.cfg.MaxCacheAge)
	if awf.Cache(cacheKey).LoadItems(ttl).Err() == nil {
		awf.Logger().Infoln("loading from cache file")
		awf.Filter(r.query).Output()
		return nil
	}

	var opts []bookmarker.Option
	if r.cfg.Firefox.Enable {
		opts = append(opts, bookmarker.OptionFirefox(r.cfg.Firefox.ProfilePath, r.cfg.Firefox.ProfileName))
	}
	if r.cfg.Chrome.Enable {
		opts = append(opts, bookmarker.OptionChrome(r.cfg.Chrome.ProfilePath, r.cfg.Chrome.ProfileName))
	}
	if r.cfg.Safari.Enable {
		opts = append(opts, bookmarker.OptionSafari())
	}

	if r.cfg.RemoveDuplicates {
		opts = append(opts, bookmarker.OptionRemoveDuplicates())
	}

	if r.folderPrefix != "" {
		// Note set empty key as to disable saving data into cache
		cacheKey = ""
		opts = append(opts, bookmarker.OptionFilterByFolder(r.folderPrefix))
	}

	manager, err := bookmarker.New(opts...)
	if err != nil {
		return err
	}

	bookmarks, err := manager.Bookmarks()
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
				).
				Variable("nextAction", "open"),
		)
	}

	defer func() {
		awf.Filter(r.query).Output()
	}()
	return awf.Cache(cacheKey).StoreItems().Err()
}

func fatal(err error) {
	awf.Fatal("a fatal error occurred", err.Error())
}
