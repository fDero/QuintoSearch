package search

import (
	"quinto/misc"
)

type Match struct {
	Success        bool
	DocId          misc.DocumentId
	StartPosition  misc.TermPosition
	EndPosition    misc.TermPosition
	InvolvedTokens misc.Set[misc.Token]
}

type Query interface {
	Run() Match
	Advance()
	Ended() bool
	Close()
	Init(misc.ReverseIndex)

	coordinates() (misc.DocumentId, misc.TermPosition)
}
