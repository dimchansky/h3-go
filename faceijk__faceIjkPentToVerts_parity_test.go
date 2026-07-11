//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_faceIjkPentToVerts_parity(t *testing.T) {
	tests := []struct {
		name string
		fijk faceIJK
		res  int32
	}{
		{
			name: "pentagon face 0 res 0",
			fijk: faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}},
			res:  0,
		},
		{
			name: "pentagon face 1 res 1 class III",
			fijk: faceIJK{Face: 1, Coord: coordIJK{I: 1, J: 0, K: 0}},
			res:  1,
		},
		{
			name: "pentagon face 2 res 2 class II",
			fijk: faceIJK{Face: 2, Coord: coordIJK{I: 0, J: 1, K: 0}},
			res:  2,
		},
		{
			name: "pentagon face 5 res 3 class III",
			fijk: faceIJK{Face: 5, Coord: coordIJK{I: 2, J: 1, K: 0}},
			res:  3,
		},
		{
			name: "pentagon face 10 res 4 class II",
			fijk: faceIJK{Face: 10, Coord: coordIJK{I: 1, J: 1, K: 1}},
			res:  4,
		},
		{
			name: "pentagon face 15 res 5 class III",
			fijk: faceIJK{Face: 15, Coord: coordIJK{I: 3, J: 2, K: 1}},
			res:  5,
		},
		{
			name: "pentagon face 19 res 6 class II",
			fijk: faceIJK{Face: 19, Coord: coordIJK{I: 4, J: 3, K: 2}},
			res:  6,
		},
		{
			name: "pentagon negative coords",
			fijk: faceIJK{Face: 8, Coord: coordIJK{I: -1, J: 2, K: 0}},
			res:  4,
		},
		{
			name: "pentagon large coords",
			fijk: faceIJK{Face: 12, Coord: coordIJK{I: 10, J: 15, K: 5}},
			res:  7,
		},
		{
			name: "pentagon zero face high res",
			fijk: faceIJK{Face: 0, Coord: coordIJK{I: 20, J: 30, K: 10}},
			res:  8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Go implementation
			goFijk := tt.fijk
			goRes := tt.res
			goVerts := make([]faceIJK, numPentVerts)
			_faceIjkPentToVerts(&goFijk, &goRes, goVerts)

			// Test C implementation
			cFijk := tt.fijk
			cRes := tt.res
			cVerts := make([]faceIJK, numPentVerts)
			_faceIjkPentToVertsC(&cFijk, &cRes, cVerts)

			// Compare resolution modification
			if goRes != cRes {
				t.Errorf("resolution mismatch: Go=%d, C=%d", goRes, cRes)
			}

			// Compare all vertices
			for i := 0; i < numPentVerts; i++ {
				if goVerts[i].Face != cVerts[i].Face {
					t.Errorf("vertex %d face mismatch: Go=%d, C=%d", i, goVerts[i].Face, cVerts[i].Face)
				}
				if goVerts[i].Coord.I != cVerts[i].Coord.I {
					t.Errorf("vertex %d I mismatch: Go=%d, C=%d", i, goVerts[i].Coord.I, cVerts[i].Coord.I)
				}
				if goVerts[i].Coord.J != cVerts[i].Coord.J {
					t.Errorf("vertex %d J mismatch: Go=%d, C=%d", i, goVerts[i].Coord.J, cVerts[i].Coord.J)
				}
				if goVerts[i].Coord.K != cVerts[i].Coord.K {
					t.Errorf("vertex %d K mismatch: Go=%d, C=%d", i, goVerts[i].Coord.K, cVerts[i].Coord.K)
				}
			}
		})
	}
}

func Test_faceIjkPentToVerts_resolutionModification(t *testing.T) {
	// Test that Class III resolutions get incremented by 1
	classIIITests := []struct {
		name    string
		fijk    faceIJK
		initRes int32
		wantRes int32
	}{
		{
			name:    "res 1 class III -> res 2",
			fijk:    faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}},
			initRes: 1,
			wantRes: 2,
		},
		{
			name:    "res 3 class III -> res 4",
			fijk:    faceIJK{Face: 5, Coord: coordIJK{I: 1, J: 1, K: 0}},
			initRes: 3,
			wantRes: 4,
		},
		{
			name:    "res 5 class III -> res 6",
			fijk:    faceIJK{Face: 10, Coord: coordIJK{I: 2, J: 1, K: 1}},
			initRes: 5,
			wantRes: 6,
		},
	}

	for _, tt := range classIIITests {
		t.Run(tt.name, func(t *testing.T) {
			fijk := tt.fijk
			res := tt.initRes
			verts := make([]faceIJK, numPentVerts)

			_faceIjkPentToVerts(&fijk, &res, verts)

			if res != tt.wantRes {
				t.Errorf("expected resolution %d, got %d", tt.wantRes, res)
			}
		})
	}

	// Test that Class II resolutions remain unchanged
	classIITests := []struct {
		name    string
		fijk    faceIJK
		initRes int32
	}{
		{
			name:    "res 0 class II unchanged",
			fijk:    faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}},
			initRes: 0,
		},
		{
			name:    "res 2 class II unchanged",
			fijk:    faceIJK{Face: 8, Coord: coordIJK{I: 1, J: 0, K: 0}},
			initRes: 2,
		},
		{
			name:    "res 4 class II unchanged",
			fijk:    faceIJK{Face: 15, Coord: coordIJK{I: 2, J: 1, K: 0}},
			initRes: 4,
		},
	}

	for _, tt := range classIITests {
		t.Run(tt.name, func(t *testing.T) {
			fijk := tt.fijk
			res := tt.initRes
			originalRes := tt.initRes
			verts := make([]faceIJK, numPentVerts)

			_faceIjkPentToVerts(&fijk, &res, verts)

			if res != originalRes {
				t.Errorf("expected resolution unchanged at %d, got %d", originalRes, res)
			}
		})
	}
}
