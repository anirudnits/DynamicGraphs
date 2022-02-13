package eulertourtree

type LLNode struct {
	value int
	prev  *LLNode
	next  *LLNode
}

type EulerTourInfo struct {
	firstInstance map[int]*LLNode
	lastInstance  map[int]*LLNode
}

func createEulerTourInfo(n int) *EulerTourInfo {
	ret := &EulerTourInfo{}

	ret.firstInstance = make(map[int]*LLNode)
	ret.lastInstance = make(map[int]*LLNode)

	return ret
}

func createLLNode(value int) *LLNode {
	return &LLNode{value, nil, nil}
}

func concatenateLL(fstHead, fstTail, sndHead, sndTail *LLNode) (*LLNode, *LLNode) {
	if fstTail == nil && sndHead == nil {
		// Both lists are nil so the return an empty nil list
		return nil, nil
	}

	if fstTail == nil || sndHead == nil {
		if fstTail != nil {
			return fstHead, fstTail
		}

		return sndHead, sndTail
	}

	fstTail.next = sndHead
	sndHead.prev = fstTail

	return fstHead, sndTail
}

func removeNodeFront(headNode, tailNode *LLNode) (*LLNode, *LLNode) {
	if headNode == nil {
		panic("Cannot remove node with nil head")
	}

	newHead := headNode.next
	if newHead != nil {
		newHead.prev = nil
	}

	return newHead, tailNode
}

func removeNodeBack(headNode, tailNode *LLNode) (*LLNode, *LLNode) {
	if tailNode == nil {
		panic("Cannot remove node from nil tail")
	}

	newTail := tailNode.prev
	if newTail != nil {
		newTail.next = nil
	}

	return headNode, newTail
}

func deleteLinkPrev(node *LLNode) {
	if node.prev == nil {
		// Prev link is not present, don't need to do anything here
		return
	}

	prevNode := node.prev

	node.prev = nil
	prevNode.next = nil
}

func deleteLinkNext(node *LLNode) {
	if node.next == nil {
		// There's no next link, so don't need to do anything
		return
	}

	nextNode := node.next

	node.next = nil
	nextNode.prev = nil
}

func searchLLLeft(node *LLNode, value int) bool {
	for node != nil {
		if node.value == value {
			return true
		}

		node = node.prev
	}

	return false
}

func searchLLRight(node *LLNode, value int) bool {
	for node != nil {
		if node.value == value {
			return true
		}

		node = node.next
	}

	return false
}

func searchLL(node *LLNode, value int) bool {
	return searchLLLeft(node, value) || searchLLRight(node, value)
}

// Traverses the connected component using an Euler Tour and returns the
// head and tail of the linked list created to store the traversal.
// Time Complexity: O(v), v is the size of the subtree rooted at vertex
func eulerTour(vertex int, graph [][]int, helper *EulerTourInfo, visited []bool) (*LLNode, *LLNode) {
	visited[vertex] = true

	vertexNode := createLLNode(vertex)

	helper.firstInstance[vertex] = vertexNode
	head, tail := vertexNode, vertexNode

	for _, neighbor := range graph[vertex] {
		if !visited[neighbor] {
			neighbor_head, neighbor_tail := eulerTour(neighbor, graph, helper, visited)

			vertexEndNode := createLLNode(vertex)
			neighbor_head, neighbor_tail = concatenateLL(
				neighbor_head,
				neighbor_tail,
				vertexEndNode,
				vertexEndNode,
			)

			head, tail = concatenateLL(head, tail, neighbor_head, neighbor_tail)
		}
	}

	helper.lastInstance[vertex] = tail
	return head, tail
}

// Time Complexity = O(h), h is the height of the tree
func getHeadLL(node *LLNode, helper *EulerTourInfo) *LLNode {
	for node.prev != nil {
		node = helper.firstInstance[node.prev.value]
	}

	return node
}

// Time Complexity = O(h), h is the height of the tree
func nearestNextNodeWithValue(node *LLNode, value int, helper *EulerTourInfo) *LLNode {
	for node != nil && node.next != nil && node.next.value != value {
		node = helper.lastInstance[node.next.value]
	}

	if node == nil {
		return nil
	} else {
		return node.next
	}
}

// Time Complexity = O(h), h is the height of the tree
func nearestPrevNodeWithValue(node *LLNode, value int, helper *EulerTourInfo) *LLNode {
	for node != nil && node.prev != nil && node.prev.value != value {
		node = helper.firstInstance[node.prev.value]
	}

	if node == nil {
		return nil
	} else {
		return node.prev
	}
}

// Time Complexity = O(h)
func reRoot(newRoot int, helper *EulerTourInfo) (*LLNode, *LLNode) {
	firstInstance, lastInstance := helper.firstInstance[newRoot], helper.lastInstance[newRoot]

	if firstInstance.prev == nil && lastInstance.next == nil {
		// This is already the root of the euler tour so don't need to do anything
		return firstInstance, lastInstance
	}

	if firstInstance.prev == nil || lastInstance.next == nil {
		panic("There's something wrong in the euler tour representation")
	}

	// cut out the section of newRoot's subtree
	// This would split the suler tour into 3 sections
	// e1 : the section before the first instance of newRoot
	// v: the section with newRoot's subtree
	// e2: the section after the last instance of newRoot

	e1Head, e1Tail := getHeadLL(firstInstance, helper), firstInstance.prev
	e2Head, e2Tail := lastInstance.next, helper.lastInstance[e1Head.value]

	deleteLinkPrev(firstInstance)
	deleteLinkNext(lastInstance)

	e1Head, e1Tail = removeNodeFront(e1Head, e1Tail)

	// Update the first and last instance of the affected vertices
	// i.e e2Tail and e1Head, considering the assumptions
	// e1Head.value(before the deletion) == e2Tail.value and e1Tail.value = e2Head.value

	helper.firstInstance[e1Tail.value] = e2Head
	helper.firstInstance[e1Tail.value] = e1Tail

	// e2Tail will be the current root before the rerooting
	helper.firstInstance[e2Tail.value] = nearestNextNodeWithValue(e2Head, e2Tail.value, helper)
	helper.lastInstance[e2Tail.value] = nearestPrevNodeWithValue(e1Tail, e1Tail.value, helper)
	if helper.lastInstance[e2Tail.value] == nil {
		helper.lastInstance[e2Tail.value] = e2Tail
	}

	// concatenate e2 and e1 sections
	aggregatedSectionHead, aggregatedSectionTail := concatenateLL(
		e2Head,
		e2Tail,
		e1Head,
		e1Tail,
	)

	vertexNode := createLLNode(newRoot)
	aggregatedSectionHead, aggregatedSectionTail = concatenateLL(
		aggregatedSectionHead,
		aggregatedSectionTail,
		vertexNode,
		vertexNode,
	)

	newTreeHead, newTreeTail := concatenateLL(firstInstance, lastInstance, aggregatedSectionHead, aggregatedSectionTail)
	helper.lastInstance[newRoot] = newTreeTail

	return newTreeHead, newTreeTail
}

// Adds the link(edge) u-v, assuming that u-v are not currently connected
// Time Complexity: O(h)
func Link(u, v int, helper *EulerTourInfo) {
	firstLLHead, firstLLTail := reRoot(u, helper)
	sndLLHead, sndLLTail := reRoot(v, helper)

	connectedHead, connectedTail := concatenateLL(firstLLHead, firstLLTail, sndLLHead, sndLLTail)

	vertexNode := createLLNode(u)
	_, connectedTail = concatenateLL(connectedHead, connectedTail, vertexNode, vertexNode)
	helper.lastInstance[u] = connectedTail
}

func cutSectionRight(cutNode *LLNode) {
	if cutNode.next != nil {
		cutNode.next.prev = nil
		cutNode.next = nil
	}
}

func cutSectionLeft(cutNode *LLNode) {
	if cutNode.prev != nil {
		cutNode.prev.next = nil
		cutNode.prev = nil
	}
}

// Cuts the link(edge) u-v, assuming that u-v exists
// Time Complexity: O(h)
func Cut(u, v int, helper *EulerTourInfo) {
	reRoot(u, helper)

	vHead, vTail := helper.firstInstance[v], helper.lastInstance[v]

	// Split the linked list into 3 parts :
	// e1: the list before the first instance of v
	// v: the subtree of v
	// e2: the list after the last instance of v
	e1Head, e1Tail := helper.firstInstance[u], vHead.prev
	e2Head, e2Tail := vTail.next, helper.lastInstance[u]

	// cut out the sections
	cutSectionLeft(vHead)
	cutSectionRight(vTail)

	e1Head, e1Tail = removeNodeBack(e1Head, e1Tail)

	helper.firstInstance[u], helper.lastInstance[u] = concatenateLL(e1Head, e1Tail, e2Head, e2Tail)
}

// Returns a bool indicating whether u and v are connected or not
// Time Complexity: O(n), n is the number of vertices in the graph
func Is_Connected(u int, v int, helper *EulerTourInfo) bool {
	return searchLL(helper.firstInstance[u], v)
}

// Initiates the initial euler tour of the graph
// Returns the EulerTourInfo, which stores information about the first and last instances
// of the vertices. Also this information is used in and is a reqd parameter in other functions
// Time Complexity: O(n), n is the number of vertices
func InitiateEulerTree(graph [][]int) *EulerTourInfo {
	n := len(graph)

	helper := createEulerTourInfo(n)
	visited := make([]bool, n)

	for vertex := 0; vertex < n; vertex++ {
		if !visited[vertex] {
			eulerTour(vertex, graph, helper, visited)
		}
	}

	return helper
}
