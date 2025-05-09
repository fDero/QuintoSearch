package search

import (
	"iter"
	"quinto/misc"
)

type Exact struct {
	peek    func() (misc.TermTracker, bool)
	advance func()
	close   func()
}

func NewExactQueryFromSlice(terms []misc.TermTracker) Exact {
	index := 0
	return Exact{
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

func NewExactQuery(iterator iter.Seq[misc.TermTracker]) Exact {
	next, stop := iter.Pull(iterator)
	value, exists := next()
	return Exact{
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

func (q *Exact) Run() Match {
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

func (q *Exact) Close() {
	if q.close != nil {
		q.close()
		q.peek = nil
		q.advance = nil
		q.close = nil
	}
}

func (q *Exact) Advance() {
	q.advance()
}

func (q *Exact) lowestDocumentId() uint64 {
	value, exists := q.peek()
	if !exists {
		return 0
	}
	return value.DocumentId
}
