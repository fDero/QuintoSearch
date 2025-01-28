package indexing

import (
	"iter"
)

func (invlst *InvertedList) Iterator() iter.Seq[TermTracker] {
	return func(yield func(TermTracker) bool) {
		for current := invlst.head; current != nil; current = current.next {
			if !yield(current.tracker) {
				break
			}
		}
	}
}

func (invlst *InvertedList) getInsertionPoint(tracker TermTracker) **invertedListNode {
	var insertionPoint **invertedListNode = nil
	for insertionPoint = &invlst.head; *insertionPoint != nil; insertionPoint = &((*insertionPoint).next) {
		overshootByDocumentId := (*insertionPoint).tracker.DocumentId > tracker.DocumentId
		onExactDocumentId := (*insertionPoint).tracker.DocumentId == tracker.DocumentId
		overshootByPosition := onExactDocumentId && (*insertionPoint).tracker.Position > tracker.Position
		if overshootByDocumentId || overshootByPosition {
			break
		}
	}
	return insertionPoint
}

func (invlst *InvertedList) Add(tracker TermTracker) {
	var insertionPoint **invertedListNode = invlst.getInsertionPoint(tracker)
	var newNextNode *invertedListNode = nil
	if *insertionPoint != nil {
		newNextNode = (*insertionPoint).next
	}
	*insertionPoint = &invertedListNode{
		tracker: tracker,
		next:    newNextNode,
	}
	if invlst.tail == nil {
		invlst.tail = *insertionPoint
	}
	invlst.size++
}
