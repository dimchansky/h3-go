// Tests ported from testH3IndexInternal.c
package h3

import "testing"

func TestFaceIjkToH3ExtremeCoordinates(t *testing.T) {
	t.Parallel()

	// Test cases for resolution 0 bounds checking
	tests := []struct {
		name     string
		fijk     FaceIJK
		res      int32
		expected H3Index
	}{
		{
			name:     "i out of bounds at res 0",
			fijk:     FaceIJK{Face: 0, Coord: CoordIJK{I: 3, J: 0, K: 0}},
			res:      0,
			expected: H3_NULL,
		},
		{
			name:     "j out of bounds at res 0",
			fijk:     FaceIJK{Face: 1, Coord: CoordIJK{I: 0, J: 4, K: 0}},
			res:      0,
			expected: H3_NULL,
		},
		{
			name:     "k out of bounds at res 0",
			fijk:     FaceIJK{Face: 2, Coord: CoordIJK{I: 2, J: 0, K: 5}},
			res:      0,
			expected: H3_NULL,
		},
		{
			name:     "i out of bounds at res 1",
			fijk:     FaceIJK{Face: 3, Coord: CoordIJK{I: 6, J: 0, K: 0}},
			res:      1,
			expected: H3_NULL,
		},
		{
			name:     "j out of bounds at res 1",
			fijk:     FaceIJK{Face: 4, Coord: CoordIJK{I: 0, J: 7, K: 1}},
			res:      1,
			expected: H3_NULL,
		},
		{
			name:     "k out of bounds at res 1",
			fijk:     FaceIJK{Face: 5, Coord: CoordIJK{I: 2, J: 0, K: 8}},
			res:      1,
			expected: H3_NULL,
		},
		{
			name:     "i out of bounds at res 2",
			fijk:     FaceIJK{Face: 6, Coord: CoordIJK{I: 18, J: 0, K: 0}},
			res:      2,
			expected: H3_NULL,
		},
		{
			name:     "j out of bounds at res 2",
			fijk:     FaceIJK{Face: 7, Coord: CoordIJK{I: 0, J: 19, K: 1}},
			res:      2,
			expected: H3_NULL,
		},
		{
			name:     "k out of bounds at res 2",
			fijk:     FaceIJK{Face: 8, Coord: CoordIJK{I: 2, J: 0, K: 20}},
			res:      2,
			expected: H3_NULL,
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
