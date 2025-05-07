package misc

import (
	"iter"
)

type TermTracker struct {
	DocumentId uint64
	Position   int
}

type TermIterator = iter.Seq[TermTracker]
