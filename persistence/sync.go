package persistence

import (
	"iter"
	"quinto/misc"
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

func (syncseg *synchronizedSegment) iterator() iter.Seq[TermTracker] {
	return func(yield func(TermTracker) bool) {
		syncseg.mutex.RLock()
		defer syncseg.mutex.RUnlock()
		syncseg.iterator()(yield)
	}
}

func (syncseg *synchronizedSegment) add(tracker TermTracker) {
	syncseg.mutex.Lock()
	defer syncseg.mutex.Unlock()
	syncseg.seg.add(tracker)
}
