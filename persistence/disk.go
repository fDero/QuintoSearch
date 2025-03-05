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
	"bufio"
	"io"
	"quinto/misc"
)

func StoreOnDisk(fileWriter *bufio.Writer, invertedList *segment) error {
	lastDocumentId := uint64(0)
	lastPosition := 0
	for tracker := range invertedList.iterator() {

		documentIdDelta := tracker.DocumentId - lastDocumentId
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

		lastDocumentId = tracker.DocumentId
		lastPosition = tracker.Position
	}
	return nil
}

func LoadFromDisk(fileReader *bufio.Reader) (*segment, error) {
	var invertedList *segment = newSegment()

	documentId := uint64(0)
	position := 0

	for {
		encodedDocumentIdDelta, idErr := loadVbyteEncodedUInt64(fileReader)
		if idErr == io.EOF {
			return invertedList, nil
		}
		if idErr != nil {
			return invertedList, idErr
		}

		encodedPositionMaybeDeltaMaybeAbsolute, posErr := loadVbyteEncodedUInt64(fileReader)
		if posErr != nil {
			return invertedList, posErr
		}

		documentIdDelta := vbyteDecodeUInt64(encodedDocumentIdDelta)
		positionMaybeDeltaMaybeAbsolute := vbyteDecodeUInt64(encodedPositionMaybeDeltaMaybeAbsolute)

		documentId += documentIdDelta
		if documentIdDelta == 0 {
			position += int(positionMaybeDeltaMaybeAbsolute)
		} else {
			position = int(positionMaybeDeltaMaybeAbsolute)
		}

		invertedList.add(misc.TermTracker{
			DocumentId: documentId,
			Position:   position,
		})
	}
}
