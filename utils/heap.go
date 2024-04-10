package utils

import hp "container/heap"

type HeapNode[T any] struct {
	TotalWeight float64
	Data        T
}

type HeapNodes[T any] []HeapNode[T]

func (h HeapNodes[T]) Len() int           { return len(h) }
func (h HeapNodes[T]) Less(i, j int) bool { return h[i].TotalWeight < h[j].TotalWeight }
func (h HeapNodes[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *HeapNodes[T]) Push(x any) {
	*h = append(*h, x.(HeapNode[T]))
}

func (h *HeapNodes[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type MinHeap[T any] struct {
	Values *HeapNodes[T]
}

func NewMinHeap[T any]() *MinHeap[T] {
	return &MinHeap[T]{Values: &HeapNodes[T]{}}
}

func (h *MinHeap[T]) Push(p HeapNode[T]) {
	hp.Push(h.Values, p)
}

func (h *MinHeap[T]) Pop() HeapNode[T] {
	n := hp.Pop(h.Values)
	return n.(HeapNode[T])
}
