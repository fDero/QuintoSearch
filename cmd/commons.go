package cmd

import (
	"fmt"
	"iter"
	"log"
	"os"
	"quinto/core"
	"quinto/data"
	"quinto/stemming"

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

func IterateTokens(cmd *cobra.Command, args []string) iter.Seq[core.Token] {
	asInlineText, _ := cmd.Flags().GetString("inline")
	asFilePath, _ := cmd.Flags().GetString("filepath")
	lang, _ := cmd.Flags().GetString("lang")

	var sourceTextIterator iter.Seq[string]

	if len(asInlineText) > 0 {
		sourceTextIterator = data.NewStringIterator(asInlineText)
	} else {
		file, err := os.Open(asFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		sourceTextIterator = data.NewFileReaderIterator(file)
	}

	switch lang {
	case "eng":
		return stemming.NewEnglishTokenIterator(sourceTextIterator)
	case "":
		return stemming.NewEnglishTokenIterator(sourceTextIterator)
	}

	panic(fmt.Sprintf("Unsupported language: %s", lang))
}
