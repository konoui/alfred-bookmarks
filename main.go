package main

import (
	"fmt"
	"log"
	"os"

	aw "github.com/deanishe/awgo"
	"github.com/konoui/alfred-firefox-bookmarks/bookmark"
)

var (
	wf *aw.Workflow
)

const (
	emptySsubtitle = "There are no resources"
	emptyTitle     = "No matching"
)

func run() {
	var query string
	if args := os.Args; len(args) > 1 {
		query = args[1]
		log.Printf("query: %s", query)
	}

	bookmarks, err := bookmark.LoadBookmarksFromCache()
	if err != nil {
		wf.FatalError(err)
	}

	log.Printf("%d total bookmark(s)", len(bookmarks))

	if query != "" {
		bookmarks = bookmarks.Filter(query)
		log.Printf("%d total filtered bookmark(s)", len(bookmarks))
	}

	for _, b := range bookmarks {
		title := fmt.Sprintf("%s/%s", b.Folder, b.Title)
		wf.NewItem(title).
			Subtitle(b.Domain).
			Arg(b.URI).
			Autocomplete(title).
			Valid(true)
	}

	wf.WarnEmpty(emptyTitle, emptySsubtitle)
	wf.SendFeedback()
}

func main() {
	wf = aw.New()

	const debugLogFile = "alfred-firefox-bookmarks.log"
	f, err := os.OpenFile(debugLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	wf.Run(run)
}
