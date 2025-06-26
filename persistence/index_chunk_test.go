package persistence

import (
	"quinto/core"
	"quinto/data"
	"testing"
)

func TestNewSortedArrayOfTermTrackers(t *testing.T) {
	sortedArray := newSortedArrayOfTermTrackers()
	termTrackers := []core.TermTracker{
		{DocId: 17, Position: 3},
		{DocId: 17, Position: 4},
		{DocId: 27, Position: 1},
		{DocId: 37, Position: 1},
	}
	for _, tracker := range termTrackers {
		if !sortedArray.Insert(tracker) {
			t.Errorf("Failed to insert term tracker: %v", tracker)
		}
	}
	if sortedArray.Size() != len(termTrackers) {
		t.Errorf("Expected size %d, got %d", len(termTrackers), sortedArray.Size())
	}
}

func TestIndexChunkWriteAndRead(t *testing.T) {
	termTrackers := []core.TermTracker{
		{DocId: 17, Position: 3},
		{DocId: 17, Position: 4},
		{DocId: 27, Position: 1},
		{DocId: 37, Position: 1},
	}

	var handler diskHandler = newMockDiskHandler()
	writerChunk := newIndexChunk("hello", "testChunk", handler)
	writerChunk.insertIterable(data.NewSliceIterator(termTrackers))
	writerChunk.writeBack()

	readerChunk := newIndexChunk("hello", "testChunk", handler)

	if readerChunk.chunkKey != "testChunk" {
		t.Errorf("Expected chunkKey to be 'testChunk', got '%s'", readerChunk.chunkKey)
	}

	if readerChunk.nextChunkKey != "" {
		t.Errorf("Expected nextChunkKey to be empty, got '%s'", readerChunk.nextChunkKey)
	}

	if readerChunk.termTrackers.Size() != len(termTrackers) {
		t.Errorf("Expected termTrackers size to be %d, got %d", len(termTrackers), readerChunk.termTrackers.Size())
	}
}
