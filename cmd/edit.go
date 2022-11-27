package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit bookie config",
	Long:  `Edit bookie config`,
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")

		editorCmd := exec.Command(editor, configPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout

		err := editorCmd.Run()
		if err != nil {
			log.Fatal("Error editing config file: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
