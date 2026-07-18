//go:build cgo && c2go && !h3v450

package h3

import (
	"testing"
)

// Test_vertexGraphToLinkedGeo_parity tests behavioral equivalence for _vertexGraphToLinkedGeo.
// This function creates a linkedGeoPolygon from a vertex graph.
// Due to the complexity of converting entire graph structures between Go and C,
// this test focuses on verifying the behavior through integration tests.
func Test_vertexGraphToLinkedGeo_parity(t *testing.T) {
	// Test 1: Empty graph produces empty polygon
	t.Run("empty_graph", func(t *testing.T) {
		graph := &vertexGraph{
			Buckets:    nil,
			NumBuckets: 0,
			Size:       0,
			Res:        0,
		}

		var out linkedGeoPolygon
		_vertexGraphToLinkedGeo(graph, &out)

		// Verify the output is empty
		if out.First != nil {
			t.Errorf("Expected empty polygon for empty graph, got First=%v", out.First)
		}
		if out.Last != nil {
			t.Errorf("Expected empty polygon for empty graph, got Last=%v", out.Last)
		}
		if out.Next != nil {
			t.Errorf("Expected empty polygon for empty graph, got Next=%v", out.Next)
		}
	})

	// Test 2: Single edge creates a simple loop
	t.Run("single_edge", func(t *testing.T) {
		// Create a simple triangle graph
		graph := &vertexGraph{
			Buckets:    make([]*vertexNode, 3),
			NumBuckets: 3,
			Size:       3,
			Res:        9,
		}

		// Create a triangle of vertices
		v1 := LatLng{Lat: 0.1, Lng: 0.1}
		v2 := LatLng{Lat: 0.1, Lng: 0.2}
		v3 := LatLng{Lat: 0.2, Lng: 0.1}

		// Build a simple closed triangle
		node1 := &vertexNode{From: v1, To: v2, Next: nil}
		node2 := &vertexNode{From: v2, To: v3, Next: nil}
		node3 := &vertexNode{From: v3, To: v1, Next: nil}

		// Place nodes in buckets based on hash
		idx1 := _hashVertex(&v1, graph.Res, graph.NumBuckets)
		idx2 := _hashVertex(&v2, graph.Res, graph.NumBuckets)
		idx3 := _hashVertex(&v3, graph.Res, graph.NumBuckets)

		// Add nodes to their respective buckets
		if graph.Buckets[idx1] == nil {
			graph.Buckets[idx1] = node1
		} else {
			node1.Next = graph.Buckets[idx1]
			graph.Buckets[idx1] = node1
		}

		if graph.Buckets[idx2] == nil {
			graph.Buckets[idx2] = node2
		} else {
			node2.Next = graph.Buckets[idx2]
			graph.Buckets[idx2] = node2
		}

		if graph.Buckets[idx3] == nil {
			graph.Buckets[idx3] = node3
		} else {
			node3.Next = graph.Buckets[idx3]
			graph.Buckets[idx3] = node3
		}

		var out linkedGeoPolygon
		_vertexGraphToLinkedGeo(graph, &out)

		// Verify we got a polygon with one loop
		if out.First == nil {
			t.Fatal("Expected non-empty polygon, got nil First")
		}
		if out.First != out.Last {
			t.Errorf("Expected single loop, First=%p Last=%p", out.First, out.Last)
		}

		// Count coordinates in the loop
		coordCount := 0
		coord := out.First.First
		for coord != nil {
			coordCount++
			coord = coord.Next
		}

		// We should have 3 coordinates forming the triangle
		if coordCount != 3 {
			t.Errorf("Expected 3 coordinates in loop, got %d", coordCount)
		}

		// Verify the graph is now empty (all edges consumed)
		if graph.Size != 0 {
			t.Errorf("Expected graph size to be 0 after conversion, got %d", graph.Size)
		}
	})

	// Test 3: Multiple disconnected edges create multiple loops
	t.Run("multiple_loops", func(t *testing.T) {
		// Create a graph with two disconnected triangles
		graph := &vertexGraph{
			Buckets:    make([]*vertexNode, 10),
			NumBuckets: 10,
			Size:       6,
			Res:        9,
		}

		// First triangle
		v1a := LatLng{Lat: 0.1, Lng: 0.1}
		v2a := LatLng{Lat: 0.1, Lng: 0.2}
		v3a := LatLng{Lat: 0.2, Lng: 0.1}

		// Second triangle (disconnected)
		v1b := LatLng{Lat: 0.5, Lng: 0.5}
		v2b := LatLng{Lat: 0.5, Lng: 0.6}
		v3b := LatLng{Lat: 0.6, Lng: 0.5}

		// Build first triangle
		node1a := &vertexNode{From: v1a, To: v2a, Next: nil}
		node2a := &vertexNode{From: v2a, To: v3a, Next: nil}
		node3a := &vertexNode{From: v3a, To: v1a, Next: nil}

		// Build second triangle
		node1b := &vertexNode{From: v1b, To: v2b, Next: nil}
		node2b := &vertexNode{From: v2b, To: v3b, Next: nil}
		node3b := &vertexNode{From: v3b, To: v1b, Next: nil}

		// Add all nodes to buckets
		nodes := []*vertexNode{node1a, node2a, node3a, node1b, node2b, node3b}
		for _, node := range nodes {
			idx := _hashVertex(&node.From, graph.Res, graph.NumBuckets)
			if graph.Buckets[idx] == nil {
				graph.Buckets[idx] = node
			} else {
				node.Next = graph.Buckets[idx]
				graph.Buckets[idx] = node
			}
		}

		var out linkedGeoPolygon
		_vertexGraphToLinkedGeo(graph, &out)

		// Verify we got a polygon with loops
		if out.First == nil {
			t.Fatal("Expected non-empty polygon, got nil First")
		}

		// Count loops
		loopCount := 0
		loop := out.First
		for loop != nil {
			loopCount++
			loop = loop.Next
		}

		// We should have 2 loops (one for each triangle)
		if loopCount != 2 {
			t.Errorf("Expected 2 loops, got %d", loopCount)
		}

		// Verify the graph is now empty (all edges consumed)
		if graph.Size != 0 {
			t.Errorf("Expected graph size to be 0 after conversion, got %d", graph.Size)
		}
	})
}
