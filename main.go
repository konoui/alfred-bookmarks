package main

import (
	"os"

	"github.com/konoui/alfred-bookmarks/cmd"
)

func main() {
	cmd.Execute(os.Args[1:]...)
}
