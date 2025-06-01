package persistence

import (
	"quinto/core"
	"quinto/data"
	"testing"
)

func TestDiskHandlerForWritingAndReadingTermTrackers(t *testing.T) {
	var handler diskHandler = newMockDiskHandler()
	handler.getWriter("testChunk")
	writer, finalize, _ := handler.getWriter("testChunk")

	termTrackers := []core.TermTracker{
		{DocId: 17, Position: 7},
		{DocId: 17, Position: 4},
		{DocId: 27, Position: 1},
		{DocId: 37, Position: 1},
	}

	encodeTermTrackersToDisk(writer, data.NewSliceIterator(termTrackers))
	finalize()

	fileReader, exists := handler.getReader("testChunk")
	if !exists {
		t.Fatalf("Expected fileReader to exist after writing")
	}

	if fileReader == nil {
		t.Error("Expected fileReader to be non-nil after writing")
	}

	chk := data.CollectAsSlice(iterateTermTrackersFromDisk(fileReader))
	if len(chk) != len(termTrackers) {
		t.Errorf("Expected 4 term trackers, got %d", len(chk))
	}
}

func TestDiskHandlerForWritingAndReadingStuff(t *testing.T) {
	var handler diskHandler = newMockDiskHandler()
	handler.getWriter("testChunk")
	writer, finalize, _ := handler.getWriter("testChunk")

	termTrackers := []core.TermTracker{
		{DocId: 17, Position: 7},
		{DocId: 17, Position: 4},
		{DocId: 27, Position: 1},
		{DocId: 37, Position: 1},
	}

	encodeStringToDisk(writer, "hello")
	encodeStringToDisk(writer, "testChunk")
	encodeStringToDisk(writer, "")
	encodeTermTrackersToDisk(writer, data.NewSliceIterator(termTrackers))
	finalize()

	reader, exists := handler.getReader("testChunk")
	if !exists {
		t.Fatalf("Expected fileReader to exist after writing")
	}

	if reader == nil {
		t.Error("Expected fileReader to be non-nil after writing")
	}

	for range 3 {
		decodeStringFromDisk(reader)
	}

	chk := data.CollectAsSlice(iterateTermTrackersFromDisk(reader))
	if len(chk) != len(termTrackers) {
		t.Errorf("Expected 4 term trackers, got %d", len(chk))
	}
}
