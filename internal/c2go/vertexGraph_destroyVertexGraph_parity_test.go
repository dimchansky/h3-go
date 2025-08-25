//go:build cgo

package c2go

import (
	"testing"
)

func Test_destroyVertexGraph_parity(t *testing.T) {
	tests := []struct {
		name       string
		numBuckets int32
		res        int32
		addNodes   []struct {
			from LatLng
			to   LatLng
		}
	}{
		{
			name:       "empty graph",
			numBuckets: 0,
			res:        5,
			addNodes:   nil,
		},
		{
			name:       "graph with buckets but no nodes",
			numBuckets: 10,
			res:        5,
			addNodes:   nil,
		},
		{
			name:       "graph with single node",
			numBuckets: 10,
			res:        5,
			addNodes: []struct {
				from LatLng
				to   LatLng
			}{
				{
					from: LatLng{Lat: 0.5, Lng: 0.5},
					to:   LatLng{Lat: 0.6, Lng: 0.6},
				},
			},
		},
		{
			name:       "graph with multiple nodes",
			numBuckets: 10,
			res:        5,
			addNodes: []struct {
				from LatLng
				to   LatLng
			}{
				{
					from: LatLng{Lat: 0.1, Lng: 0.1},
					to:   LatLng{Lat: 0.2, Lng: 0.2},
				},
				{
					from: LatLng{Lat: 0.3, Lng: 0.3},
					to:   LatLng{Lat: 0.4, Lng: 0.4},
				},
				{
					from: LatLng{Lat: 0.5, Lng: 0.5},
					to:   LatLng{Lat: 0.6, Lng: 0.6},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Go graph
			goGraph := &VertexGraph{}
			initVertexGraph(goGraph, tt.numBuckets, tt.res)

			// Add nodes to Go graph
			for _, node := range tt.addNodes {
				addVertexNode(goGraph, &node.from, &node.to)
			}

			// Create C graph
			cGraph := &VertexGraph{}
			initVertexGraphC(cGraph, tt.numBuckets, tt.res)

			// Note: For this test, we're mainly verifying that destroyVertexGraph
			// properly clears the graph. The C wrapper simulates the behavior.

			// Destroy Go graph
			destroyVertexGraph(goGraph)

			// Destroy C graph
			destroyVertexGraphC(cGraph)

			// Verify both graphs are in the same state after destruction
			if goGraph.Buckets != nil {
				t.Errorf("Go graph buckets not nil after destroy")
			}
			if cGraph.Buckets != nil {
				t.Errorf("C graph buckets not nil after destroy")
			}

			// After destroy, size should be 0
			if goGraph.Size != 0 {
				t.Errorf("Go graph size = %d, want 0", goGraph.Size)
			}
			if cGraph.Size != 0 {
				t.Errorf("C graph size = %d, want 0", cGraph.Size)
			}
		})
	}
}
