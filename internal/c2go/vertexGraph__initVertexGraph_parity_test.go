//go:build cgo

package c2go

import (
	"testing"
)

func Test_initVertexGraph_parity(t *testing.T) {
	testCases := []struct {
		name       string
		numBuckets int
		res        int
	}{
		{
			name:       "zero buckets",
			numBuckets: 0,
			res:        5,
		},
		{
			name:       "small bucket count",
			numBuckets: 5,
			res:        7,
		},
		{
			name:       "medium bucket count",
			numBuckets: 100,
			res:        9,
		},
		{
			name:       "large bucket count",
			numBuckets: 1000,
			res:        12,
		},
		{
			name:       "max resolution",
			numBuckets: 50,
			res:        15,
		},
		{
			name:       "min resolution",
			numBuckets: 10,
			res:        0,
		},
		{
			name:       "single bucket",
			numBuckets: 1,
			res:        6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goGraph VertexGraph
			initVertexGraph(&goGraph, tc.numBuckets, tc.res)

			// Test C implementation
			var cGraph VertexGraph
			initVertexGraphC(&cGraph, tc.numBuckets, tc.res)

			// Compare results
			if goGraph.NumBuckets != cGraph.NumBuckets {
				t.Errorf("NumBuckets mismatch: Go=%d, C=%d", goGraph.NumBuckets, cGraph.NumBuckets)
			}
			if goGraph.Size != cGraph.Size {
				t.Errorf("Size mismatch: Go=%d, C=%d", goGraph.Size, cGraph.Size)
			}
			if goGraph.Res != cGraph.Res {
				t.Errorf("Res mismatch: Go=%d, C=%d", goGraph.Res, cGraph.Res)
			}

			// Check buckets allocation behavior
			if tc.numBuckets > 0 {
				if goGraph.Buckets == nil {
					t.Error("Go: Buckets should not be nil for numBuckets > 0")
				}
				if len(goGraph.Buckets) != tc.numBuckets {
					t.Errorf("Go: Buckets length mismatch: expected %d, got %d", tc.numBuckets, len(goGraph.Buckets))
				}
				if cGraph.Buckets == nil {
					t.Error("C: Buckets should not be nil for numBuckets > 0")
				}
			} else {
				if goGraph.Buckets != nil {
					t.Error("Go: Buckets should be nil for numBuckets == 0")
				}
				if cGraph.Buckets != nil {
					t.Error("C: Buckets should be nil for numBuckets == 0")
				}
			}
		})
	}
}
