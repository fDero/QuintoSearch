package data

type Set[T comparable] struct {
	storage map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		storage: make(map[T]bool),
	}
}

func (s *Set[T]) InsertOne(value T) {
	s.storage[value] = true
}

func (s *Set[T]) InsertAll(other *Set[T]) {
	for value := range other.storage {
		s.storage[value] = true
	}
}

func (s *Set[T]) Contains(value T) bool {
	flag, exists := s.storage[value]
	return exists && flag
}

func (s *Set[T]) Size() int {
	return len(s.storage)
}
