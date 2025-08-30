package persistence

import (
	"iter"
	"quinto/core"
	"quinto/data"
	"sync/atomic"
)

type wrappedIndexChunk struct {
	chunk     *indexChunk
	listEntry data.ConcurrentListEntry[string]
}

type PersistenceConfig struct {
	MaxCachedChunks int64
	MaxChunkSize    int
	IoHandler       diskHandler
}

type PersistenceManager struct {
	cacheSize   atomic.Int64
	config      PersistenceConfig
	chunkPool   data.ConcurrentMap[string, wrappedIndexChunk]
	accessList  data.ConcurrentList[string]
	pendingSync data.ConcurrentQueue[string]
}

func NewPersistenceManager(config PersistenceConfig) *PersistenceManager {
	return &PersistenceManager{
		config:      config,
		chunkPool:   *data.NewConcurrentMap[string, wrappedIndexChunk](),
		accessList:  *data.NewLinkedList[string](),
		pendingSync: *data.NewConcurrentQueue[string](),
	}
}

func (pm *PersistenceManager) evictNotPendingLRU() {
	for {
		for listEntry := range pm.accessList.IterateBackwards() {
			wrappedChunk, exists := pm.chunkPool.Get(listEntry.Value())
			if exists && !wrappedChunk.chunk.pendingWriteBack {
				pm.chunkPool.Delete(listEntry.Value())
				listEntry.Remove()
				pm.cacheSize.Add(-1)
				return
			}
		}
	}
}

func (pm *PersistenceManager) retrieveChunkFromCache(key string) *indexChunk {
	wrappedChunk, exists := pm.chunkPool.Get(key)
	if exists {
		wrappedChunk.listEntry.Remove()
		pm.accessList.InsertFront(key)
		return wrappedChunk.chunk
	}
	return nil
}

func (pm *PersistenceManager) retrieveChunkFromDisk(key string) *indexChunk {
	var chunkPtr *indexChunk = nil
	for !pm.chunkPool.Contains(key) {
		if pm.cacheSize.Load() >= pm.config.MaxCachedChunks {
			pm.evictNotPendingLRU()
		}
		chunkPtr = newIndexChunk(key, pm.config.IoHandler)
		pm.cacheSize.Add(1)
		pm.chunkPool.Set(key, wrappedIndexChunk{
			listEntry: pm.accessList.InsertFront(key),
			chunk:     chunkPtr,
		})
		if chunkPtr.pendingWriteBack {
			pm.pendingSync.Push(key)
		}
	}
	return chunkPtr
}

func (pm *PersistenceManager) retrieveChunk(key string) *indexChunk {
	chunk := pm.retrieveChunkFromCache(key)
	if chunk == nil {
		return pm.retrieveChunkFromDisk(key)
	}
	if chunk.termTrackers.Size() > int(pm.config.MaxChunkSize) {
		new_chunk := chunk.split()
		pm.chunkPool.Set(new_chunk.chunkKey, wrappedIndexChunk{
			listEntry: pm.accessList.InsertFront(new_chunk.chunkKey),
			chunk:     new_chunk,
		})
		pm.cacheSize.Add(1)
		pm.pendingSync.Push(new_chunk.chunkKey)
	}
	return chunk
}

func (pm *PersistenceManager) IterateTerms(term string) iter.Seq[core.TermTracker] {
	return func(yield func(core.TermTracker) bool) {
		chunk := pm.retrieveChunk("term-" + term)
		for chunk != nil && chunk.termTrackers.Size() > 0 {
			for term := range chunk.iterate() {
				if !yield(term) {
					break
				}
			}
			chunk = pm.retrieveChunk(chunk.nextChunkKey)
		}
	}
}
