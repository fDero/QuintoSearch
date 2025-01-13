package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tokenizeCmd = &cobra.Command{
	Use:   "tokenize",
	Short: "Used to inspect the tokens generated from the input",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateInputFlags(cmd, args)
	},

	Run: func(cmd *cobra.Command, args []string) {
		tokens := ParseInputTokens(cmd, args)
		for _, token := range tokens {
			fmt.Printf("[%s] ", token)
		}
	},
}

func init() {
	rootCmd.AddCommand(tokenizeCmd)
	RegisterInputFlags(tokenizeCmd)
}
