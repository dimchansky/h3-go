//go:build cgo

package c2go

import (
	"testing"
)

func Test_firstVertexNode_parity(t *testing.T) {
	testCases := []struct {
		name        string
		graph       VertexGraph
		expectNil   bool
		description string
	}{
		{
			name: "empty graph",
			graph: VertexGraph{
				Buckets:    nil,
				NumBuckets: 0,
				Size:       0,
				Res:        9,
			},
			expectNil:   true,
			description: "Empty graph should return nil",
		},
		{
			name: "graph with empty buckets",
			graph: VertexGraph{
				Buckets:    make([]*VertexNode, 3),
				NumBuckets: 3,
				Size:       0,
				Res:        9,
			},
			expectNil:   true,
			description: "Graph with all empty buckets should return nil",
		},
		{
			name: "graph with node in first bucket",
			graph: VertexGraph{
				Buckets: []*VertexNode{
					{
						From: LatLng{Lat: 0.0, Lng: 0.0},
						To:   LatLng{Lat: 0.1, Lng: 0.1},
						Next: nil,
					},
					nil,
					nil,
				},
				NumBuckets: 3,
				Size:       1,
				Res:        9,
			},
			expectNil:   false,
			description: "Graph with node in first bucket should return that node",
		},
		{
			name: "graph with node in second bucket",
			graph: VertexGraph{
				Buckets: []*VertexNode{
					nil,
					{
						From: LatLng{Lat: 0.2, Lng: 0.2},
						To:   LatLng{Lat: 0.3, Lng: 0.3},
						Next: nil,
					},
					nil,
				},
				NumBuckets: 3,
				Size:       1,
				Res:        9,
			},
			expectNil:   false,
			description: "Graph with node in second bucket should return that node",
		},
		{
			name: "graph with multiple nodes, first in second bucket",
			graph: VertexGraph{
				Buckets: []*VertexNode{
					nil,
					{
						From: LatLng{Lat: 0.4, Lng: 0.4},
						To:   LatLng{Lat: 0.5, Lng: 0.5},
						Next: &VertexNode{
							From: LatLng{Lat: 0.6, Lng: 0.6},
							To:   LatLng{Lat: 0.7, Lng: 0.7},
							Next: nil,
						},
					},
					{
						From: LatLng{Lat: 0.8, Lng: 0.8},
						To:   LatLng{Lat: 0.9, Lng: 0.9},
						Next: nil,
					},
				},
				NumBuckets: 3,
				Size:       3,
				Res:        9,
			},
			expectNil:   false,
			description: "Graph with multiple nodes should return first node found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goResult := firstVertexNode(&tc.graph)

			// Test C implementation
			cResult := firstVertexNodeC(&tc.graph)

			// Verify both return nil or both return non-nil
			if tc.expectNil {
				if goResult != nil {
					t.Errorf("Go implementation expected nil, got non-nil")
				}
				if cResult != nil {
					t.Errorf("C implementation expected nil, got non-nil")
				}
			} else {
				if goResult == nil {
					t.Errorf("Go implementation expected non-nil, got nil")
				}
				if cResult == nil {
					t.Errorf("C implementation expected non-nil, got nil")
				}

				// If both are non-nil, verify they match
				if goResult != nil && cResult != nil {
					// For this function, we expect the first node found to have the same coordinates
					// Check the From coordinates match
					if !latLngsAlmostEqual(&goResult.From, &cResult.From) {
						t.Errorf("From coordinates mismatch: Go=%+v, C=%+v", goResult.From, cResult.From)
					}
					// Check the To coordinates match
					if !latLngsAlmostEqual(&goResult.To, &cResult.To) {
						t.Errorf("To coordinates mismatch: Go=%+v, C=%+v", goResult.To, cResult.To)
					}
				}
			}
		})
	}
}

// latLngsAlmostEqual checks if two LatLng values are approximately equal
func latLngsAlmostEqual(a, b *LatLng) bool {
	const tolerance = 1e-12
	return abs(a.Lat.Rad()-b.Lat.Rad()) < tolerance && abs(a.Lng.Rad()-b.Lng.Rad()) < tolerance
}
