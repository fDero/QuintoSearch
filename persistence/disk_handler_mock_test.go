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

	inputs := []core.TermTracker{
		{DocId: 17, Position: 1},
		{DocId: 17, Position: 2},
		{DocId: 27, Position: 1},
		{DocId: 37, Position: 1},
	}

	encodeTermTrackersToDisk(writer, data.NewSliceIterator(inputs))
	finalize()

	fileReader, exists := handler.getReader("testChunk")
	if !exists {
		t.Fatalf("Expected fileReader to exist after writing")
	}

	if fileReader == nil {
		t.Error("Expected fileReader to be non-nil after writing")
	}

	outputs := data.CollectAsSlice(iterateTermTrackersFromDisk(fileReader))
	if len(outputs) != len(inputs) {
		t.Errorf("Expected 4 term trackers, got %d", len(outputs))
	}

	for in, out := range data.ZipSlices(inputs, outputs) {
		if in.DocId != out.DocId || in.Position != out.Position {
			t.Errorf("Expected term tracker be %v, got %v", in, out)
		}
	}
}

func TestDiskHandlerForWritingAndReadingStuff(t *testing.T) {
	var handler diskHandler = newMockDiskHandler()
	handler.getWriter("testChunk")
	writer, finalize, _ := handler.getWriter("testChunk")

	termTrackers := []core.TermTracker{
		{DocId: 17, Position: 3},
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

func TestDiskHandlerForWritingAndReadingStuffWithBigNumbers(t *testing.T) {
	var handler diskHandler = newMockDiskHandler()
	handler.getWriter("testChunk")
	writer, finalize, _ := handler.getWriter("testChunk")

	bigNumber := 4800034432111100120
	bigDocId := core.DocumentId(bigNumber)
	bigPosition := core.TermPosition(bigNumber)

	termTrackers := []core.TermTracker{
		{DocId: 24, Position: 4},
		{DocId: 30, Position: 9},
		{DocId: 39, Position: 1},
		{DocId: 39, Position: bigPosition},
		{DocId: bigDocId, Position: 1},
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
