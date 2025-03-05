/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

A synchronizedSegment is just a convenient wrapper around a segment that provides
synchronization primitives to ensure that the segment can be safely accessed by
multiple goroutines. This is useful when the segment is shared between multiple
goroutines, such as when it is stored in a cache.
==================================================================================*/

package persistence

import (
	"iter"
	"quinto/misc"
	"unsafe"
)

type synchronizedSegment struct {
	underlyng segment
	mutex     misc.WritersFirstRWMutex
}

func newSynchronizedSegment(seg *segment) *synchronizedSegment {
	if seg == nil {
		seg = newSegment()
	}
	return &synchronizedSegment{
		underlyng: *seg,
		mutex:     *misc.NewWritersFirstRWMutex(),
	}
}

func (syncseg *synchronizedSegment) estimateSize() int64 {
	syncseg.mutex.RLock()
	defer syncseg.mutex.RUnlock()
	return int64(unsafe.Sizeof(syncseg.mutex)) + syncseg.underlyng.estimateSize()
}

func (syncseg *synchronizedSegment) iterator() iter.Seq[misc.TermTracker] {
	return func(yield func(misc.TermTracker) bool) {
		syncseg.mutex.RLock()
		defer syncseg.mutex.RUnlock()
		syncseg.underlyng.iterator()(yield)
	}
}

func (syncseg *synchronizedSegment) add(tracker misc.TermTracker) {
	syncseg.mutex.Lock()
	defer syncseg.mutex.Unlock()
	syncseg.underlyng.add(tracker)
}
