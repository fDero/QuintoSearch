
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quinto",
	Short: "Minimal full-text-search engine",
	Long:  "quinto: version 0.0.0 pre-release)",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
