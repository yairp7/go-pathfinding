package utils

type Queue[T any] struct {
	size  int
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: []T{},
	}
}

func (q *Queue[T]) Push(item T) {
	if q.size < len(q.items) {
		q.items[q.size] = item
	} else {
		q.items = append(q.items, item)
	}
	q.size++
}

func (q *Queue[T]) Pop() T {
	item := q.items[0]
	q.items = q.items[1:]
	q.size--
	return item
}

func (q *Queue[T]) Size() int {
	return q.size
}
