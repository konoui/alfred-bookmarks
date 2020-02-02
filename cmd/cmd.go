package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/konoui/alfred-bookmarks/pkg/bookmarker"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

// Execute Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetOutput(outStream)
	if err := rootCmd.Execute(); err != nil {
		log.Printf("command execution failed: %+v", err)
		os.Exit(1)
	}
}

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "alfred-bookmarks <query>",
		Short: "search bookmarks",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			return run(query)
		},
		SilenceUsage: true,
	}

	return rootCmd
}

const (
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
	firefoxImage   = "firefox.png"
	chromeImage    = "chrome.png"
)

func run(query string) error {
	awf := alfred.NewWorkflow()
	awf.SetStreams(outStream, errStream)
	awf.EmptyWarning(emptyTitle, emptySsubtitle)

	c, err := newConfig()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
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

	browsers := bookmarker.New(
		firefoxOption,
		chromeOption,
		duplicateOption,
		bookmarker.OptionCacheMaxAge(c.MaxCacheAge),
	)

	bookmarks, err := browsers.Bookmarks()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		return err
	}

	log.Printf("%d total bookmark(s)", len(bookmarks))
	log.Printf("query: %s", query)
	if query != "" {
		bookmarks = bookmarks.Filter(query)
		log.Printf("%d total filtered bookmark(s)", len(bookmarks))
	}

	for _, b := range bookmarks {
		var image string
		if b.Bookmarker == bookmarker.Firefox {
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
