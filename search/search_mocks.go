package search

import (
	"iter"
	"quinto/misc"
	"sync/atomic"
)

var helloWorldDocument = createDummyDocument([]string{
	"hello",
	"world",
})

var guitarDocument = createDummyDocument([]string{
	"guitar",
	"string",
	"instrument",
	"band",
	"important",
	"music",
	"instrument",
})

var hobbyDocument = createDummyDocument([]string{
	"love",
	"music",
	"chess",
	"science",
})

var toolsDocument = createDummyDocument([]string{
	"screwdriver",
	"hammer",
	"instrument",
	"drill",
	"wrench",
})

func createDummyDocument(tokens []string) []misc.Token {
	out := []misc.Token{}
	for i, token := range tokens {
		out = append(out, misc.Token{
			StemmedText:  token,
			Position:     i,
		})
	}
	return out
}

type NaiveReverseIndex struct {
	terms     map[string][]misc.TermTracker
	IdCounter atomic.Uint64
}

func NewNaiveReverseIndex() *NaiveReverseIndex {
	return &NaiveReverseIndex{
		terms:     make(map[string][]misc.TermTracker),
		IdCounter: atomic.Uint64{},
	}
}

func (q *NaiveReverseIndex) IterateOverTerms(term string) iter.Seq[misc.TermTracker] {
	return func(yield func(misc.TermTracker) bool) {
		termTrackers, exists := q.terms[term]
		if !exists {
			return
		}
		for _, termTracker := range termTrackers {
			if !yield(termTracker) {
				return
			}
		}
	}
}

func (q *NaiveReverseIndex) StoreNewDocument(toks iter.Seq[misc.Token]) (uint64, error) {
	id := q.IdCounter.Add(1)
	for tok := range toks {
		newTracker := misc.TermTracker{DocumentId: id, Position: tok.Position}
		q.terms[tok.StemmedText] = append(q.terms[tok.StemmedText], newTracker)
	}
	return id, nil
}
