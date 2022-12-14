package cmd

import (
	"log"
	"os"

	"github.com/dithmer/bookie/bookmarks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bookie",
	Short: "A bookmark manager",
	Long: `A bookmark manager that allows you to open bookmarks from the command line.
It is written in Go and uses a TOML config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := bookmarks.NewConfig(configPath)
		if err != nil {
			log.Fatal("Error while reading config from ", configPath, ": ", err)
		}

		if query != "" {
			err = config.OpenBookmarkWithQuery(bookmarks.ParseQuery(query), openWithParam, openWithFile)
		} else {
			err = config.OpenBookmark(openWithParam, openWithFile)
		}
		if err != nil {
			log.Fatal("Error while opening bookmark:", err)
		}
	},
}

var configPath string
var query string
var openWithParam string
var openWithFile string

func init() {
	standardConfigPath := os.Getenv("HOME") + "/.config/bookie/config.toml"

	rootCmd.PersistentFlags().StringVarP(&configPath, "config-path", "p", standardConfigPath, "Path to the config file")
	rootCmd.PersistentFlags().StringVarP(&openWithParam, "open-with-type", "o", "", "Open bookmark with a specific application")
	rootCmd.PersistentFlags().StringVarP(&openWithFile, "open-with-file", "f", "", "Open bookmark with a specific file")

	rootCmd.Flags().StringVarP(&query, "query", "q", "", "Query to search for bookmarks")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Error while executing root command: ", err)
	}
}
