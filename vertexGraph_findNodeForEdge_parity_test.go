//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_findNodeForEdge_parity(t *testing.T) {
	tests := []struct {
		name        string
		numBuckets  int32
		res         int32
		edges       []struct{ from, to LatLng }
		searchFrom  LatLng
		searchTo    *LatLng // nil for wildcard search
		expectFound bool
		description string
	}{
		{
			name:       "find existing edge exact match",
			numBuckets: 10,
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
				{LatLng{Lat: 40.7128, Lng: -74.0060}, LatLng{Lat: 40.7129, Lng: -74.0061}},
			},
			searchFrom:  LatLng{Lat: 37.775, Lng: -122.418},
			searchTo:    &LatLng{Lat: 37.776, Lng: -122.419},
			expectFound: true,
			description: "Find exact edge match in graph with multiple edges",
		},
		{
			name:       "find nonexistent edge",
			numBuckets: 10,
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
			},
			searchFrom:  LatLng{Lat: 99.0, Lng: 99.0},
			searchTo:    &LatLng{Lat: 98.0, Lng: 98.0},
			expectFound: false,
			description: "Search for edge that doesn't exist",
		},
		{
			name:       "wildcard search with nil toVtx",
			numBuckets: 5,
			res:        8,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 51.5074, Lng: -0.1278}, LatLng{Lat: 51.5080, Lng: -0.1285}},
				{LatLng{Lat: 51.5074, Lng: -0.1278}, LatLng{Lat: 51.5090, Lng: -0.1290}},
			},
			searchFrom:  LatLng{Lat: 51.5074, Lng: -0.1278},
			searchTo:    nil, // Wildcard search
			expectFound: true,
			description: "Find any edge starting from specific vertex (wildcard)",
		},
		{
			name:       "wildcard search no match",
			numBuckets: 5,
			res:        8,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 51.5074, Lng: -0.1278}, LatLng{Lat: 51.5080, Lng: -0.1285}},
			},
			searchFrom:  LatLng{Lat: 99.0, Lng: 99.0},
			searchTo:    nil, // Wildcard search
			expectFound: false,
			description: "Wildcard search for vertex that doesn't exist",
		},
		{
			name:       "hash collision search",
			numBuckets: 1, // Force hash collisions
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
				{LatLng{Lat: 37.780, Lng: -122.420}, LatLng{Lat: 37.781, Lng: -122.421}},
				{LatLng{Lat: 37.785, Lng: -122.425}, LatLng{Lat: 37.786, Lng: -122.426}},
			},
			searchFrom:  LatLng{Lat: 37.780, Lng: -122.420},
			searchTo:    &LatLng{Lat: 37.781, Lng: -122.421},
			expectFound: true,
			description: "Find edge in bucket with hash collisions (linked list search)",
		},
		{
			name:        "empty graph search",
			numBuckets:  10,
			res:         9,
			edges:       []struct{ from, to LatLng }{}, // Empty
			searchFrom:  LatLng{Lat: 37.775, Lng: -122.418},
			searchTo:    &LatLng{Lat: 37.776, Lng: -122.419},
			expectFound: false,
			description: "Search in empty graph",
		},
		{
			name:       "zero coordinates",
			numBuckets: 15,
			res:        10,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 0.0, Lng: 0.0}, LatLng{Lat: 0.001, Lng: 0.001}},
			},
			searchFrom:  LatLng{Lat: 0.0, Lng: 0.0},
			searchTo:    &LatLng{Lat: 0.001, Lng: 0.001},
			expectFound: true,
			description: "Find edge with zero start coordinates",
		},
		{
			name:       "high resolution precision",
			numBuckets: 20,
			res:        15,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 40.712800001, Lng: -74.006000001}, LatLng{Lat: 40.712800002, Lng: -74.006000002}},
			},
			searchFrom:  LatLng{Lat: 40.712800001, Lng: -74.006000001},
			searchTo:    &LatLng{Lat: 40.712800002, Lng: -74.006000002},
			expectFound: true,
			description: "Find edge at high resolution with small coordinate differences",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize Go graph
			goGraph := &vertexGraph{
				Buckets:    make([]*vertexNode, tt.numBuckets),
				NumBuckets: tt.numBuckets,
				Size:       0,
				Res:        tt.res,
			}

			// Add all test edges to the graph
			for _, edge := range tt.edges {
				addVertexNode(goGraph, &edge.from, &edge.to)
			}

			// Test Go implementation
			goResult := findNodeForEdge(goGraph, &tt.searchFrom, tt.searchTo)
			goFound := goResult != nil

			// Verify the search result
			if goFound != tt.expectFound {
				t.Errorf("findNodeForEdge: expected found=%t, got found=%t", tt.expectFound, goFound)
			}

			if tt.expectFound && goResult != nil {
				// Verify the found node has the correct from vertex
				if !geoAlmostEqual(&goResult.From, &tt.searchFrom) {
					t.Errorf("Found node from vertex mismatch: got (%f, %f), want (%f, %f)",
						goResult.From.Lat, goResult.From.Lng, tt.searchFrom.Lat, tt.searchFrom.Lng)
				}

				// Verify to vertex if not a wildcard search
				if tt.searchTo != nil {
					if !geoAlmostEqual(&goResult.To, tt.searchTo) {
						t.Errorf("Found node to vertex mismatch: got (%f, %f), want (%f, %f)",
							goResult.To.Lat, goResult.To.Lng, tt.searchTo.Lat, tt.searchTo.Lng)
					}
				}
			}

			// Test C implementation for parity
			cResult := findNodeForEdgeC(goGraph, &tt.searchFrom, tt.searchTo)
			cFound := cResult != nil

			if cFound != goFound {
				t.Errorf("Go/C parity mismatch: Go found=%t, C found=%t", goFound, cFound)
			}

			if tt.expectFound && goResult != nil && cResult != nil {
				// Compare Go and C results for found vertices
				if !geoAlmostEqual(&goResult.From, &cResult.From) {
					t.Errorf("Go/C from vertex parity mismatch: Go=(%f, %f), C=(%f, %f)",
						goResult.From.Lat, goResult.From.Lng, cResult.From.Lat, cResult.From.Lng)
				}

				if !geoAlmostEqual(&goResult.To, &cResult.To) {
					t.Errorf("Go/C to vertex parity mismatch: Go=(%f, %f), C=(%f, %f)",
						goResult.To.Lat, goResult.To.Lng, cResult.To.Lat, cResult.To.Lng)
				}
			}
		})
	}
}

func Test_findNodeForEdge_linked_list_traversal_parity(t *testing.T) {
	// Test specifically for linked list traversal when multiple nodes are in same bucket
	numBuckets := int32(1) // Force all items into same bucket
	res := int32(9)

	goGraph := &vertexGraph{
		Buckets:    make([]*vertexNode, numBuckets),
		NumBuckets: numBuckets,
		Size:       0,
		Res:        res,
	}

	// Add multiple edges that will all hash to the same bucket
	edges := []struct {
		from LatLng
		to   LatLng
	}{
		{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
		{LatLng{Lat: 37.780, Lng: -122.420}, LatLng{Lat: 37.781, Lng: -122.421}},
		{LatLng{Lat: 37.785, Lng: -122.425}, LatLng{Lat: 37.786, Lng: -122.426}},
		{LatLng{Lat: 37.790, Lng: -122.430}, LatLng{Lat: 37.791, Lng: -122.431}},
	}

	// Add all edges
	for _, edge := range edges {
		addVertexNode(goGraph, &edge.from, &edge.to)
	}

	// Test finding each edge - should find all of them despite hash collisions
	for i, edge := range edges {
		goResult := findNodeForEdge(goGraph, &edge.from, &edge.to)
		if goResult == nil {
			t.Errorf("Failed to find edge %d: from (%f, %f) to (%f, %f)",
				i, edge.from.Lat, edge.from.Lng, edge.to.Lat, edge.to.Lng)
			continue
		}

		// Verify found node matches expected edge
		if !geoAlmostEqual(&goResult.From, &edge.from) {
			t.Errorf("Edge %d from vertex mismatch: got (%f, %f), want (%f, %f)",
				i, goResult.From.Lat, goResult.From.Lng, edge.from.Lat, edge.from.Lng)
		}

		if !geoAlmostEqual(&goResult.To, &edge.to) {
			t.Errorf("Edge %d to vertex mismatch: got (%f, %f), want (%f, %f)",
				i, goResult.To.Lat, goResult.To.Lng, edge.to.Lat, edge.to.Lng)
		}

		// Test wildcard search (should find the same or compatible edge)
		wildcardResult := findNodeForEdge(goGraph, &edge.from, nil)
		if wildcardResult == nil {
			t.Errorf("Failed to find edge %d with wildcard search from (%f, %f)",
				i, edge.from.Lat, edge.from.Lng)
		} else if !geoAlmostEqual(&wildcardResult.From, &edge.from) {
			t.Errorf("Edge %d wildcard from vertex mismatch: got (%f, %f), want (%f, %f)",
				i, wildcardResult.From.Lat, wildcardResult.From.Lng, edge.from.Lat, edge.from.Lng)
		}
	}

	// Test searching for non-existent edge in populated bucket
	nonExistentFrom := LatLng{Lat: 99.0, Lng: 99.0}
	nonExistentTo := LatLng{Lat: 98.0, Lng: 98.0}

	result := findNodeForEdge(goGraph, &nonExistentFrom, &nonExistentTo)
	if result != nil {
		t.Errorf("Expected not to find non-existent edge, but found: from (%f, %f) to (%f, %f)",
			result.From.Lat, result.From.Lng, result.To.Lat, result.To.Lng)
	}

	// Test wildcard search for non-existent from vertex
	wildcardResult := findNodeForEdge(goGraph, &nonExistentFrom, nil)
	if wildcardResult != nil {
		t.Errorf("Expected not to find non-existent vertex in wildcard search, but found: from (%f, %f)",
			wildcardResult.From.Lat, wildcardResult.From.Lng)
	}
}
