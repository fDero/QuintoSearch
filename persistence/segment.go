/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

A segment is a sub-unit of an inverted list for a given term. It is implemented as a
linked list of term trackers, which are used to store the document ID and position of
a term in a document. Inverted lists are stored in segments so that they can be locked
and updated independently of other segments. This allows for concurrent reads and writes
to the inverted index without the need for a global lock.
==================================================================================*/

package persistence

import (
	"iter"
	"math"
	"quinto/misc"
	"unsafe"
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
		lowestDocumentId:  math.MaxUint64,
		highestDocumentId: 0,
	}
}

func (seg *segment) estimateSize() int64 {
	estimatedSize := uintptr(0)
	if seg.size != 0 {
		estimatedSize += unsafe.Sizeof(*seg.head)
		estimatedSize *= uintptr(seg.size)
	}
	estimatedSize += unsafe.Sizeof(*seg)
	return int64(estimatedSize)
}

func (seg *segment) iterator() iter.Seq[misc.TermTracker] {
	return func(yield func(misc.TermTracker) bool) {
		for current := seg.head; current != nil; current = current.next {
			if !yield(current.tracker) {
				break
			}
		}
	}
}

func (seg *segment) getInsertionPoint(tracker misc.TermTracker) **segmentNode {
	var insertionPoint **segmentNode = nil
	for insertionPoint = &seg.head; *insertionPoint != nil; insertionPoint = &((*insertionPoint).next) {
		overshootByDocumentId := (*insertionPoint).tracker.DocumentId > tracker.DocumentId
		onExactDocumentId := (*insertionPoint).tracker.DocumentId == tracker.DocumentId
		overshootByPosition := onExactDocumentId && (*insertionPoint).tracker.Position > tracker.Position
		if overshootByDocumentId || overshootByPosition {
			break
		}
	}
	return insertionPoint
}

func (seg *segment) add(tracker misc.TermTracker) {
	var insertionPoint **segmentNode = seg.getInsertionPoint(tracker)
	var newNextNode *segmentNode = nil
	if *insertionPoint != nil {
		newNextNode = *insertionPoint
	}
	*insertionPoint = &segmentNode{
		tracker: tracker,
		next:    newNextNode,
	}
	if newNextNode == nil {
		seg.tail = *insertionPoint
	}
	seg.size++
	seg.highestDocumentId = max(seg.highestDocumentId, tracker.DocumentId)
	seg.lowestDocumentId = min(seg.lowestDocumentId, tracker.DocumentId)
}
