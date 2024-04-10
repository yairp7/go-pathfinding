package pf

type navGraphEdge[T comparable] struct {
	node   T
	weight float64
}

type navGraph[T comparable] struct {
	nodes map[T][]navGraphEdge[T]
}

func newNavigationGraph[T comparable]() navGraph[T] {
	return navGraph[T]{nodes: make(map[T][]navGraphEdge[T])}
}

func (g *navGraph[T]) addEdge(origin, destination T, weight float64) {
	g.nodes[origin] = append(g.nodes[origin], navGraphEdge[T]{node: destination, weight: weight})
	g.nodes[destination] = append(g.nodes[destination], navGraphEdge[T]{node: origin, weight: weight})
}

func (g *navGraph[T]) getEdges(node T) []navGraphEdge[T] {
	return g.nodes[node]
}
