//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test__initVertexNode_parity(t *testing.T) {
	tests := []struct {
		name    string
		fromVtx LatLng
		toVtx   LatLng
	}{
		{
			name:    "zero coordinates",
			fromVtx: LatLng{Lat: 0.0, Lng: 0.0},
			toVtx:   LatLng{Lat: 0.0, Lng: 0.0},
		},
		{
			name:    "positive coordinates",
			fromVtx: LatLng{Lat: 0.5, Lng: 0.7},
			toVtx:   LatLng{Lat: 0.8, Lng: 0.9},
		},
		{
			name:    "negative coordinates",
			fromVtx: LatLng{Lat: -0.5, Lng: -0.7},
			toVtx:   LatLng{Lat: -0.8, Lng: -0.9},
		},
		{
			name:    "mixed coordinates",
			fromVtx: LatLng{Lat: -0.5, Lng: 0.7},
			toVtx:   LatLng{Lat: 0.8, Lng: -0.9},
		},
		{
			name:    "large values",
			fromVtx: LatLng{Lat: 1.5707963267948966, Lng: 3.141592653589793}, // PI/2 and PI
			toVtx:   LatLng{Lat: -1.5707963267948966, Lng: -3.141592653589793},
		},
		{
			name:    "very small values",
			fromVtx: LatLng{Lat: 1e-10, Lng: 1e-10},
			toVtx:   LatLng{Lat: -1e-10, Lng: -1e-10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and initialize Go node
			goNode := &VertexNode{}
			_initVertexNode(goNode, &tt.fromVtx, &tt.toVtx)

			// Create and initialize C node
			cNode := &VertexNode{}
			_initVertexNodeC(cNode, &tt.fromVtx, &tt.toVtx)

			// Compare from vertex
			if goNode.From.Lat != cNode.From.Lat {
				t.Errorf("from.Lat mismatch: Go = %v, C = %v", goNode.From.Lat, cNode.From.Lat)
			}
			if goNode.From.Lng != cNode.From.Lng {
				t.Errorf("from.Lng mismatch: Go = %v, C = %v", goNode.From.Lng, cNode.From.Lng)
			}

			// Compare to vertex
			if goNode.To.Lat != cNode.To.Lat {
				t.Errorf("to.Lat mismatch: Go = %v, C = %v", goNode.To.Lat, cNode.To.Lat)
			}
			if goNode.To.Lng != cNode.To.Lng {
				t.Errorf("to.Lng mismatch: Go = %v, C = %v", goNode.To.Lng, cNode.To.Lng)
			}

			// Compare next pointer (should be nil)
			if goNode.Next != nil {
				t.Errorf("Go node.Next is not nil")
			}
			if cNode.Next != nil {
				t.Errorf("C node.Next is not nil")
			}
		})
	}
}
