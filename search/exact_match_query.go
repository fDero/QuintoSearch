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
	"quinto/misc"
)

type ExactQuery struct {
	term    string
	peek    func() (misc.TermTracker, bool)
	advance func()
	close   func()
}

func NewExactQueryFromSlice(terms []misc.TermTracker) ExactQuery {
	index := 0
	return ExactQuery{
		peek: func() (misc.TermTracker, bool) {
			if index >= len(terms) {
				return misc.TermTracker{}, false
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

func NewExactQuery(iterator iter.Seq[misc.TermTracker]) ExactQuery {
	next, stop := iter.Pull(iterator)
	value, exists := next()
	return ExactQuery{
		peek: func() (misc.TermTracker, bool) {
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

func (q *ExactQuery) Init(index misc.ReverseIndex) {
	tmp := NewExactQuery(index.IterateOverTerms(q.term))
	q.peek = tmp.peek
	q.advance = tmp.advance
	q.close = tmp.close
}

func (q *ExactQuery) Run() Match {
	value, exists := q.peek()
	if !exists {
		return Match{success: false}
	}
	return Match{
		success:       true,
		DocumentId:    value.DocumentId,
		StartPosition: value.Position,
		EndPosition:   value.Position,
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

func (q *ExactQuery) coordinates() (uint64, int) {
	value, exists := q.peek()
	if !exists {
		return 0, -1
	}
	return value.DocumentId, value.Position
}
