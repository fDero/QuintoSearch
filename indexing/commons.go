package indexing

type TermTracker struct {
	DocumentId uint64
	Position   int
}

type InvertedList struct {
	head *invertedListNode
	tail *invertedListNode
	size uint64
}

type invertedListNode struct {
	tracker TermTracker
	next    *invertedListNode
}

type SearchResults struct {
	scoreByDocumentId           map[uint64]uint64
	sortedBestMatchingDocuments []uint64
}
