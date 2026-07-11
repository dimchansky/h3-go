// Tests ported from H3 v4.4.0: src/apps/testapps/testH3IndexInternal.c.
package h3

import "testing"

func TestFaceIjkToH3ExtremeCoordinates(t *testing.T) {
	t.Parallel()

	// Test cases for resolution 0 bounds checking
	tests := []struct {
		name     string
		fijk     faceIJK
		res      int32
		expected h3Index
	}{
		{
			name:     "i out of bounds at res 0",
			fijk:     faceIJK{Face: 0, Coord: coordIJK{I: 3, J: 0, K: 0}},
			res:      0,
			expected: h3Null,
		},
		{
			name:     "j out of bounds at res 0",
			fijk:     faceIJK{Face: 1, Coord: coordIJK{I: 0, J: 4, K: 0}},
			res:      0,
			expected: h3Null,
		},
		{
			name:     "k out of bounds at res 0",
			fijk:     faceIJK{Face: 2, Coord: coordIJK{I: 2, J: 0, K: 5}},
			res:      0,
			expected: h3Null,
		},
		{
			name:     "i out of bounds at res 1",
			fijk:     faceIJK{Face: 3, Coord: coordIJK{I: 6, J: 0, K: 0}},
			res:      1,
			expected: h3Null,
		},
		{
			name:     "j out of bounds at res 1",
			fijk:     faceIJK{Face: 4, Coord: coordIJK{I: 0, J: 7, K: 1}},
			res:      1,
			expected: h3Null,
		},
		{
			name:     "k out of bounds at res 1",
			fijk:     faceIJK{Face: 5, Coord: coordIJK{I: 2, J: 0, K: 8}},
			res:      1,
			expected: h3Null,
		},
		{
			name:     "i out of bounds at res 2",
			fijk:     faceIJK{Face: 6, Coord: coordIJK{I: 18, J: 0, K: 0}},
			res:      2,
			expected: h3Null,
		},
		{
			name:     "j out of bounds at res 2",
			fijk:     faceIJK{Face: 7, Coord: coordIJK{I: 0, J: 19, K: 1}},
			res:      2,
			expected: h3Null,
		},
		{
			name:     "k out of bounds at res 2",
			fijk:     faceIJK{Face: 8, Coord: coordIJK{I: 2, J: 0, K: 20}},
			res:      2,
			expected: h3Null,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := _faceIjkToH3(&tt.fijk, tt.res)
			if result != tt.expected {
				t.Errorf("_faceIjkToH3() = %v, want %v", result, tt.expected)
			}
		})
	}
}
