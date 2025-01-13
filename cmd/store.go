package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Used to store documents in the database",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateInputFlags(cmd, args)
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("store called")
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	RegisterInputFlags(storeCmd)
}
