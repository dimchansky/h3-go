//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_removeVertexNode_parity(t *testing.T) {
	tests := []struct {
		name        string
		numBuckets  int32
		res         int32
		edges       []struct{ from, to LatLng }
		removeIndex int
		description string
	}{
		{
			name:       "remove single node",
			numBuckets: 10,
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
			},
			removeIndex: 0,
			description: "Remove the only node from a single-node graph",
		},
		{
			name:       "remove first node from chain",
			numBuckets: 1, // Force all into same bucket to create chain
			res:        8,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
				{LatLng{Lat: 37.780, Lng: -122.420}, LatLng{Lat: 37.781, Lng: -122.421}},
				{LatLng{Lat: 37.785, Lng: -122.425}, LatLng{Lat: 37.786, Lng: -122.426}},
			},
			removeIndex: 0,
			description: "Remove first node from a chain in hash bucket",
		},
		{
			name:       "remove middle node from chain",
			numBuckets: 1, // Force all into same bucket to create chain
			res:        8,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
				{LatLng{Lat: 37.780, Lng: -122.420}, LatLng{Lat: 37.781, Lng: -122.421}},
				{LatLng{Lat: 37.785, Lng: -122.425}, LatLng{Lat: 37.786, Lng: -122.426}},
			},
			removeIndex: 1,
			description: "Remove middle node from a chain in hash bucket",
		},
		{
			name:       "remove last node from chain",
			numBuckets: 1, // Force all into same bucket to create chain
			res:        8,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
				{LatLng{Lat: 37.780, Lng: -122.420}, LatLng{Lat: 37.781, Lng: -122.421}},
				{LatLng{Lat: 37.785, Lng: -122.425}, LatLng{Lat: 37.786, Lng: -122.426}},
			},
			removeIndex: 2,
			description: "Remove last node from a chain in hash bucket",
		},
		{
			name:       "remove from distributed graph",
			numBuckets: 10,
			res:        5,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 0.0, Lng: 0.0}, LatLng{Lat: 0.001, Lng: 0.001}},
				{LatLng{Lat: 10.0, Lng: 10.0}, LatLng{Lat: 10.001, Lng: 10.001}},
				{LatLng{Lat: 20.0, Lng: 20.0}, LatLng{Lat: 20.001, Lng: 20.001}},
			},
			removeIndex: 1,
			description: "Remove node from graph with distributed hash buckets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and populate Go graph
			goGraph := &vertexGraph{
				Buckets:    make([]*vertexNode, tt.numBuckets),
				NumBuckets: tt.numBuckets,
				Size:       0,
				Res:        tt.res,
			}

			// Add all edges and collect node references
			nodes := make([]*vertexNode, len(tt.edges))
			for i, edge := range tt.edges {
				nodes[i] = addVertexNode(goGraph, &edge.from, &edge.to)
				if nodes[i] == nil {
					t.Fatalf("addVertexNode returned nil for edge %d", i)
				}
			}

			// Verify initial graph size
			expectedInitialSize := int32(len(tt.edges))
			if goGraph.Size != expectedInitialSize {
				t.Fatalf("Initial graph size: got %d, want %d", goGraph.Size, expectedInitialSize)
			}

			// Test removal with Go implementation
			nodeToRemove := nodes[tt.removeIndex]
			goResult := removeVertexNode(goGraph, nodeToRemove)

			// Test removal with C implementation (simulated)
			cGraph := &vertexGraph{
				Buckets:    make([]*vertexNode, tt.numBuckets),
				NumBuckets: tt.numBuckets,
				Size:       0,
				Res:        tt.res,
			}

			// Recreate the same graph structure for C test
			cNodes := make([]*vertexNode, len(tt.edges))
			for i, edge := range tt.edges {
				cNodes[i] = addVertexNode(cGraph, &edge.from, &edge.to)
			}

			cNodeToRemove := cNodes[tt.removeIndex]
			cResult := removeVertexNodeC(cGraph, cNodeToRemove)

			// Compare results
			if goResult != cResult {
				t.Errorf("Result mismatch: Go=%d, C=%d", goResult, cResult)
			}

			// Both should succeed (return 0)
			if goResult != 0 {
				t.Errorf("Expected removal to succeed (return 0), got %d", goResult)
			}

			// Verify graph size after removal
			expectedFinalSize := expectedInitialSize - 1
			if goGraph.Size != expectedFinalSize {
				t.Errorf("Final graph size: got %d, want %d", goGraph.Size, expectedFinalSize)
			}

			// Verify the removed node is no longer findable in the graph
			removedNodeFound := false
			for i := 0; i < int(goGraph.NumBuckets); i++ {
				node := goGraph.Buckets[i]
				for node != nil {
					if node == nodeToRemove {
						removedNodeFound = true
						break
					}
					node = node.Next
				}
				if removedNodeFound {
					break
				}
			}

			if removedNodeFound {
				t.Error("Removed node was still found in the graph")
			}

			// Verify remaining nodes are still present
			for i, originalNode := range nodes {
				if i == tt.removeIndex {
					continue // Skip the removed node
				}

				nodeFound := false
				for j := 0; j < int(goGraph.NumBuckets); j++ {
					node := goGraph.Buckets[j]
					for node != nil {
						if geoAlmostEqual(&node.From, &originalNode.From) &&
							geoAlmostEqual(&node.To, &originalNode.To) {
							nodeFound = true
							break
						}
						node = node.Next
					}
					if nodeFound {
						break
					}
				}

				if !nodeFound {
					t.Errorf("Expected node %d to remain in graph after removal", i)
				}
			}
		})
	}
}

func Test_removeVertexNode_not_found_parity(t *testing.T) {
	// Test removing a node that doesn't exist in the graph
	goGraph := &vertexGraph{
		Buckets:    make([]*vertexNode, 10),
		NumBuckets: 10,
		Size:       0,
		Res:        8,
	}

	// Add one node
	fromVtx := LatLng{Lat: 37.775, Lng: -122.418}
	toVtx := LatLng{Lat: 37.776, Lng: -122.419}
	_ = addVertexNode(goGraph, &fromVtx, &toVtx)

	// Create a different node that's not in the graph
	differentFromVtx := LatLng{Lat: 40.0, Lng: -75.0}
	differentToVtx := LatLng{Lat: 40.1, Lng: -75.1}
	nonExistentNode := &vertexNode{
		From: differentFromVtx,
		To:   differentToVtx,
		Next: nil,
	}

	// Test Go implementation - should fail to find the node
	goResult := removeVertexNode(goGraph, nonExistentNode)

	// Test C implementation - should also fail
	cGraph := &vertexGraph{
		Buckets:    make([]*vertexNode, 10),
		NumBuckets: 10,
		Size:       1, // Simulate having one node
		Res:        8,
	}
	cResult := removeVertexNodeC(cGraph, nonExistentNode)

	// Both should fail (return 1)
	if goResult != 1 {
		t.Errorf("Go implementation: expected failure (return 1), got %d", goResult)
	}

	if cResult != 1 {
		t.Errorf("C implementation: expected failure (return 1), got %d", cResult)
	}

	// Graph size should remain unchanged
	if goGraph.Size != 1 {
		t.Errorf("Graph size should remain 1, got %d", goGraph.Size)
	}

	// The original node should still be there
	index := _hashVertex(&fromVtx, goGraph.Res, goGraph.NumBuckets)
	if goGraph.Buckets[index] == nil {
		t.Error("Original node should still exist in the graph")
	}
}

func Test_removeVertexNode_empty_graph_parity(t *testing.T) {
	// Test removing from an empty graph
	goGraph := &vertexGraph{
		Buckets:    make([]*vertexNode, 5),
		NumBuckets: 5,
		Size:       0,
		Res:        8,
	}

	// Create a node to "remove" (even though it's not in the graph)
	fromVtx := LatLng{Lat: 37.775, Lng: -122.418}
	toVtx := LatLng{Lat: 37.776, Lng: -122.419}
	nonExistentNode := &vertexNode{
		From: fromVtx,
		To:   toVtx,
		Next: nil,
	}

	// Test Go implementation
	goResult := removeVertexNode(goGraph, nonExistentNode)

	// Test C implementation
	cGraph := &vertexGraph{
		Buckets:    make([]*vertexNode, 5),
		NumBuckets: 5,
		Size:       0,
		Res:        8,
	}
	cResult := removeVertexNodeC(cGraph, nonExistentNode)

	// Both should fail (return 1)
	if goResult != 1 {
		t.Errorf("Go implementation: expected failure (return 1), got %d", goResult)
	}

	if cResult != 1 {
		t.Errorf("C implementation: expected failure (return 1), got %d", cResult)
	}

	// Graph size should remain 0
	if goGraph.Size != 0 {
		t.Errorf("Graph size should remain 0, got %d", goGraph.Size)
	}
}
