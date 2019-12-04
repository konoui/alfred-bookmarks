package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/konoui/alfred-bookmarks/pkg/bookmark"
	"github.com/konoui/go-alfred"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

const (
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
)

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

	bookmarks, err := bookmark.LoadBookmarksFromCache()
	if err != nil {
		awf.Fatal(fmt.Sprintf("A Error Occurs: %s", err), "")
		return
	}

	log.Printf("%d total bookmark(s)", len(bookmarks))

	if query != "" {
		bookmarks = bookmarks.Filter(query)
		log.Printf("%d total filtered bookmark(s)", len(bookmarks))
	}

	for _, b := range bookmarks {
		subTitle := fmt.Sprintf("[%s] %s", b.Folder, b.Domain)
		awf.Append(alfred.Item{
			Title:        b.Title,
			Subtitle:     subTitle,
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
