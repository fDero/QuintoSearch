package persistence

type PersistenceManager struct {
	recapDirectoty   string
	segmentDirectory string
	latestDocumentId uint64
	currentTick      uint64
}

func NewPersistenceManager(dbDirectory string) *PersistenceManager {
	return &PersistenceManager{
		segmentDirectory: dbDirectory + "/segments",
		recapDirectoty:   dbDirectory + "/recap",
		latestDocumentId: 0,
		currentTick:      0,
	}
}

func (pm *PersistenceManager) GenerateDocumentId() uint64 {
	documentId := 1 + pm.latestDocumentId
	pm.latestDocumentId = documentId
	return documentId
}
