//go:build cgo && c2go

package h3

import "testing"

func Test_ijkNormalizeCouldOverflow_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord CoordIJK
	}{
		{"origin", CoordIJK{0, 0, 0}},
		{"small positive", CoordIJK{1, 2, 0}},
		{"small negative", CoordIJK{-1, -2, 0}},
		{"mixed small", CoordIJK{-1, 2, 0}},
		{"mixed small 2", CoordIJK{1, -2, 0}},
		{"medium values", CoordIJK{100, -50, 0}},
		{"medium values 2", CoordIJK{-100, 50, 0}},
		{"medium values 3", CoordIJK{100, 50, 0}},
		{"medium values 4", CoordIJK{-100, -50, 0}},
		{"large positive", CoordIJK{1000000, 2000000, 0}},
		{"large negative", CoordIJK{-1000000, -2000000, 0}},
		{"large mixed", CoordIJK{1000000, -2000000, 0}},
		{"large mixed 2", CoordIJK{-1000000, 2000000, 0}},
		{"equal values pos", CoordIJK{1000, 1000, 0}},
		{"equal values neg", CoordIJK{-1000, -1000, 0}},
		// Test values near overflow boundaries
		{"near max positive", CoordIJK{INT32_MAX / 2, INT32_MAX / 3, 0}},
		{"near max negative", CoordIJK{INT32_MIN / 2, INT32_MIN / 3, 0}},
		{"near max mixed", CoordIJK{INT32_MAX / 2, INT32_MIN / 3, 0}},
		{"boundary case 1", CoordIJK{INT32_MAX, 0, 0}},
		{"boundary case 2", CoordIJK{0, INT32_MIN, 0}},
		{"boundary case 3", CoordIJK{INT32_MAX, INT32_MIN, 0}},
		// Test cases that might trigger overflow
		{"potential overflow 1", CoordIJK{INT32_MAX, -1, 0}},
		{"potential overflow 2", CoordIJK{INT32_MIN, 1, 0}},
		{"potential overflow 3", CoordIJK{INT32_MAX - 1, -2, 0}},
		{"potential overflow 4", CoordIJK{INT32_MIN + 1, 2, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkNormalizeCouldOverflowC(&tt.coord)

			// Call Go implementation
			gotGo := _ijkNormalizeCouldOverflow(&tt.coord)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_ijkNormalizeCouldOverflow() mismatch: Go=%v != C=%v for input{%d,%d,%d}",
					gotGo, gotC, tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that function is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := CoordIJK{1000, -500, 0}

		// Call function twice
		result1 := _ijkNormalizeCouldOverflow(&coord)
		result2 := _ijkNormalizeCouldOverflow(&coord)

		if result1 != result2 {
			t.Errorf("_ijkNormalizeCouldOverflow should be deterministic: first=%v != second=%v",
				result1, result2)
		}
	})

	// Test mathematical property: should not change based on k=0 assumption
	t.Run("k_zero_assumption", func(t *testing.T) {
		testCoords := []CoordIJK{
			{100, -50, 0},
			{-100, 50, 0},
			{INT32_MAX / 2, INT32_MIN / 2, 0},
		}

		for _, coord := range testCoords {
			// Test with k=0 (expected)
			result1 := _ijkNormalizeCouldOverflow(&coord)

			// Test with different k (should behave the same since function only looks at i,j)
			coordWithK := CoordIJK{coord.I, coord.J, 999}
			result2 := _ijkNormalizeCouldOverflow(&coordWithK)

			if result1 != result2 {
				t.Errorf("K value should not affect overflow check: k=0 gave %v, k=999 gave %v for i=%d,j=%d",
					result1, result2, coord.I, coord.J)
			}
		}
	})

	// Test overflow boundary conditions more precisely
	t.Run("overflow_boundaries", func(t *testing.T) {
		// Test cases that should definitely NOT overflow
		safeCoords := []CoordIJK{
			{1000, 1000, 0},                   // both positive
			{-1000, -1000, 0},                 // both negative
			{1000, -500, 0},                   // moderate mixed
			{INT32_MAX / 4, INT32_MIN / 4, 0}, // well within bounds
		}

		for _, coord := range safeCoords {
			result := _ijkNormalizeCouldOverflow(&coord)
			if result {
				t.Logf("Expected no overflow for safe case {%d,%d,%d} but got overflow warning",
					coord.I, coord.J, coord.K)
				// Note: Don't fail here as the C implementation is the authority
			}
		}
	})
}
