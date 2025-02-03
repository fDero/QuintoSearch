package persistence

import (
	"bufio"
)

func StoreOnDisk(fileWriter bufio.Writer, invertedList segment) error {
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

func LoadFromDisk(fileReader bufio.Reader) (segment, error) {
	var invertedList segment

	documentId := uint64(0)
	position := 0

	for {
		encodedDocumentIdDelta, idErr := loadVbyteEncodedUInt64(fileReader)
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

		invertedList.add(TermTracker{
			DocumentId: documentId,
			Position:   position,
		})
	}
}
