package cmd

import (
	"fmt"
	"log"

	"github.com/dithmer/bookie/bookmarks"
	"github.com/spf13/cobra"
)

var listTagsCmd = &cobra.Command{
	Use:   "list-tags",
	Short: "List all tags",
	Long:  `List all tags`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := bookmarks.NewConfig(configPath)
		if err != nil {
			log.Fatal("Error loading config: ", err)
		}

		tags, err := config.ListTags()
		if err != nil {
			log.Fatal("Error listing tags: ", err)
		}

		for _, tag := range tags {
			fmt.Println(tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(listTagsCmd)
}
