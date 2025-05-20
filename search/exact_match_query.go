/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

An "ExactQuery" query is a type of search query that looks for an exact match of a
term in a document. It is intended to be used as the leaf node of a query-tree for
every structred query. "ExactQuery" implements the Query interface, which defines the
"Run", "Advance", and "Close" methods. Please refer to the documentation of the
"Query" interface for more details about its methods and their intended usage.
==================================================================================*/

package search

import (
	"iter"
	"quinto/core"
	"quinto/data"
)

type ExactQuery struct {
	term    string
	peek    func() (core.TermTracker, bool)
	advance func()
	close   func()
}

func NewExactQueryFromSlice(terms []core.TermTracker) ExactQuery {
	index := 0
	return ExactQuery{
		peek: func() (core.TermTracker, bool) {
			if index >= len(terms) {
				return core.TermTracker{}, false
			}
			return terms[index], true
		},
		advance: func() {
			index++
		},
		close: func() {
			index = 0
		},
	}
}

func NewExactQuery(iterator iter.Seq[core.TermTracker]) ExactQuery {
	next, stop := iter.Pull(iterator)
	value, exists := next()
	return ExactQuery{
		peek: func() (core.TermTracker, bool) {
			return value, exists
		},
		advance: func() {
			value, exists = next()
		},
		close: func() {
			stop()
		},
	}
}

func (q *ExactQuery) Init(index core.ReverseIndex) {
	tmp := NewExactQuery(index.IterateOverTerms(q.term))
	q.peek = tmp.peek
	q.advance = tmp.advance
	q.close = tmp.close
}

func (q *ExactQuery) Run() core.Match {
	value, exists := q.peek()
	if !exists {
		return core.Match{Success: false}
	}

	involvedTokens := data.NewSet[core.Token]()
	involvedTokens.InsertOne(core.Token{
		StemmedText: q.term,
		Position:    value.Position,
	})

	return core.Match{
		Success:        true,
		DocId:          value.DocId,
		StartPosition:  value.Position,
		EndPosition:    value.Position,
		InvolvedTokens: involvedTokens,
	}
}

func (q *ExactQuery) Close() {
	if q.close != nil {
		q.close()
		q.peek = nil
		q.advance = nil
		q.close = nil
	}
}

func (q *ExactQuery) Advance() {
	q.advance()
}

func (q *ExactQuery) Ended() bool {
	_, exists := q.peek()
	return !exists
}

func (q *ExactQuery) Coordinates() (core.DocumentId, core.TermPosition) {
	value, exists := q.peek()
	if !exists {
		return 0, 0
	}
	return value.DocId, value.Position
}
