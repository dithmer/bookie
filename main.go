package main

import (
	"github.com/dithmer/bookie/bookmarks"
)

func main() {
	config, err := bookmarks.NewConfig()
	handleMainError(err)

	err = config.OpenBookmark()
	handleMainError(err)
}

func handleMainError(err error) {
	if err != nil {
		panic(err)
	}
}
