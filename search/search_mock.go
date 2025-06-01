package search

import (
	"iter"
	"quinto/core"
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

func createDummyDocument(tokens []string) []core.Token {
	out := []core.Token{}
	for i, token := range tokens {
		out = append(out, core.Token{
			StemmedText: token,
			Position:    core.TermPosition(i),
		})
	}
	return out
}

type NaiveReverseIndex struct {
	terms     map[string][]core.TermTracker
	IdCounter atomic.Uint64
}

func NewNaiveReverseIndex() *NaiveReverseIndex {
	return &NaiveReverseIndex{
		terms:     make(map[string][]core.TermTracker),
		IdCounter: atomic.Uint64{},
	}
}

func (q *NaiveReverseIndex) IterateOverTerms(term string) iter.Seq[core.TermTracker] {
	return func(yield func(core.TermTracker) bool) {
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

func (q *NaiveReverseIndex) StoreNewDocument(toks iter.Seq[core.Token]) (core.DocumentId, error) {
	id := core.DocumentId(q.IdCounter.Add(1))
	for tok := range toks {
		newTracker := core.TermTracker{DocId: id, Position: tok.Position}
		q.terms[tok.StemmedText] = append(q.terms[tok.StemmedText], newTracker)
	}
	return id, nil
}
