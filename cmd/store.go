
package cmd

import (
	"fmt"
	"quinto/frontend"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Used to store documents in the database",

	PreRunE: func(cmd *cobra.Command, args []string) error {

		asInlineText, _ := cmd.Flags().GetString("inline")
		asFilePath, _ := cmd.Flags().GetString("filepath")

		if len(asInlineText) > 0 && len(asFilePath) > 0 {
			return fmt.Errorf("conflicting flags: --as-inline-text and --as-file-path cannot be set at the same time")
		}

		if len(asInlineText) + len(asFilePath) == 0 {
			return fmt.Errorf("missing flags: --as-inline-text or --as-file-path must be set")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		asInlineText, _ := cmd.Flags().GetString("inline")
		asFilePath, _ := cmd.Flags().GetString("filepath")
		
		var tokens []string
		if (len(asInlineText) > 0) {
			tokens = frontend.ProcessInputText(asInlineText)
		} else if (len(asFilePath) > 0) {
			tokens = frontend.ProcessInputFile(asFilePath)
		}

		fmt.Println(tokens)
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.Flags().String("inline", "", "Treat inputs as inline text")
	storeCmd.Flags().String("filepath", "", "Treat inputs as local file-paths")
	storeCmd.Flags().String("lang", "eng", "Select language: eng->English")
}
