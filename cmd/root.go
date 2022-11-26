package cmd

import (
	"log"

	"github.com/dithmer/bookie/bookmarks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bookie",
	Short: "A bookmark manager",
	Long: `A bookmark manager that allows you to open bookmarks from the command line.
It is written in Go and uses a TOML config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := "bookmarks.toml"
		config, err := bookmarks.NewConfig(configPath)
		if err != nil {
			log.Fatal("Error while reading config from", configPath, ":", err)
		}

		err = config.OpenBookmark()
		if err != nil {
			log.Fatal("Error while opening bookmark:", err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Error while executing root command:", err)
	}
}
