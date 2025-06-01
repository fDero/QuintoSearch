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
	input := []core.TermTracker{
		{DocId: 1, Position: 2},
	}

	buffer := new(bytes.Buffer)
	inputIterator := data.NewSliceIterator(input)
	encodeTermTrackersToDisk(buffer, inputIterator)
	output := data.CollectAsSlice(iterateTermTrackersFromDisk(buffer))

	if len(output) != 1 {
		t.Errorf("Expected one element in the output, got %d elements", len(output))
	}

	badDocId := output[0].DocId != input[0].DocId
	badPosition := output[0].Position != input[0].Position
	if badDocId || badPosition {
		t.Errorf("Expected %v, got %v", input[0], output[0])
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

	for in, out := range data.ZipSlices(input, output) {
		if in.DocId != out.DocId || in.Position != out.Position {
			t.Errorf("Mismatch: expected %v, got %v", in, out)
		}
	}
}

func TestWritingAndReadingMultipleElementsExtended(t *testing.T) {

	bigNumber := uint64(12345678901234567890)
	bigDocId := core.DocumentId(bigNumber)
	bigPosition := core.TermPosition(bigNumber)

	if enc := vbyteEncodeUInt64(bigNumber); len(enc) <= 1 {
		t.Fatalf("we want `bigNumber` to have a large encoding in variable length: %v", bigNumber)
	}

	input := []core.TermTracker{
		{DocId: 1, Position: 2},
		{DocId: 1, Position: 3},
		{DocId: 2, Position: 4},
		{DocId: 2, Position: 5},
		{DocId: 2, Position: bigPosition},
		{DocId: 3, Position: 7},
		{DocId: 3, Position: 8},
		{DocId: bigDocId, Position: 1},
		{DocId: bigDocId, Position: 2},
	}

	buffer := new(bytes.Buffer)
	inputIterator := data.NewSliceIterator(input)
	encodeTermTrackersToDisk(buffer, inputIterator)
	output := data.CollectAsSlice(iterateTermTrackersFromDisk(buffer))

	if len(output) != len(input) {
		t.Errorf("Expected %d elements in the output, got %d elements", len(input), len(output))
	}
}

func TestEncodeDecodeString(t *testing.T) {
	samples := []string{
		"", "hello", "quinto", "120\n", "aaaa \n \raaa \r",
		"a very long string that exceeds the usual length",
	}
	for _, sample := range samples {
		buffer := new(bytes.Buffer)
		if err := encodeStringToDisk(buffer, sample); err != nil {
			t.Fatalf("Failed to encode string: %v", err)
		}
		reader := bytes.NewReader(buffer.Bytes())
		result, err := decodeStringFromDisk(reader)
		if err != nil {
			t.Fatalf("Failed to decode string: %v", err)
		}
		if result != sample {
			t.Errorf("Expected '%s', got '%s'", sample, result)
		}
	}
}
