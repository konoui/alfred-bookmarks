package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/konoui/alfred-bookmarks/pkg/bookmark"
	"github.com/konoui/go-alfred"
	"github.com/spf13/viper"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

const (
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
)

// Config configuration which browser bookmark read
type Config struct {
	Firefox struct {
		Enable bool   `json:"enable"`
		Path   string `json:"path,omitempty"`
	} `json:"firefox"`
	Chrome struct {
		Enable bool   `json:"enable"`
		Path   string `json:"path,omitempty"`
	} `json:"chrome"`
}

// NewConfig return alfred bookmark configuration
func newConfig() (*Config, error) {
	var c Config
	viper.SetConfigType("yaml")
	viper.SetConfigName(".alfred-bookmarks")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return &Config{}, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return &Config{}, err
	}

	return &c, nil
}

func run() {
	awf := alfred.NewWorkflow()
	awf.SetStdStream(outStream)
	awf.SetErrStream(outStream)
	awf.EmptyWarning(emptyTitle, emptySsubtitle)

	var query string
	if args := os.Args; len(args) > 1 {
		query = args[1]
		log.Printf("query: %s", query)
	}

	c, err := newConfig()
	if err != nil {
		awf.Fatal(fmt.Sprintf("a error occurs: %s", err), "")
		return

	}

	foption, coption := bookmark.OptionNone(), bookmark.OptionNone()
	if c.Firefox.Enable {
		foption = bookmark.OptionFirefox(c.Firefox.Path)
	}

	if c.Chrome.Enable {
		coption = bookmark.OptionChrome(c.Chrome.Path)
	}

	browsers := bookmark.NewBrowsers(
		foption,
		coption,
	)

	bookmarks, err := browsers.Bookmarks()
	if err != nil {
		awf.Fatal(fmt.Sprintf("an error occurs: %s", err), "")
		return
	}

	log.Printf("%d total bookmark(s)", len(bookmarks))

	if query != "" {
		bookmarks = bookmarks.Filter(query)
		log.Printf("%d total filtered bookmark(s)", len(bookmarks))
	}

	for _, b := range bookmarks {
		subtitle := fmt.Sprintf("[%s] %s", b.Folder, b.Domain)
		awf.Append(alfred.Item{
			Title:        b.Title,
			Subtitle:     subtitle,
			Autocomplete: b.Title,
			Arg:          b.URI,
		})
	}

	awf.Output()
}

func main() {
	log.SetOutput(ioutil.Discard)
	run()
}
