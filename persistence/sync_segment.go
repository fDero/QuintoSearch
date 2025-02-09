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

func newSynchronizedSegment() *synchronizedSegment {
	return &synchronizedSegment{
		underlyng: *newSegment(),
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
