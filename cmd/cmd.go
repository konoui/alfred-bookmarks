package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/konoui/alfred-bookmarks/pkg/bookmark"
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
			if err := run(query); err != nil {
				log.Printf("run failed: %+v", err)
			}
			return nil
		},
		SilenceUsage: true,
	}

	return rootCmd
}

const (
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
)

func run(query string) error {
	awf := alfred.NewWorkflow()
	awf.SetStdStream(outStream)
	awf.SetErrStream(outStream)
	awf.EmptyWarning(emptyTitle, emptySsubtitle)

	c, err := newConfig()
	if err != nil {
		awf.Fatal(fmt.Sprintf("an error occurs: %s", err), "")
		return err
	}

	foption, coption := bookmark.OptionNone(), bookmark.OptionNone()
	duplicate := bookmark.OptionNone()
	if c.Firefox.Enable {
		foption = bookmark.OptionFirefox(c.Firefox.Path)
	}
	if c.Chrome.Enable {
		coption = bookmark.OptionChrome(c.Chrome.Path)
	}
	if c.RemoveDuplicate {
		duplicate = bookmark.OptionRemoveDuplicate()
	}

	browsers := bookmark.NewBrowsers(
		foption,
		coption,
		duplicate,
	)

	bookmarks, err := browsers.Bookmarks()
	if err != nil {
		awf.Fatal(fmt.Sprintf("an error occurs: %s", err), "")
		return err
	}

	log.Printf("%d total bookmark(s)", len(bookmarks))
	log.Printf("query: %s", query)
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
	return nil
}
