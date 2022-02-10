package dynamicgraphs

import (
	"reflect"
	"testing"
)

func TestDynamicConnectivityInForest(t *testing.T) {
	testCases := []struct {
		graph    [][]int
		queries  []Query
		expected []bool
	}{
		{
			graph: [][]int{{1, 2}, {0}, {0}},
			queries: []Query{
				{"check", 0, 1},
				{"check", 0, 2},
				{"check", 1, 2},
				{"cut", 0, 1},
				{"check", 0, 1},
				{"check", 1, 2},
				{"check", 0, 2},
			},
			expected: []bool{true, true, true, false, false, true},
		},
	}

	for _, testCase := range testCases {
		t.Run("Running tests", func(t *testing.T) {
			calculatedResult := CheckDynamicConnectivity(testCase.queries, testCase.graph)
			if !reflect.DeepEqual(testCase.expected, calculatedResult) {
				t.Fatalf("Calculated and the expected result don't match")
			}
		})
	}
}
