package indexing

type node struct {
	documentId uint64
	position   int
	nextNode   *node
}

type InvertedIndex struct {
	invertedLists map[string]*node
}

func MakeNewEmptyInvertedIndex() InvertedIndex {
	return InvertedIndex{
		invertedLists: make(map[string]*node),
	}
}

func (ridx *InvertedIndex) generateCursorsMap(query []string) map[string]*node {
	cursorsMap := make(map[string]*node)
	for _, term := range query {
		if listRoot, ok := ridx.invertedLists[term]; ok {
			cursorsMap[term] = listRoot
		} else {
			cursorsMap[term] = nil
		}
	}
	return cursorsMap
}

func (ridx *InvertedIndex) Search(query []string, weightsMap map[string]uint64) SearchResults {
	search_results := SearchResults{
		scoreByDocumentId:           make(map[uint64]uint64),
		sortedBestMatchingDocuments: make([]uint64, 0),
	}
	cursorsMap := ridx.generateCursorsMap(query)
	for term, cursor := range cursorsMap {
		for cursor != nil {
			weight := weightsMap[term]
			search_results.incrementScore(cursor.documentId, weight)
			cursor = cursor.nextNode
		}
	}
	return search_results
}

func (ridx *InvertedIndex) generateInsertionMap(documentTokens []string, documentId uint64) map[string]*node {
	cursorsMap := ridx.generateCursorsMap(documentTokens)
	insertionMap := make(map[string]*node)
	for term, cursor := range cursorsMap {
		if cursor != nil && cursor.documentId > documentId {
			cursor = nil
		}
		for cursor != nil && cursor.nextNode != nil && cursor.nextNode.documentId < documentId {
			cursor = cursor.nextNode
		}
		insertionMap[term] = cursor
	}
	return insertionMap
}

func (ridx *InvertedIndex) Store(documentTokens []string, documentId uint64) {
	insertionMap := ridx.generateInsertionMap(documentTokens, documentId)
	for termPositionInQuery, term := range documentTokens {
		ptrToNewNode := &node{documentId: documentId, position: termPositionInQuery, nextNode: nil}
		insertionCursor := insertionMap[term]
		if insertionCursor == nil {
			ptrToNewNode.nextNode = ridx.invertedLists[term]
			ridx.invertedLists[term] = ptrToNewNode
		} else {
			ptrToNewNode.nextNode = insertionCursor.nextNode
			insertionCursor.nextNode = ptrToNewNode
		}
		insertionMap[term] = ptrToNewNode
	}
}
