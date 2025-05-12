/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

A "ReverseIndex" is an interface that describes every entity that can be used to
store documents and initialize queries to retrieve them. It is a key component of the
search engine, as it allows for efficient storage and retrieval of documents based
on their content. The "ReverseIndex" is responsable for managing the mapping between
a single term and the documents that contain it, togheter with the positions of the
term in the document. Is essential for the well functioning of the search engine
that terms are iterated in ascending order of document-id and position within a given
document. The "IterateOverTerms" must herby work in this way in every implementation.
==================================================================================*/

package misc

import (
	"iter"
	"sync/atomic"
)

type Token struct {
	StemmedText  string
	OriginalText string
	Position     int
}

type TermTracker struct {
	DocumentId uint64
	Position   int
}

type ReverseIndex interface {
	IterateOverTerms(term string) iter.Seq[TermTracker]
	StoreNewDocument(toks iter.Seq[Token]) (uint64, error)
}

type NaiveReverseIndex struct {
	terms     map[string][]TermTracker
	IdCounter atomic.Uint64
}

func NewNaiveReverseIndex() *NaiveReverseIndex {
	return &NaiveReverseIndex{
		terms:     make(map[string][]TermTracker),
		IdCounter: atomic.Uint64{},
	}
}

func (q *NaiveReverseIndex) IterateOverTerms(term string) iter.Seq[TermTracker] {
	return func(yield func(TermTracker) bool) {
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

func (q *NaiveReverseIndex) StoreNewDocument(toks iter.Seq[Token]) (uint64, error) {
	id := q.IdCounter.Add(1)
	for tok := range toks {
		newTracker := TermTracker{DocumentId: id, Position: tok.Position}
		q.terms[tok.StemmedText] = append(q.terms[tok.StemmedText], newTracker)
	}
	return id, nil
}
