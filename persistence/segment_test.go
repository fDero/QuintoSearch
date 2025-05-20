package persistence

import (
	"quinto/core"
	"testing"
)

func TestAddToEmptySegment(t *testing.T) {
	segment := newSegment()
	tracker := core.TermTracker{DocId: 1, Position: 2}
	segment.add(tracker)
	if segment.head == nil || segment.tail == nil {
		t.Errorf("Expected segment to have head and tail, but they are nil")
	}
	if segment.head.tracker != tracker {
		t.Errorf("Expected head tracker to be %v, but got %v", tracker, segment.head.tracker)
	}
	if segment.tail.tracker != tracker {
		t.Errorf("Expected tail tracker to be %v, but got %v", tracker, segment.tail.tracker)
	}
	if segment.size != 1 {
		t.Errorf("Expected segment size to be 1, but got %d", segment.size)
	}
	if segment.lowestDocumentId != 1 {
		t.Errorf("Expected lowestDocumentId to be 1, but got %d", segment.lowestDocumentId)
	}
	if segment.highestDocumentId != 1 {
		t.Errorf("Expected highestDocumentId to be 1, but got %d", segment.highestDocumentId)
	}
	if segment.head != segment.tail {
		t.Errorf("Expected head and tail to be the same node, but they are different")
	}
}

func TestAddToNonEmptySegment(t *testing.T) {
	segment := newSegment()
	tracker := core.TermTracker{DocId: 1, Position: 2}
	tracker2 := core.TermTracker{DocId: 50, Position: 2}
	segment.add(tracker)
	segment.add(tracker2)
	if segment.size != 2 {
		t.Errorf("Expected segment size to be 2, but got %d", segment.size)
	}
	if segment.lowestDocumentId != 1 {
		t.Errorf("Expected lowestDocumentId to be 1, but got %d", segment.lowestDocumentId)
	}
	if segment.highestDocumentId != 50 {
		t.Errorf("Expected highestDocumentId to be 50, but got %d", segment.highestDocumentId)
	}
	if segment.head.tracker != tracker {
		t.Errorf("Expected head tracker to be %v, but got %v", tracker, segment.head.tracker)
	}
	if segment.tail.tracker != tracker2 {
		t.Errorf("Expected tail tracker to be %v, but got %v", tracker2, segment.tail.tracker)
	}
	if segment.tail.next != nil {
		t.Errorf("Expected tail.next to be nil, but got %v", segment.tail.next)
	}
	if segment.tail != segment.head.next {
		t.Errorf("Expected tail to be head.next, but got %v", segment.tail)
	}
}

func TestAddToSegmentWithSameDocumentId(t *testing.T) {
	segment := newSegment()
	tracker := core.TermTracker{DocId: 1, Position: 2}
	tracker2 := core.TermTracker{DocId: 50, Position: 2}
	segment.add(tracker)
	segment.add(tracker2)
	segment.add(tracker)
	segment.add(tracker)
	if segment.size != 4 {
		t.Errorf("Expected segment size to be 4, but got %d", segment.size)
	}
	if segment.lowestDocumentId != 1 {
		t.Errorf("Expected lowestDocumentId to be 1, but got %d", segment.lowestDocumentId)
	}
	if segment.highestDocumentId != 50 {
		t.Errorf("Expected highestDocumentId to be 50, but got %d", segment.highestDocumentId)
	}
}

func TestAddInTheMiddleOfSegment(t *testing.T) {
	segment := newSegment()
	tracker := core.TermTracker{DocId: 1, Position: 2}
	tracker3 := core.TermTracker{DocId: 50, Position: 2}
	tracker2 := core.TermTracker{DocId: 25, Position: 2}
	segment.add(tracker)
	segment.add(tracker3)
	segment.add(tracker2)
	if segment.size != 3 {
		t.Errorf("Expected segment size to be 2, but got %d", segment.size)
	}
	if segment.lowestDocumentId != 1 {
		t.Errorf("Expected lowestDocumentId to be 1, but got %d", segment.lowestDocumentId)
	}
	if segment.highestDocumentId != 50 {
		t.Errorf("Expected highestDocumentId to be 50, but got %d", segment.highestDocumentId)
	}
	if segment.head.tracker != tracker {
		t.Errorf("Expected head tracker to be %v, but got %v", tracker, segment.head.tracker)
	}
	if segment.head.next.tracker != tracker2 {
		t.Errorf("Expected head.next tracker to be %v, but got %v", tracker2, segment.head.next.tracker)
	}
	if segment.tail.tracker != tracker3 {
		t.Errorf("Expected head.next tracker to be %v, but got %v", tracker3, segment.head.next.tracker)
	}
	for it := range segment.iterator() {
		t.Log(it)
	}
}
