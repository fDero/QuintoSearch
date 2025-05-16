package data

type Heap[T any] struct {
	orderingPredicate func(a, b T) bool
	storage           []T
}

func NewHeap[T any](orderingPredicate func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		orderingPredicate: orderingPredicate,
		storage:           make([]T, 0),
	}
}

func (h *Heap[T]) Push(value T) {
	h.storage = append(h.storage, value)
	h.siftUp()
}

func (h *Heap[T]) Peek() (T, bool) {
	if h.Size() == 0 {
		var zero T
		return zero, false
	}
	return h.storage[0], true
}

func (h *Heap[T]) Pop() (T, bool) {
	if h.Size() == 0 {
		var zero T
		return zero, false
	}
	n := h.Size()
	h.swap(0, n-1)
	popped := h.storage[n-1]
	h.storage = h.storage[:n-1]
	h.siftDown(0, n-2)
	return popped, true
}

func (h *Heap[T]) Size() int {
	return len(h.storage)
}

func (h *Heap[T]) siftDown(currentIdx int, endIdx int) {
	leftChildIdx := currentIdx*2 + 1
	for leftChildIdx <= endIdx {

		rightChildIdx := currentIdx*2 + 2
		if rightChildIdx > endIdx {
			rightChildIdx = -1
		}

		idxToSwap := leftChildIdx
		if h.compareAtIndex(rightChildIdx, leftChildIdx) {
			idxToSwap = rightChildIdx
		}

		if h.compareAtIndex(idxToSwap, currentIdx) {
			h.swap(idxToSwap, currentIdx)
			currentIdx = idxToSwap
			leftChildIdx = currentIdx*2 + 1

		} else {
			return
		}
	}
}

func (h *Heap[T]) siftUp() {
	currentIdx := h.Size() - 1
	parentIdx := (currentIdx - 1) / 2
	for h.compareAtIndex(currentIdx, parentIdx) {
		h.swap(currentIdx, parentIdx)
		currentIdx = parentIdx
		parentIdx = (currentIdx - 1) / 2
	}
}

func (h *Heap[T]) compareAtIndex(i, j int) bool {
	return i >= 0 && j >= 0 &&
		i < len(h.storage) &&
		j < len(h.storage) &&
		h.orderingPredicate(h.storage[i], h.storage[j])
}

func (h *Heap[T]) swap(i, j int) {
	h.storage[i], h.storage[j] =
		h.storage[j], h.storage[i]
}
