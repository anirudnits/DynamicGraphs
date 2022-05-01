package dynamicgraphs

import (
	"dynamicgraphs/eulertourtree"
)

type Query struct {
	operation string
	vertex1   int
	vertex2   int
}

func processQueries(queries []Query) []bool {
	connectivity := make([]bool, 0)

	for _, query := range queries {
		if query.operation == "check" {
			connectivityCheck := eulertourtree.Is_Connected(
				query.vertex1,
				query.vertex2,
			)
			connectivity = append(connectivity, connectivityCheck)
		} else if query.operation == "link" {
			eulertourtree.Link(
				query.vertex1,
				query.vertex2,
			)
		} else if query.operation == "cut" {
			eulertourtree.Cut(
				query.vertex1,
				query.vertex2,
			)
		}
	}

	return connectivity
}

// Processes the list of queries sequentially and records the check
// connectivity operation as eith true/false and returns this list
// This is a special case of the dynamic connectivity problem, wherein cut operations
// are assumed to split a tree in two and each link operation is assumed to connect
// two previously unconnected vertices and not close any loops
func CheckDynamicConnectivity(queries []Query, graph [][]int) []bool {
	eulertourtree.InitiateEulerTree(graph)

	return processQueries(queries)
}
