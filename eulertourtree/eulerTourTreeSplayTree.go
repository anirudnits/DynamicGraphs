package eulertourtree

type Edge struct {
	from int
	to   int
}

func isElement(edge Edge, mapper map[Edge]*treeNode) bool {
	_, contains := mapper[edge]

	return contains
}

type EulerTourInfoSplay struct {
	// TODO: examine whether I could map an edge pointer to a treenode pointer
	edgeNodeMapper map[Edge]*treeNode
	visited        []bool
}

func createEdge(from, to int) *Edge {
	return &Edge{from, to}
}

func createEulerTourInfoSplay(n int) *EulerTourInfoSplay {
	ret := &EulerTourInfoSplay{}

	ret.visited = make([]bool, n)
	ret.edgeNodeMapper = make(map[Edge]*treeNode)

	return ret
}

func eulerTourSplay(vertex int, graph [][]int, helper *EulerTourInfoSplay, prevRoot *treeNode) {
	helper.visited[vertex] = true

	vertexEdge := createEdge(vertex, vertex)
	vertexNode := createNode(vertexEdge)
	root := insert(vertexNode, prevRoot)

	helper.edgeNodeMapper[*vertexEdge] = vertexNode

	for _, neighbor := range graph[vertex] {
		if !helper.visited[neighbor] {
			// create an edge for vertex-neighbor edge and
			// insert into the tree
			goEdge := createEdge(vertex, neighbor)
			goNode := createNode(goEdge)
			root = insert(goNode, root)
			helper.edgeNodeMapper[*goEdge] = goNode

			eulerTourSplay(neighbor, graph, helper, root)

			// create an edge for neighbor-vertex edge and
			// insert into the tree
			returnEdge := createEdge(neighbor, vertex)
			returnNode := createNode(returnEdge)
			root = insert(returnNode, root)
			helper.edgeNodeMapper[*returnEdge] = returnNode

		}
	}
}

// Re-roots the tree on the edges {u, u}
// I have some lemmas on how the representation remains
// consistent after this operation. I will prove them
// and provide them in a different file.
// Amortised time complexity: O(logn)
func reRootSplay(newRoot int, helper *EulerTourInfoSplay) {
	newRootNode := helper.edgeNodeMapper[Edge{newRoot, newRoot}]

	// split the tree into 2 parts
	// A: the part of the tree before edge {newRoot, newRoot}
	// B: remaining tree
	A, B := leftSplitTreeAtNode(newRootNode)
	joinTrees(A, B)
}

// Connects vertices u and v
// Amortised time complexity: O(logn)
func link(u, v int, helper *EulerTourInfoSplay) {
	if isConnectedSplay(u, v, helper) {
		// the two vertices are already connected
		// This implementation is based on the idea that
		// this never happens. Just return
		// TODO: Provide an error/warning on this condition

		return
	}
	uEdge := helper.edgeNodeMapper[Edge{u, u}]
	vEdge := helper.edgeNodeMapper[Edge{v, v}]

	reRootSplay(u, helper)
	reRootSplay(v, helper)

	// Add edge {u, v} to the rightmost node of uTree
	newEdgeUV := createEdge(u, v)
	newNodeUV := createNode(newEdgeUV)
	insertAtRightMost(newNodeUV, uEdge)
	helper.edgeNodeMapper[*newEdgeUV] = newNodeUV

	// Add edge {v, u} to the rightmost node of vTree
	newEdgeVU := createEdge(v, u)
	newNodeVU := createNode(newEdgeVU)
	insertAtRightMost(newNodeVU, uEdge)
	helper.edgeNodeMapper[*newEdgeVU] = newNodeVU

	joinTrees(uEdge, vEdge)
}

func cut(u, v int, helper *EulerTourInfoSplay) {
	// Check if u-v and v-u are present

	edgeUV, edgeVU := Edge{u, v}, Edge{v, u}
	edgeU, edgeV := Edge{u, u}, Edge{v, v}

	if !(isElement(edgeUV, helper.edgeNodeMapper) || isElement(edgeVU, helper.edgeNodeMapper)) {
		// TODO: provide an error/warning saying that u-v is not an edge currectly
		return
	}

	if !isElement(edgeUV, helper.edgeNodeMapper) || !isElement(edgeVU, helper.edgeNodeMapper) {
		// raise an error because this condition should never happen
		panic("Grpah is inconsistent")
	}

	// Split the tree into 3 sections
	// A: Before edge {u, v}
	// B: Between edge {u, v} and {v, u}
	// C: After edge {v, u}

	A, B1 := leftSplitTreeAtNode(helper.edgeNodeMapper[edgeUV])
	B2, C := rightSplitTreeAtNode(helper.edgeNodeMapper[edgeVU])

	joinTrees(B1, B2) // to get B
	joinTrees(A, C)

	// Delete edges {u, v} and {v, u}
	delete(helper.edgeNodeMapper, Edge{u, v})
	delete(helper.edgeNodeMapper, Edge{v, u})

	splayToRoot(helper.edgeNodeMapper[edgeU])
	splayToRoot(helper.edgeNodeMapper[edgeV])

}

// Checks if vertices u and v are connected
// Amortised time complexity: O(logn)
func isConnectedSplay(u, v int, helper *EulerTourInfoSplay) bool {
	uRoot := getRoot(helper.edgeNodeMapper[Edge{u, u}])
	vRoot := getRoot(helper.edgeNodeMapper[Edge{v, v}])

	return (uRoot == vRoot)
}

// Initiates the initial euler tour of the graph
// Returns the EulerTourInfo, which stores information about the first and last instances
// of the vertices. Also this information is used in and is a reqd parameter in other functions
// Time Complexity: O(n), n is the number of vertices
func InitiateEulerTreeSplay(graph [][]int) *EulerTourInfoSplay {
	n := len(graph)

	helper := createEulerTourInfoSplay(n)

	for vertex := 0; vertex < n; vertex++ {
		// TODO: effective way to splay the tree on random nodes
		// Now as each edge is appended the tree is splayed at that edge
		if !helper.visited[vertex] {
			eulerTourSplay(vertex, graph, helper, nil)
		}
	}

	return helper
}
