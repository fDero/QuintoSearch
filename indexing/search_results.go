package indexing

import (
	"sort"
)

type SearchResults struct {
	scoreByDocumentId           map[uint64]uint64
	sortedBestMatchingDocuments []uint64
}

func (sr *SearchResults) incrementScore(documentId uint64, additiveIncrement uint64) {
	newScore := additiveIncrement
	if oldScore, ok := sr.scoreByDocumentId[documentId]; ok {
		newScore += oldScore
	}
	sr.scoreByDocumentId[documentId] = newScore
	if len(sr.sortedBestMatchingDocuments) > 0 {
		sr.sortedBestMatchingDocuments = []uint64{}
	}
}

func (sr *SearchResults) updateBestMatches() {
	if len(sr.sortedBestMatchingDocuments) == 0 {
		for documentId := range sr.scoreByDocumentId {
			sr.sortedBestMatchingDocuments = append(sr.sortedBestMatchingDocuments, documentId)
		}
		sort.Slice(sr.sortedBestMatchingDocuments, func(lx, rx int) bool {
			lxid := sr.sortedBestMatchingDocuments[lx]
			rxid := sr.sortedBestMatchingDocuments[rx]
			lxscore := sr.scoreByDocumentId[lxid]
			rxscore := sr.scoreByDocumentId[rxid]
			return lxscore > rxscore
		})
	}
}

func (sr *SearchResults) GetBestMatches(pageSize int, pageIndex int) []uint64 {
	sr.updateBestMatches()
	totLen := len(sr.sortedBestMatchingDocuments)
	var selectedResults []uint64
	firstIndex := pageSize * pageIndex
	lastIndex := pageSize * (pageIndex + 1)
	for i := firstIndex; i < totLen && i < lastIndex; i++ {
		docId := sr.sortedBestMatchingDocuments[i]
		selectedResults = append(selectedResults, docId)
	}
	return selectedResults
}

func (sr *SearchResults) GetSizeInPages(pageSize int) int {
	sr.updateBestMatches()
	totLen := len(sr.sortedBestMatchingDocuments)
	roundedDown := totLen / pageSize
	if totLen%pageSize != 0 {
		return roundedDown + 1
	} else {
		return roundedDown
	}
}
