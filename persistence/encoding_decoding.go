/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

Storing and loading data on/from disk is a crucial part of the persistence layer.
This file contains the functions that are used to store and load data from disk. The
API is designed around the `persistence.diskHandler` interface, which provides a
layer of abstraction over the disk operations. This turns out to be useful for
testing purposes, as it allows us to mock the disk operations and test the
persistence layer without actually writing to disk.
==================================================================================*/

package persistence

import (
	"io"
	"iter"
	"quinto/core"
)

func encodeStringToDisk(fileWriter io.Writer, text string) error {
	encodedLen := vbyteEncodeUInt64(uint64(len(text)))
	if _, err := fileWriter.Write(encodedLen); err != nil {
		return err
	}
	if _, err := fileWriter.Write([]byte(text)); err != nil {
		return err
	}
	return nil
}

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

func decodeStringFromDisk(fileReader io.ByteReader) (string, error) {
	encodedLen, err := loadVbyteEncodedUInt64(fileReader)
	decodedLen := vbyteDecodeUInt64(encodedLen)
	if err != nil || decodedLen == 0 {
		return "", err
	}
	bytes := make([]byte, decodedLen)
	for i := range decodedLen {
		bytes[i], err = fileReader.ReadByte()
		if err != nil {
			return "", err
		}
	}
	return string(bytes), nil
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
			position = core.TermPosition(positionMaybeDeltaMaybeAbsolute)
		} else {
			position += core.TermPosition(positionMaybeDeltaMaybeAbsolute)
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
