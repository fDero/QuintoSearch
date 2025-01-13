package cmd

import (
	"fmt"
	"quinto/frontend"

	"github.com/spf13/cobra"
)

func ValidateInputFlags(cmd *cobra.Command, args []string) error {
	asInlineText, _ := cmd.Flags().GetString("inline")
	asFilePath, _ := cmd.Flags().GetString("filepath")

	if len(asInlineText) > 0 && len(asFilePath) > 0 {
		return fmt.Errorf("conflicting flags: --as-inline-text and --as-file-path cannot be set at the same time")
	}

	if len(asInlineText)+len(asFilePath) == 0 {
		return fmt.Errorf("missing flags: --as-inline-text or --as-file-path must be set")
	}

	return nil
}

func RegisterInputFlags(cmd *cobra.Command) {
	cmd.Flags().String("inline", "", "Treat inputs as inline text")
	cmd.Flags().String("filepath", "", "Treat inputs as local file-paths")
	cmd.Flags().String("lang", "eng", "Select language: eng->English")
}

func ParseInputTokens(cmd *cobra.Command, args []string) []string {
	asInlineText, _ := cmd.Flags().GetString("inline")
	asFilePath, _ := cmd.Flags().GetString("filepath")

	if len(asInlineText) > 0 {
		return frontend.ProcessInputText(asInlineText)
	} else {
		return frontend.ProcessInputFile(asFilePath)
	}
}
