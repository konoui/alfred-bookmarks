package main

import (
	"github.com/konoui/alfred-bookmarks/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	cmd.Execute(rootCmd)
}
