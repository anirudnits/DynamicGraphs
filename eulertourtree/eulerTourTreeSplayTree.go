package eulertourtree

type edge struct {
	from int
	to   int
}

type SplayEulerTourInfo struct {
	// TODO: examine whether I could map an edge pointer to a treenode pointer
	edgeNodeMapper map[edge]*treeNode
	visited        []bool
}

func createEdge(from, to int) *edge {
	return &edge{from, to}
}
