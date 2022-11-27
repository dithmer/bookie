package cmd

import (
	"bufio"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
)

var chooserCmd = &cobra.Command{
	Use:   "chooser",
	Short: "Run inbuilt chooser",
	Long:  `Run inbuilt chooser`,
	Run: func(cmd *cobra.Command, args []string) {
		var bookmarks []string
		var incomingBookmarksChan = make(chan string)

		// read stdin line by line
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				incomingBookmarksChan <- scanner.Text()
			}
		}()

		a := app.New()

		var width, height float32 = 800, 200
		w := a.NewWindow("Bookie Chooser")
		w.Resize(fyne.NewSize(width, height))

		go func() {
			for {
				bookmark := <-incomingBookmarksChan
				bookmarks = append(bookmarks, bookmark)

				// update UI
				list := widget.NewList(
					func() int {
						return len(bookmarks)
					},
					func() fyne.CanvasObject {
						return widget.NewLabel("")
					},
					func(id widget.ListItemID, item fyne.CanvasObject) {
						item.(*widget.Label).SetText(bookmarks[id])
					},
				)

				w.SetContent(list)
			}
		}()

		w.Show()
		a.Run()
	},
}

func init() {
	rootCmd.AddCommand(chooserCmd)
}
