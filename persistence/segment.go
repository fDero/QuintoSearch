package persistence

import (
	"iter"
	"quinto/misc"
)

type segment struct {
	head *segmentNode
	tail *segmentNode

	size              uint64
	lowestDocumentId  uint64
	highestDocumentId uint64
}

type segmentNode struct {
	tracker misc.TermTracker
	next    *segmentNode
}

func newSegment() *segment {
	return &segment{
		head:              nil,
		tail:              nil,
		size:              0,
		lowestDocumentId:  0,
		highestDocumentId: 0,
	}
}

func (invlst *segment) iterator() iter.Seq[misc.TermTracker] {
	return func(yield func(misc.TermTracker) bool) {
		for current := invlst.head; current != nil; current = current.next {
			if !yield(current.tracker) {
				break
			}
		}
	}
}

func (invlst *segment) getInsertionPoint(tracker misc.TermTracker) **segmentNode {
	var insertionPoint **segmentNode = nil
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

func (invlst *segment) add(tracker misc.TermTracker) {
	var insertionPoint **segmentNode = invlst.getInsertionPoint(tracker)
	var newNextNode *segmentNode = nil
	if *insertionPoint != nil {
		newNextNode = (*insertionPoint).next
	}
	*insertionPoint = &segmentNode{
		tracker: tracker,
		next:    newNextNode,
	}
	if invlst.tail == nil {
		invlst.tail = *insertionPoint
	}
	invlst.size++
	invlst.highestDocumentId = max(invlst.highestDocumentId, tracker.DocumentId)
	invlst.lowestDocumentId = min(invlst.lowestDocumentId, tracker.DocumentId)
}
