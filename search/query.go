package search

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

	lowestDocumentId() uint64
}
