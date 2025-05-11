package search

import (
	"quinto/misc"
)

type Match struct {
	success       bool
	DocumentId    uint64
	StartPosition int
	EndPosition   int
}

type Query interface {
	Run() Match
	Advance()
	Close()
	Init(misc.ReverseIndex)

	lowestDocumentId() uint64
}
