package utils

type Stack[T any] struct {
	size  int
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: []T{},
	}
}

func (s *Stack[T]) Push(item T) {
	if s.size < len(s.items) {
		s.items[s.size] = item
	} else {
		s.items = append(s.items, item)
	}
	s.size++
}

func (s *Stack[T]) Pop() T {
	item := s.items[s.size-1]
	s.items = s.items[:s.size-1]
	s.size--
	return item
}

func (s *Stack[T]) Size() int {
	return s.size
}
