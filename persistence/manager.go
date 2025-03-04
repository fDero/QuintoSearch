package persistence

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"quinto/misc"
	"sync"

	"github.com/dgraph-io/ristretto/v2"
)

type PersistenceManager struct {
	segments    *ristretto.Cache[string, *synchronizedSegment]
	mutex       sync.Mutex
	dbDirectory string
}

func NewPersistenceManager(dbDirectory string) *PersistenceManager {
	const oneGigabyteOfStorageCapacity = 1 << 30
	const maximumNumberOfCacheBuckets = 1e7
	cache, err := ristretto.NewCache(&ristretto.Config[string, *synchronizedSegment]{
		NumCounters: maximumNumberOfCacheBuckets,
		MaxCost:     oneGigabyteOfStorageCapacity,
		BufferItems: 64,
	})
	if err != nil {
		panic("Failed to create cache while constructing 'PersistenceManager'")
	}
	return &PersistenceManager{
		segments:    cache,
		mutex:       sync.Mutex{},
		dbDirectory: dbDirectory,
	}
}

func (pm *PersistenceManager) SearchOnDisk(term string, blockCounter int) (*synchronizedSegment, bool) {
	currentKey := fmt.Sprint(term, "_", blockCounter)
	files, err := os.ReadDir(pm.dbDirectory)
	if err != nil {
		panic("Failed to access db-directory")
	}
	for _, file := range files {
		if file.Name() == currentKey {
			completePath := fmt.Sprint(pm.dbDirectory, "/", currentKey)
			file, _ := os.Open(completePath)
			reader := bufio.NewReader(file)
			extracted, _ := LoadFromDisk(reader)
			syncseg := newSynchronizedSegment(extracted)
			pm.segments.Set(currentKey, syncseg, syncseg.estimateSize())
			return syncseg, true
		}
	}
	return nil, false
}

func (pm *PersistenceManager) GetInvertedListIterator(term string) iter.Seq[misc.TermTracker] {
	return func(yield func(misc.TermTracker) bool) {
		for counter := 0; true; counter++ {
			currentKey := fmt.Sprint(term, "_", counter)
			syncseg, found := pm.segments.Get(currentKey)
			if !found {
				syncseg, found = pm.SearchOnDisk(term, counter)
			}
			if !found {
				break
			}
			for term := range syncseg.iterator() {
				yield(term)
			}
		}
	}
}

func (pm *PersistenceManager) getCacheKey(term string, documentId uint64) string {
	for counter := 0; true; counter++ {
		currentKey := fmt.Sprint(term, "_", counter)
		segment, found := pm.segments.Get(currentKey)
		if !found || segment.underlyng.size < 1 {
			return currentKey
		}
		if segment.underlyng.tail.tracker.DocumentId <= documentId {
			return currentKey
		}
	}
	panic("illegal state during execution of 'PersistenceManager.getCacheKey'")
}

func (pm *PersistenceManager) StoreNewDocument(documentId uint64, documentInputTokenStream iter.Seq[misc.Token]) {
	segmentsForGivenTermInCurrentDocument := make(map[string]string)
	for token := range documentInputTokenStream {
		key, found := segmentsForGivenTermInCurrentDocument[token.StemmedText]
		if !found {
			key = pm.getCacheKey(token.StemmedText, documentId)
			segmentsForGivenTermInCurrentDocument[token.StemmedText] = key
		}
		pm.mutex.Lock()
		syncseg, found := pm.segments.Get(key)
		if !found {
			syncseg = newSynchronizedSegment(nil)
			pm.segments.Set(key, syncseg, syncseg.estimateSize())
			pm.segments.Wait()
		}
		pm.mutex.Unlock()
		syncseg.add(misc.TermTracker{
			DocumentId: documentId,
			Position:   token.Position,
		})
		pm.segments.Set(key, syncseg, syncseg.estimateSize())
	}
}
