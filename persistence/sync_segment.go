package persistence

import (
	"iter"
	"quinto/misc"
	"unsafe"
)

type synchronizedSegment struct {
	seg   segment
	mutex misc.WritersFirstRWMutex
}

func newSynchronizedSegment() *synchronizedSegment {
	return &synchronizedSegment{
		seg:   *newSegment(),
		mutex: *misc.NewWritersFirstRWMutex(),
	}
}

func (syncseg *synchronizedSegment) estimateSize() int64 {
	syncseg.mutex.RLock()
	defer syncseg.mutex.RUnlock()
	return int64(unsafe.Sizeof(syncseg.mutex)) + syncseg.seg.estimateSize()
}

func (syncseg *synchronizedSegment) iterator() iter.Seq[misc.TermTracker] {
	return func(yield func(misc.TermTracker) bool) {
		syncseg.mutex.RLock()
		defer syncseg.mutex.RUnlock()
		syncseg.seg.iterator()(yield)
	}
}

func (syncseg *synchronizedSegment) add(tracker misc.TermTracker) {
	syncseg.mutex.Lock()
	defer syncseg.mutex.Unlock()
	syncseg.seg.add(tracker)
}
