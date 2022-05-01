package eulertourtreesplay

type treeNode struct {
	value     interface{}
	dummyComp int // dummy comparison value used in insertion
	left      *treeNode
	right     *treeNode
	parent    *treeNode
}

func createNode(value interface{}) *treeNode {
	return &treeNode{
		value:     value,
		dummyComp: 0,
		left:      nil,
		right:     nil,
		parent:    nil,
	}
}

// Right rotate on node, lifting its left child
func rightRotate(node *treeNode) {
	toBeliftedNode := node.left
	nodeParent := node.parent

	node.left = toBeliftedNode.right
	toBeliftedNode.right = node

	node.parent = toBeliftedNode
	toBeliftedNode.parent = nodeParent
}

// Left rotate on node, lifting its right child
func leftRotate(node *treeNode) {
	toBeliftedNode := node.right
	nodeParent := node.parent

	node.right = toBeliftedNode.left
	toBeliftedNode.left = node

	node.parent = toBeliftedNode
	toBeliftedNode.parent = nodeParent
}

func isLeaf(node *treeNode) bool {
	return node.left == nil && node.right == nil
}

func isRoot(node *treeNode) bool {
	return node.parent == nil
}

// Time Complexity: O(h), h is the height of the tree
// For splay tree, its amortized O(logn), n is the number of nodes
func getRightMostNode(root *treeNode) *treeNode {
	for root.right != nil {
		root = root.right
	}

	return root
}

// Time Complexity: O(h), h is the height of the tree
// For splay tree, its amortized O(logn): n is the number of nodes
func getRoot(node *treeNode) *treeNode {
	for !isRoot(node) {
		node = node.parent
	}

	return node
}

// Time Complexity: O(h), h is the height of the tree
// For splay tree, its amortized O(logn), n is the number of nodes
func getLeftMostNode(root *treeNode) *treeNode {
	for root.left != nil {
		root = root.left
	}

	return root
}

// Time Complexity: O(h), h is the height of the tree
// For splay tree, its amortized O(logn), n is the number of nodes
func inorderSuccessor(node *treeNode) *treeNode {
	if node.right != nil {
		return getLeftMostNode(node.right)
	}

	for node.parent != nil && node.parent.left != node {
		node = node.parent
	}

	return node.parent
}

// Time Complexity: O(h), h is the height of the tree
// For splay tree, its amortized O(logn), n is the number of nodes
func inorderPredecessor(node *treeNode) *treeNode {
	if node.left != nil {
		return getRightMostNode(node.left)
	}

	for node.parent != nil && node.parent.right != node {
		node = node.parent
	}

	return node.parent
}

// The function takes the node and returns one of the mutiple splay conditions that might apply:
// 1. node is left child of parent and parent is left child of grandparent
// 2. node is right child of parent and parent is right child of grandparent
// 3. node is right child of parent and parent is left child of grandparent
// 4. node is left child of parent and parent is right child of grandparent
func checkSplayCondition(node *treeNode) int {
	parent, grandParent := node.parent, node.parent.parent

	if node == parent.left && parent == grandParent.left {
		return 1
	}

	if node == parent.right && parent == grandParent.right {
		return 2
	}

	if node == parent.right && parent == grandParent.left {
		return 3
	}

	return 4
}

func splay(node *treeNode) {
	if node.parent == nil {
		// node is already at the root, so just return
		return
	}

	if node.parent.parent == nil {
		// the node has no grandparent, so just one left/right rotate is required
		if node == node.parent.left {
			rightRotate(node.parent)
		} else {
			leftRotate(node.parent)
		}

		return
	}

	splayCondition := checkSplayCondition(node)
	parent, grandParent := node.parent, node.parent.parent

	switch splayCondition {
	case 1:
		leftRotate(grandParent)
		leftRotate(parent)
	case 2:
		rightRotate(grandParent)
		rightRotate(parent)
	case 3:
		leftRotate(parent)
		rightRotate(grandParent)
	case 4:
		rightRotate(parent)
		leftRotate(grandParent)
	}
}

func splayToRoot(node *treeNode) {
	for node.parent != nil {
		splay(node)
	}
}

func insertAtRightMost(node *treeNode, root *treeNode) {
	rightMostNode := getRightMostNode(root)

	rightMostNode.right = node
	node.parent = rightMostNode
}

func insert(node *treeNode, root *treeNode) *treeNode {
	if root == nil {
		return node
	}

	insertAtRightMost(node, root)

	splayToRoot(node)

	return node
}

func deleteTreeNode(node *treeNode) {
	splayToRoot(node)

	leftSubtree := node.left
	rightSubtree := node.right

	node.left = nil
	node.right = nil

	leftSubtree.parent = nil
	rightSubtree.parent = nil

	rightMostNodeLeftSubtree := getRightMostNode(leftSubtree)
	splayToRoot(rightMostNodeLeftSubtree)

	rightMostNodeLeftSubtree.right = rightSubtree
}

// Spliting the tree at node, into 2 parts
// A: The part of the tree before node(including the node) and
// B: the part of the tree after node
// Amortized Time Complexity: O(logn)
// Returns the roots of the two trees formed in order (A, B)
func rightSplitTreeAtNode(node *treeNode) (*treeNode, *treeNode) {
	successor := inorderSuccessor(node)
	splayToRoot(successor)

	// cut the left link to split the tree node
	node.parent = nil
	successor.left = nil

	return node, successor
}

// Spliting the tree at node, into 2 parts
// A: The part of the tree before node and
// B: the part of the tree after node(including the node)
// Amortized Time Complexity: O(logn)
// Returns the roots of the two trees formed in order (A, B)
func leftSplitTreeAtNode(node *treeNode) (*treeNode, *treeNode) {
	predecessor := inorderPredecessor(node)
	splayToRoot(predecessor)

	// cut the right link to split the tree
	node.parent = nil
	predecessor.right = nil

	return predecessor, node
}

// Joining two trees, assuming all nodes in root2 are larger
// than nodes in root1
func joinTrees(root1, root2 *treeNode) {
	maxElementTree1 := getRightMostNode(root1)
	splayToRoot(maxElementTree1)

	maxElementTree1.right = root2
	root2.parent = maxElementTree1
}