package cmd

import (
	"github.com/dithmer/bookie/bookmarks"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new bookmark",
	Long:  `Add a new bookmark to the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := bookmarks.NewConfig("bookmarks.toml")
		if err != nil {
			panic(err)
		}

		err = config.AddBookmark(bookmarks.Bookmark{
			Content:     content,
			Description: description,
			Type:        bType,
			Tags:        tags,
		})
		if err != nil {
			panic(err)
		}
	},
}

var (
	content     string
	description string
	tags        []string
	bType       string
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&content, "content", "c", "", "Content of the bookmark")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the bookmark")
	addCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Tags of the bookmark")
	addCmd.Flags().StringVarP(&bType, "type", "y", "", "Type of the bookmark")

	var err error

	err = addCmd.MarkFlagRequired("content")
	if err != nil {
		panic(err)
	}

	err = addCmd.MarkFlagRequired("description")
	if err != nil {
		panic(err)
	}

	err = addCmd.MarkFlagRequired("type")
	if err != nil {
		panic(err)
	}
}
