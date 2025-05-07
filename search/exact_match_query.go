package search

import (
	"quinto/misc"
)

type Exact struct {
	peek    func() (misc.TermTracker, bool)
	advance func()
	close   func()
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
