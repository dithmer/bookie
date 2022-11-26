package cmd

import (
	"github.com/dithmer/bookie/bookmarks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bookie",
	Short: "A bookmark manager",
	Long: `A bookmark manager that allows you to open bookmarks from the command line.
It is written in Go and uses a TOML config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := bookmarks.NewConfig("bookmarks.toml")
		if err != nil {
			panic(err)
		}

		err = config.OpenBookmark()
		if err != nil {
			panic(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
