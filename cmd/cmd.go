package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/konoui/alfred-bookmarks/pkg/bookmarker"
	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/initialize"
)

var (
	awf     *alfred.Workflow
	version = "*"
)

const (
	emptyTitle    = "No matching"
	emptySubtitle = ""
	firefoxImage  = "firefox.png"
	chromeImage   = "chrome.png"
	safariImage   = "safari.png"
)

func init() {
	awf = alfred.NewWorkflow(
		alfred.WithMaxResults(20),
		alfred.WithGitHubUpdater(
			"konoui", "alfred-bookmarks",
			version,
			14*24*time.Hour,
		),
		alfred.WithCacheSuffix("-alfred-bookmarks.cache"),
		alfred.WithOutWriter(os.Stdout),
		alfred.WithLogWriter(os.Stderr),
		alfred.WithInitializers(
			initialize.NewUpdateRecommendation(5*time.Second),
			initialize.NewUpdateExecution(2*time.Minute),
		),
	)
	awf.SetEmptyWarning(emptyTitle, emptySubtitle)
}

type runtime struct {
	cfg           *Config
	query         string
	folderPrefixF func(subtitle string) bool
	clear         bool
}

// Execute runs cmd
func Execute(args ...string) {
	cfg, err := newConfig()
	if err != nil {
		awf.Fatal("a fatal error occurred", err.Error())
	}

	r, err := parse(cfg, args...)
	if err != nil {
		awf.Clear().Append(
			alfred.NewItem().
				Title("-f option: filter by bookmark folder name").
				Icon(awf.Assets().IconAlertNote()).
				Valid(false),
			alfred.NewItem().
				Title("--clear option: clear existing cache data").
				Icon(awf.Assets().IconAlertNote()).
				Valid(false),
		).Output()
		return
	}
	os.Exit(awf.RunSimple(r.run))
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
		cfg:           cfg,
		query:         strings.Join(fs.Args(), " "),
		folderPrefixF: filterBySubtitle(folderPrefix),
		clear:         clear,
	}
	return r, nil
}

func (r *runtime) run() error {
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
		awf.FilterByItemProperty(r.folderPrefixF, alfred.ItemPropertySubtitle).
			Filter(r.query).Output()
		return nil
	}

	opts := make([]bookmarker.Option, 0, 4)
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
		awf.FilterByItemProperty(r.folderPrefixF, alfred.ItemPropertySubtitle).
			Filter(r.query).Output()
	}()
	return awf.Cache(cacheKey).StoreItems().Err()
}

func filterBySubtitle(prefixQuery string) func(subtitle string) bool {
	f := func(subtitle string) bool {
		// Note: if input is empty return true
		if prefixQuery == "" {
			return true
		}

		leftTrimedSubtitle := strings.TrimLeft(subtitle, "[")
		idx := strings.LastIndex(leftTrimedSubtitle, "]")
		if idx < 0 {
			return false
		}
		folder := leftTrimedSubtitle[:idx]
		return hasFolderPrefix(folder, prefixQuery)
	}
	return f
}

func hasFolderPrefix(folder, prefix string) bool {
	folder = strings.ToLower(folder)
	folder = strings.ReplaceAll(folder, " ", "")
	prefix = strings.ToLower(prefix)
	prefix = strings.ReplaceAll(prefix, " ", "")
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	if strings.HasPrefix(folder, prefix) {
		return true
	}

	return strings.HasPrefix(folder+"/", prefix)
}
