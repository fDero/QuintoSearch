/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"quinto/pkg/tokenizer"

	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Used to store documents in the database",

	PreRunE: func(cmd *cobra.Command, args []string) error {

		asInlineText, _ := cmd.Flags().GetBool("as-inline-text")
		asFilePath, _ := cmd.Flags().GetBool("as-file-path")

		if asInlineText && asFilePath {
			return fmt.Errorf("conflicting flags: --as-inline-text and --as-file-path cannot be set at the same time")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		asInlineText, _ := cmd.Flags().GetBool("as-inline-text")
		asFilePath, _ := cmd.Flags().GetBool("as-file-path")

		fmt.Println(asInlineText)
		fmt.Println(asFilePath)

		var ss []string = tokenizer.Split("Hello, World!")
		fmt.Println(ss)
		fmt.Println(args)
		fmt.Println(len(args))
		fmt.Println("store called")
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.Flags().BoolP("as-inline-text", "i", false, "Treat inputs as inline text")
	storeCmd.Flags().BoolP("as-file-path", "f", false, "Treat inputs as local file-paths")
}
