package persistence

import (
	"fmt"
	"iter"
	"quinto/misc"
	"sync"

	"github.com/dgraph-io/ristretto/v2"
)

type PersistenceManager struct {
	segments_cache *ristretto.Cache[string, *synchronizedSegment]
	cache_mutex    sync.Mutex
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
		segments_cache: cache,
	}
}

func (pm *PersistenceManager) GetInvertedListIterator(term string) iter.Seq[misc.TermTracker] {
	return func(yield func(misc.TermTracker) bool) {
		for counter := 0; true; counter++ {
			currentKey := fmt.Sprint(term, "_", counter)
			syncseg, found := pm.segments_cache.Get(currentKey)
			if !found {
				// TODO: actually it should be fetched from disk
				// only if not found, iteration should stop
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
		segment, found := pm.segments_cache.Get(currentKey)
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
		pm.cache_mutex.Lock()
		syncseg, found := pm.segments_cache.Get(key)
		if !found {
			syncseg = newSynchronizedSegment()
			pm.segments_cache.Set(key, syncseg, syncseg.estimateSize())
			pm.segments_cache.Wait()
		}
		pm.cache_mutex.Unlock()
		syncseg.add(misc.TermTracker{
			DocumentId: documentId,
			Position:   token.Position,
		})
		pm.segments_cache.Set(key, syncseg, syncseg.estimateSize())
	}
}
