package persistence

import (
	"bytes"
	"quinto/core"
	"quinto/data"
	"testing"
)

func TestWritingAndReadingEmpty(t *testing.T) {
	buffer := new(bytes.Buffer)
	inputIterator := data.NewSliceIterator([]core.TermTracker{})
	encodeTermTrackersToDisk(buffer, inputIterator)
	output := data.CollectAsSlice(iterateTermTrackersFromDisk(buffer))

	if len(output) != 0 {
		t.Errorf("Expected empty output, got %d elements", len(output))
	}
}

func TestWritingAndReadingOneElement(t *testing.T) {
	buffer := new(bytes.Buffer)
	inputIterator := data.NewSliceIterator([]core.TermTracker{
		{DocId: 1, Position: 2},
	})
	encodeTermTrackersToDisk(buffer, inputIterator)
	output := data.CollectAsSlice(iterateTermTrackersFromDisk(buffer))

	if len(output) != 1 {
		t.Errorf("Expected one element in the output, got %d elements", len(output))
	}
}

func TestWritingAndReadingMultipleElements(t *testing.T) {
	input := []core.TermTracker{
		{DocId: 1, Position: 2},
		{DocId: 1, Position: 3},
		{DocId: 2, Position: 4},
		{DocId: 2, Position: 5},
		{DocId: 2, Position: 6},
		{DocId: 3, Position: 7},
		{DocId: 3, Position: 8},
		{DocId: 4, Position: 9},
	}

	buffer := new(bytes.Buffer)
	inputIterator := data.NewSliceIterator(input)
	encodeTermTrackersToDisk(buffer, inputIterator)
	output := data.CollectAsSlice(iterateTermTrackersFromDisk(buffer))

	if len(output) != len(input) {
		t.Errorf("Expected %d elements in the output, got %d elements", len(input), len(output))
	}
}
