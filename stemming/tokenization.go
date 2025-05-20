package stemming

import (
	"iter"
	"quinto/core"
	"quinto/data"
	"strings"
)

func NewTokenIterator(
	sourceTextIterator iter.Seq[string],
	stopWords data.Set[string],
	stemmer func(string) string,
) iter.Seq[core.Token] {

	return func(yield func(core.Token) bool) {

		if sourceTextIterator == nil {
			return
		}

		position := core.TermPosition(0)
		for originalTokenText := range sourceTextIterator {

			if originalTokenText == "" {
				continue
			}

			lowerCasedTokenText := strings.ToLower(originalTokenText)
			if stopWords.Contains(lowerCasedTokenText) {
				continue
			}

			mustContinue := yield(core.Token{
				Position:     position,
				OriginalText: originalTokenText,
				StemmedText:  stemmer(lowerCasedTokenText),
			})

			if !mustContinue {
				break
			}
		}
	}
}

func NewEnglishTokenIterator(sourceTextIterator iter.Seq[string]) iter.Seq[core.Token] {
	return NewTokenIterator(
		sourceTextIterator,
		stopWordsEnglish(),
		stemEnglish,
	)
}
