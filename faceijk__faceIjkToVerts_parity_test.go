//go:build cgo && c2go

package h3

import "testing"

func Test_faceijk_faceIjkToVerts_ParityWithC(t *testing.T) {
	testCases := []struct {
		fijk faceIJK
		res  int32
		desc string
	}{
		// Class II resolutions
		{faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}}, 0, "origin face 0 res 0"},
		{faceIJK{Face: 1, Coord: coordIJK{I: 1, J: 2, K: 3}}, 2, "face 1 res 2"},
		{faceIJK{Face: 5, Coord: coordIJK{I: 10, J: 15, K: 20}}, 4, "face 5 res 4"},
		{faceIJK{Face: 9, Coord: coordIJK{I: 100, J: 200, K: 150}}, 6, "face 9 res 6"},
		{faceIJK{Face: 15, Coord: coordIJK{I: 500, J: 750, K: 250}}, 8, "face 15 res 8"},

		// Class III resolutions
		{faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}}, 1, "origin face 0 res 1"},
		{faceIJK{Face: 2, Coord: coordIJK{I: 2, J: 4, K: 6}}, 3, "face 2 res 3"},
		{faceIJK{Face: 7, Coord: coordIJK{I: 25, J: 50, K: 75}}, 5, "face 7 res 5"},
		{faceIJK{Face: 12, Coord: coordIJK{I: 200, J: 400, K: 300}}, 7, "face 12 res 7"},
		{faceIJK{Face: 19, Coord: coordIJK{I: 1000, J: 1500, K: 500}}, 9, "face 19 res 9"},

		// Edge cases
		{faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 0, K: 0}}, 0, "i-axis face 0"},
		{faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 1, K: 0}}, 0, "j-axis face 0"},
		{faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 1}}, 0, "k-axis face 0"},

		// Higher resolutions
		{faceIJK{Face: 10, Coord: coordIJK{I: 5000, J: 7500, K: 2500}}, 10, "high res even"},
		{faceIJK{Face: 11, Coord: coordIJK{I: 10000, J: 15000, K: 5000}}, 11, "high res odd"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Make copies for Go implementation
			goFijk := tc.fijk
			goRes := tc.res
			goVerts := make([]faceIJK, numHexVerts)

			// Make copies for C implementation
			cFijk := tc.fijk
			cRes := tc.res
			cVerts := make([]faceIJK, numHexVerts)

			// Test Go implementation
			_faceIjkToVerts(&goFijk, &goRes, goVerts)

			// Test C implementation
			_faceIjkToVertsC(&cFijk, &cRes, cVerts)

			// Compare resolution (may be modified for Class III)
			if goRes != cRes {
				t.Fatalf("resolution mismatch: go=%d c=%d", goRes, cRes)
			}

			// Compare all vertices
			for v := 0; v < numHexVerts; v++ {
				if goVerts[v].Face != cVerts[v].Face {
					t.Fatalf("vertex %d face mismatch: go=%d c=%d", v, goVerts[v].Face, cVerts[v].Face)
				}
				if goVerts[v].Coord.I != cVerts[v].Coord.I {
					t.Fatalf("vertex %d I mismatch: go=%d c=%d", v, goVerts[v].Coord.I, cVerts[v].Coord.I)
				}
				if goVerts[v].Coord.J != cVerts[v].Coord.J {
					t.Fatalf("vertex %d J mismatch: go=%d c=%d", v, goVerts[v].Coord.J, cVerts[v].Coord.J)
				}
				if goVerts[v].Coord.K != cVerts[v].Coord.K {
					t.Fatalf("vertex %d K mismatch: go=%d c=%d", v, goVerts[v].Coord.K, cVerts[v].Coord.K)
				}
			}
		})
	}
}

func Test_faceijk_faceIjkToVerts_ResolutionModification_ParityWithC(t *testing.T) {
	// Test that Class III resolutions get modified correctly
	classIIICases := []struct {
		originalRes int32
		expectedRes int32
		desc        string
	}{
		{1, 2, "res 1 -> 2"},
		{3, 4, "res 3 -> 4"},
		{5, 6, "res 5 -> 6"},
		{7, 8, "res 7 -> 8"},
		{9, 10, "res 9 -> 10"},
		{11, 12, "res 11 -> 12"},
		{13, 14, "res 13 -> 14"},
		{15, 16, "res 15 -> 16"}, // Edge case - max resolution
	}

	baseFijk := faceIJK{Face: 0, Coord: coordIJK{I: 100, J: 200, K: 150}}

	for _, tc := range classIIICases {
		t.Run(tc.desc, func(t *testing.T) {
			// Test Go implementation
			goFijk := baseFijk
			goRes := tc.originalRes
			goVerts := make([]faceIJK, numHexVerts)
			_faceIjkToVerts(&goFijk, &goRes, goVerts)

			// Test C implementation
			cFijk := baseFijk
			cRes := tc.originalRes
			cVerts := make([]faceIJK, numHexVerts)
			_faceIjkToVertsC(&cFijk, &cRes, cVerts)

			// Both should modify resolution the same way
			if goRes != cRes {
				t.Fatalf("resolution modification mismatch: go=%d c=%d", goRes, cRes)
			}

			if goRes != tc.expectedRes {
				t.Fatalf("resolution not modified as expected: got=%d expected=%d", goRes, tc.expectedRes)
			}
		})
	}
}

func Test_faceijk_faceIjkToVerts_ClassII_NoResolutionChange_ParityWithC(t *testing.T) {
	// Test that Class II resolutions don't get modified
	classIICases := []int32{0, 2, 4, 6, 8, 10, 12, 14}
	baseFijk := faceIJK{Face: 5, Coord: coordIJK{I: 50, J: 100, K: 75}}

	for _, originalRes := range classIICases {
		// Test Go implementation
		goFijk := baseFijk
		goRes := originalRes
		goVerts := make([]faceIJK, numHexVerts)
		_faceIjkToVerts(&goFijk, &goRes, goVerts)

		// Test C implementation
		cFijk := baseFijk
		cRes := originalRes
		cVerts := make([]faceIJK, numHexVerts)
		_faceIjkToVertsC(&cFijk, &cRes, cVerts)

		// Resolution should not change for Class II
		if goRes != originalRes {
			t.Fatalf("Class II resolution incorrectly modified: original=%d new=%d", originalRes, goRes)
		}
		if cRes != originalRes {
			t.Fatalf("C Class II resolution incorrectly modified: original=%d new=%d", originalRes, cRes)
		}
		if goRes != cRes {
			t.Fatalf("resolution mismatch: go=%d c=%d", goRes, cRes)
		}
	}
}
