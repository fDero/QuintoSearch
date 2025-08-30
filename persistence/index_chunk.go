/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains the implementation of `indexChunk`, which is wrapper around a
sorted array of `core.TermTracker` objects. It is supposed to be the smallest unit
in which a full inverted list can be concurrently accessed. A full inverted list is
composed of one or more index chunks. An `indexChunk` must be read and written
on disk. Multiple readers can read from it concurrently, but only one writer can
update it. After an update, the `indexChunk` must be written back to disk.
==================================================================================*/

package persistence

import (
	"fmt"
	"iter"
	"quinto/core"
	"quinto/data"
	"strconv"
)

type indexChunk struct {
	termTrackers     data.SortedArray[core.TermTracker]
	chunkKey         string
	nextChunkKey     string
	pendingWriteBack bool
	splitCounter     uint64
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

func newIndexChunk(chunkKey string, handler diskHandler) *indexChunk {
	chunk := &indexChunk{
		termTrackers:     newSortedArrayOfTermTrackers(),
		chunkKey:         chunkKey,
		nextChunkKey:     "",
		pendingWriteBack: false,
		handler:          handler,
		splitCounter:     0,
		rwMutex:          core.NewWritersFirstRWMutex(),
	}
	reader, exists := handler.getReader(chunkKey)
	if !exists || reader == nil {
		return chunk
	}
	errors := [3]error{}
	var splitCounterString string = ""
	chunk.chunkKey, errors[0] = decodeStringFromDisk(reader)
	chunk.nextChunkKey, errors[1] = decodeStringFromDisk(reader)
	splitCounterString, errors[2] = decodeStringFromDisk(reader)
	chunk.splitCounter, _ = strconv.ParseUint(splitCounterString, 10, 64)
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
	encodeStringToDisk(writer, chunk.chunkKey)
	encodeStringToDisk(writer, chunk.nextChunkKey)
	encodeStringToDisk(writer, fmt.Sprint(chunk.splitCounter))
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

func (chunk *indexChunk) split() *indexChunk {
	chunk.rwMutex.Lock()
	chunk.splitCounter++
	defer chunk.rwMutex.Unlock()
	newChunk := &indexChunk{
		termTrackers:     newSortedArrayOfTermTrackers(),
		chunkKey:         chunk.chunkKey + "-" + fmt.Sprint(chunk.splitCounter),
		nextChunkKey:     chunk.nextChunkKey,
		pendingWriteBack: false,
		handler:          chunk.handler,
		rwMutex:          core.NewWritersFirstRWMutex(),
	}
	chunk.nextChunkKey = newChunk.chunkKey
	counter := 0
	for tracker := range chunk.iterate() {
		if counter >= chunk.termTrackers.Size()/2 {
			newChunk.termTrackers.Insert(tracker)
			chunk.termTrackers.Remove(tracker)
		}
		counter++
	}
	return newChunk
}
