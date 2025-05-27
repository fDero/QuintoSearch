/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

Storing and loading inverted lists on/from disk is a crucial part of the persistence
layer. This file contains the functions that are used to store and load inverted lists
from disk. The inverted lists are stored in a custom binary format that uses v-byte
encoding to compress the document IDs and positions. This format is optimized for
space efficiency and fast read/write operations.
==================================================================================*/

package persistence

import (
	"io"
	"iter"
	"quinto/core"
)

func encodeTermTrackersToDisk(fileWriter io.Writer, invertedListIterator iter.Seq[core.TermTracker]) error {
	lastDocumentId := core.DocumentId(0)
	lastPosition := core.TermPosition(0)
	for tracker := range invertedListIterator {

		documentIdDelta := tracker.DocId - lastDocumentId
		positionDelta := tracker.Position - lastPosition

		positionToEncode := positionDelta
		if documentIdDelta == 0 {
			positionToEncode = tracker.Position
		}

		encodedDocumentId := vbyteEncodeUInt64(documentIdDelta)
		encodedPosition := vbyteEncodeUInt64(uint64(positionToEncode))

		if _, err := fileWriter.Write(encodedDocumentId); err != nil {
			return err
		}
		if _, err := fileWriter.Write(encodedPosition); err != nil {
			return err
		}

		lastDocumentId = tracker.DocId
		lastPosition = tracker.Position
	}
	return nil
}

func processTermTrackersFromDisk(fileReader io.ByteReader, yield func(core.TermTracker) bool) error {

	documentId := core.DocumentId(0)
	position := core.TermPosition(0)

	for {
		encodedDocumentIdDelta, idErr := loadVbyteEncodedUInt64(fileReader)
		if idErr == io.EOF {
			return nil
		}
		if idErr != nil {
			return idErr
		}

		encodedPositionMaybeDeltaMaybeAbsolute, posErr := loadVbyteEncodedUInt64(fileReader)
		if posErr != nil {
			return posErr
		}

		documentIdDelta := core.DocumentId(vbyteDecodeUInt64(encodedDocumentIdDelta))
		positionMaybeDeltaMaybeAbsolute := vbyteDecodeUInt64(encodedPositionMaybeDeltaMaybeAbsolute)

		documentId += documentIdDelta
		if documentIdDelta == 0 {
			position += core.TermPosition(positionMaybeDeltaMaybeAbsolute)
		} else {
			position = core.TermPosition(positionMaybeDeltaMaybeAbsolute)
		}

		keepGoing := yield(core.TermTracker{
			DocId:    documentId,
			Position: position,
		})

		if !keepGoing {
			return nil
		}
	}
}

func iterateTermTrackersFromDisk(fileReader io.ByteReader) iter.Seq[core.TermTracker] {
	return func(yield func(core.TermTracker) bool) {
		processTermTrackersFromDisk(fileReader, yield)
	}
}
