//go:build cgo && c2go && !h3v450

package h3

import (
	"testing"
)

func Test_findNodeForVertex_parity(t *testing.T) {
	tests := []struct {
		name        string
		numBuckets  int32
		res         int32
		edges       []struct{ from, to LatLng }
		searchFrom  LatLng
		expectFound bool
		description string
	}{
		{
			name:       "find vertex with single edge",
			numBuckets: 10,
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
			},
			searchFrom:  LatLng{Lat: 37.775, Lng: -122.418},
			expectFound: true,
			description: "Find vertex that has exactly one outgoing edge",
		},
		{
			name:       "find vertex with multiple edges",
			numBuckets: 10,
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.777, Lng: -122.420}},
				{LatLng{Lat: 40.7128, Lng: -74.0060}, LatLng{Lat: 40.7129, Lng: -74.0061}},
			},
			searchFrom:  LatLng{Lat: 37.775, Lng: -122.418},
			expectFound: true,
			description: "Find vertex that has multiple outgoing edges",
		},
		{
			name:       "find nonexistent vertex",
			numBuckets: 10,
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
			},
			searchFrom:  LatLng{Lat: 99.0, Lng: 99.0},
			expectFound: false,
			description: "Search for vertex that doesn't exist as start of any edge",
		},
		{
			name:       "find vertex in hash collision bucket",
			numBuckets: 1, // Force hash collisions
			res:        9,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 37.775, Lng: -122.418}, LatLng{Lat: 37.776, Lng: -122.419}},
				{LatLng{Lat: 37.780, Lng: -122.420}, LatLng{Lat: 37.781, Lng: -122.421}},
				{LatLng{Lat: 37.785, Lng: -122.425}, LatLng{Lat: 37.786, Lng: -122.426}},
			},
			searchFrom:  LatLng{Lat: 37.780, Lng: -122.420},
			expectFound: true,
			description: "Find vertex in bucket with hash collisions (linked list search)",
		},
		{
			name:        "search in empty graph",
			numBuckets:  10,
			res:         9,
			edges:       []struct{ from, to LatLng }{}, // Empty
			searchFrom:  LatLng{Lat: 37.775, Lng: -122.418},
			expectFound: false,
			description: "Search for vertex in empty graph",
		},
		{
			name:       "find vertex with zero coordinates",
			numBuckets: 15,
			res:        10,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 0.0, Lng: 0.0}, LatLng{Lat: 0.001, Lng: 0.001}},
			},
			searchFrom:  LatLng{Lat: 0.0, Lng: 0.0},
			expectFound: true,
			description: "Find vertex with zero start coordinates",
		},
		{
			name:       "find vertex at high resolution",
			numBuckets: 20,
			res:        15,
			edges: []struct{ from, to LatLng }{
				{LatLng{Lat: 40.712800001, Lng: -74.006000001}, LatLng{Lat: 40.712800002, Lng: -74.006000002}},
			},
			searchFrom:  LatLng{Lat: 40.712800001, Lng: -74.006000001},
			expectFound: true,
			description: "Find vertex at high resolution with small coordinate differences",
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
			goResult := findNodeForVertex(goGraph, &tt.searchFrom)
			goFound := goResult != nil

			// Verify the search result
			if goFound != tt.expectFound {
				t.Errorf("findNodeForVertex: expected found=%t, got found=%t", tt.expectFound, goFound)
			}

			if tt.expectFound && goResult != nil {
				// Verify the found node has the correct from vertex
				if !geoAlmostEqual(&goResult.From, &tt.searchFrom) {
					t.Errorf("Found node from vertex mismatch: got (%f, %f), want (%f, %f)",
						goResult.From.Lat, goResult.From.Lng, tt.searchFrom.Lat, tt.searchFrom.Lng)
				}
			}

			// Test C implementation for parity
			cResult := findNodeForVertexC(goGraph, &tt.searchFrom)
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

func Test_findNodeForVertex_multiple_edges_same_vertex_parity(t *testing.T) {
	// Test specifically when a vertex has multiple outgoing edges
	numBuckets := int32(10)
	res := int32(9)

	goGraph := &vertexGraph{
		Buckets:    make([]*vertexNode, numBuckets),
		NumBuckets: numBuckets,
		Size:       0,
		Res:        res,
	}

	// Add multiple edges from the same starting vertex
	startVertex := LatLng{Lat: 37.775, Lng: -122.418}
	edges := []LatLng{
		{Lat: 37.776, Lng: -122.419},
		{Lat: 37.777, Lng: -122.420},
		{Lat: 37.778, Lng: -122.421},
	}

	// Add all edges from the same start vertex
	for _, toVertex := range edges {
		addVertexNode(goGraph, &startVertex, &toVertex)
	}

	// findNodeForVertex should find one of the edges (typically the first one found)
	goResult := findNodeForVertex(goGraph, &startVertex)
	if goResult == nil {
		t.Errorf("Expected to find vertex with multiple edges")
		return
	}

	// Verify the found node has the correct from vertex
	if !geoAlmostEqual(&goResult.From, &startVertex) {
		t.Errorf("Found node from vertex mismatch: got (%f, %f), want (%f, %f)",
			goResult.From.Lat, goResult.From.Lng, startVertex.Lat, startVertex.Lng)
	}

	// Verify that the found edge's 'to' vertex is one of our expected edges
	found := false
	for _, expectedTo := range edges {
		if geoAlmostEqual(&goResult.To, &expectedTo) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Found node's 'to' vertex (%f, %f) is not one of the expected edges",
			goResult.To.Lat, goResult.To.Lng)
	}

	// Test C implementation for parity
	cResult := findNodeForVertexC(goGraph, &startVertex)
	if cResult == nil {
		t.Errorf("C implementation failed to find vertex with multiple edges")
		return
	}

	// Verify Go and C implementations find vertices with the same from coordinates
	if !geoAlmostEqual(&goResult.From, &cResult.From) {
		t.Errorf("Go/C from vertex parity mismatch: Go=(%f, %f), C=(%f, %f)",
			goResult.From.Lat, goResult.From.Lng, cResult.From.Lat, cResult.From.Lng)
	}

	// Both should find valid edges (though potentially different ones)
	goToValid := false
	cToValid := false
	for _, expectedTo := range edges {
		if geoAlmostEqual(&goResult.To, &expectedTo) {
			goToValid = true
		}
		if geoAlmostEqual(&cResult.To, &expectedTo) {
			cToValid = true
		}
	}

	if !goToValid {
		t.Errorf("Go result's 'to' vertex (%f, %f) is not one of the expected edges",
			goResult.To.Lat, goResult.To.Lng)
	}
	if !cToValid {
		t.Errorf("C result's 'to' vertex (%f, %f) is not one of the expected edges",
			cResult.To.Lat, cResult.To.Lng)
	}
}
