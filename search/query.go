package search

import (
	"quinto/misc"
)

type Match struct {
	Success        bool
	DocumentId     uint64
	StartPosition  int
	EndPosition    int
	InvolvedTokens misc.Set[misc.Token]
}

type Query interface {
	Run() Match
	Advance()
	Ended() bool
	Close()
	Init(misc.ReverseIndex)

	coordinates() (uint64, int)
}
