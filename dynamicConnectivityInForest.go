package dynamicgraphs

import (
	"dynamicgraphs/eulertourtree"
)

type Query struct {
	operation string
	vertex1   int
	vertex2   int
}

func processQueries(queries []Query, eulerTourInfo *eulertourtree.EulerTourInfo) []bool {
	connectivity := make([]bool, 0)

	for _, query := range queries {
		if query.operation == "check" {
			connectivityCheck := eulertourtree.Is_Connected(
				query.vertex1,
				query.vertex2,
				eulerTourInfo,
			)
			connectivity = append(connectivity, connectivityCheck)
		} else if query.operation == "link" {
			eulertourtree.Link(
				query.vertex1,
				query.vertex2,
				eulerTourInfo,
			)
		} else if query.operation == "cut" {
			eulertourtree.Cut(
				query.vertex1,
				query.vertex2,
				eulerTourInfo,
			)
		}
	}

	return connectivity
}

// Processes the list of queries sequentially and records the check
// connectivity operation as eith true/false and returns this list
// This is a special case of the dynmaic connectivity problem, wherein cut operations
// are assumed to split a tree in two and each link operation is assumed to connect
// two previously unconnected vertices and not close any loops
func CheckDynamicConnectivity(queries []Query, graph [][]int) []bool {
	eulerTourInfo := eulertourtree.InitiateEulerTree(graph)

	return processQueries(queries, eulerTourInfo)
}
