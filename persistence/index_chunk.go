/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains the implementation of `indexChunk`, which is wrapper around a
sorted array of `core.TermTracker` objects. It is supposed to be the smallest unit
in which a full inverted list is stored. An `indexChunk` must be read and written
on disk. Multiple readers can read from it concurrently, but only one writer can
update it. After an update, the `indexChunk` must be written back to disk.
==================================================================================*/

package persistence

import (
	"iter"
	"quinto/core"
	"quinto/data"
)

type indexChunk struct {
	termTrackers     data.SortedArray[core.TermTracker]
	termAsText       string
	chunkKey         string
	nextChunkKey     string
	pendingWriteBack bool
	handler          diskHandler
	rwMutex          core.ReadWriteMutex
}

func panicWhenSomeErrorsOccurred(err []error) {
	for _, e := range err {
		if e != nil {
			panic(e)
		}
	}
}

func newSortedArrayOfTermTrackers() data.SortedArray[core.TermTracker] {
	equalityPredicate := func(this, other core.TermTracker) bool {
		return this.DocId == other.DocId &&
			this.Position == other.Position
	}
	orderingPredicate := func(this, other core.TermTracker) bool {
		return this.DocId < other.DocId ||
			(this.DocId == other.DocId && this.Position < other.Position)
	}
	return *data.NewSortedArray(orderingPredicate, equalityPredicate)
}

func newIndexChunk(term string, chunkKey string, handler diskHandler) *indexChunk {
	chunk := &indexChunk{
		termTrackers:     newSortedArrayOfTermTrackers(),
		termAsText:       term,
		chunkKey:         chunkKey,
		nextChunkKey:     "",
		pendingWriteBack: false,
		handler:          handler,
		rwMutex:          core.NewWritersFirstRWMutex(),
	}
	reader, exists := handler.getReader(chunkKey)
	if !exists || reader == nil {
		return chunk
	}
	errors := [3]error{}
	chunk.termAsText, errors[0] = decodeStringFromDisk(reader)
	chunk.chunkKey, errors[2] = decodeStringFromDisk(reader)
	chunk.nextChunkKey, errors[1] = decodeStringFromDisk(reader)
	panicWhenSomeErrorsOccurred(errors[:])
	for tracker := range iterateTermTrackersFromDisk(reader) {
		chunk.termTrackers.Insert(tracker)
	}
	return chunk
}

func (chunk *indexChunk) writeBack() {
	chunk.rwMutex.RLock()
	defer chunk.rwMutex.RUnlock()
	if !chunk.pendingWriteBack {
		return
	}
	writer, finalize, _ := chunk.handler.getWriter(chunk.chunkKey)
	defer finalize()
	encodeStringToDisk(writer, chunk.termAsText)
	encodeStringToDisk(writer, chunk.chunkKey)
	encodeStringToDisk(writer, chunk.nextChunkKey)
	encodeTermTrackersToDisk(writer, chunk.iterate())
	chunk.pendingWriteBack = false
}

func (chunk *indexChunk) iterate() iter.Seq[core.TermTracker] {
	chunk.rwMutex.RLock()
	return func(yield func(core.TermTracker) bool) {
		for tracker := range chunk.termTrackers.Iterate() {
			if !yield(tracker) {
				break
			}
		}
		chunk.rwMutex.RUnlock()
	}
}

func (chunk *indexChunk) insertIterable(termsIterator iter.Seq[core.TermTracker]) {
	chunk.rwMutex.Lock()
	defer chunk.rwMutex.Unlock()
	for term := range termsIterator {
		inserted := chunk.termTrackers.Insert(term)
		chunk.pendingWriteBack = chunk.pendingWriteBack || inserted
	}
}

func (chunk *indexChunk) removeFromDocument(docId core.DocumentId) {
	chunk.rwMutex.Lock()
	defer chunk.rwMutex.Unlock()
	predicate := func(tracker core.TermTracker) bool {
		return tracker.DocId == docId
	}
	removed := chunk.termTrackers.RemoveIf(predicate)
	chunk.pendingWriteBack = chunk.pendingWriteBack || removed
}
