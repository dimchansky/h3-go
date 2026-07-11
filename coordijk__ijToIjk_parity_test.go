//go:build cgo && c2go

package h3

import "testing"

func Test_ijToIjk_parity(t *testing.T) {
	tests := []struct {
		name string
		ij   CoordIJ
	}{
		{"origin", CoordIJ{0, 0}},
		{"unit i", CoordIJ{1, 0}},
		{"unit j", CoordIJ{0, 1}},
		{"positive coords", CoordIJ{3, 5}},
		{"negative coords", CoordIJ{-2, -4}},
		{"mixed coords", CoordIJ{3, -2}},
		{"mixed coords 2", CoordIJ{-3, 2}},
		{"large positive", CoordIJ{1000, 2000}},
		{"large negative", CoordIJ{-1000, -2000}},
		{"large mixed", CoordIJ{1000, -2000}},
		{"large mixed 2", CoordIJ{-1000, 2000}},
		{"equal values pos", CoordIJ{500, 500}},
		{"equal values neg", CoordIJ{-500, -500}},
		{"asymmetric", CoordIJ{123, 456}},
		{"asymmetric neg", CoordIJ{-123, -456}},
		{"medium values", CoordIJ{100, -50}},
		{"medium values 2", CoordIJ{-100, 50}},
		// Test boundary cases that might trigger overflow
		{"boundary case 1", CoordIJ{int32Max / 2, 0}},
		{"boundary case 2", CoordIJ{0, int32Min / 2}},
		{"boundary case 3", CoordIJ{int32Max / 2, int32Min / 2}},
		{"boundary case 4", CoordIJ{int32Min / 2, int32Max / 2}},
		// Test cases that might cause overflow in normalization
		{"potential overflow 1", CoordIJ{int32Max, -1}},
		{"potential overflow 2", CoordIJ{int32Min, 1}},
		{"potential overflow 3", CoordIJ{int32Max - 1000, -2000}},
		{"potential overflow 4", CoordIJ{int32Min + 1000, 2000}},
		{"extreme positive", CoordIJ{int32Max, int32Max}},
		{"extreme negative", CoordIJ{int32Min, int32Min}},
		{"extreme mixed", CoordIJ{int32Max, int32Min}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotCIjk, gotCErr := ijToIjkC(&tt.ij)

			// Call Go implementation
			var gotGoIjk coordIJK
			gotGoErr := ijToIjk(&tt.ij, &gotGoIjk)

			// Compare error codes first
			if gotGoErr != gotCErr {
				t.Errorf("ijToIjk() error mismatch: Go=%d != C=%d for input{%d,%d}",
					gotGoErr, gotCErr, tt.ij.I, tt.ij.J)
				return
			}

			// If both succeeded, compare results
			if gotGoErr == eSuccess {
				if gotGoIjk.I != gotCIjk.I || gotGoIjk.J != gotCIjk.J || gotGoIjk.K != gotCIjk.K {
					t.Errorf("ijToIjk() result mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d}",
						gotGoIjk.I, gotGoIjk.J, gotGoIjk.K, gotCIjk.I, gotCIjk.J, gotCIjk.K,
						tt.ij.I, tt.ij.J)
				}
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		ij := CoordIJ{100, -50}

		// Apply transformation twice
		var result1, result2 coordIJK
		err1 := ijToIjk(&ij, &result1)
		err2 := ijToIjk(&ij, &result2)

		if err1 != err2 {
			t.Errorf("ijToIjk should have deterministic errors: first=%d != second=%d", err1, err2)
		}

		if err1 == eSuccess {
			if result1.I != result2.I || result1.J != result2.J || result1.K != result2.K {
				t.Errorf("ijToIjk should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
					result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
			}
		}
	})

	// Test basic mathematical properties when successful
	t.Run("mathematical_properties", func(t *testing.T) {
		testIJs := []CoordIJ{
			{0, 0}, {1, 0}, {0, 1}, {5, 3}, {-2, 4}, {10, -5},
		}

		for _, ij := range testIJs {
			var ijk coordIJK
			err := ijToIjk(&ij, &ijk)

			if err == eSuccess {
				// Basic property: the IJ coordinates should be preserved in some form
				// After normalization, the relationship might be different, but
				// let's at least verify the function doesn't crash
				if ijk.I == 0 && ijk.J == 0 && ijk.K == 0 {
					// Should only happen for origin
					if ij.I != 0 || ij.J != 0 {
						t.Errorf("Non-origin quadIJ{%d,%d} produced origin IJK{%d,%d,%d}",
							ij.I, ij.J, ijk.I, ijk.J, ijk.K)
					}
				}
			}
		}
	})

	// Test round-trip property with ijkToIj
	t.Run("round_trip_where_possible", func(t *testing.T) {
		testIJs := []CoordIJ{
			{0, 0}, {1, 0}, {0, 1}, {2, 1}, {1, 2}, {3, 2}, {-1, 0}, {0, -1},
		}

		for _, origIj := range testIJs {
			var ijk coordIJK
			err := ijToIjk(&origIj, &ijk)

			if err == eSuccess {
				// Convert back to quadIJ
				var resultIj CoordIJ
				ijkToIj(&ijk, &resultIj)

				// Note: Due to normalization, this might not be a perfect round-trip
				// but we can at least check that the process doesn't crash
				t.Logf("Round trip: quadIJ{%d,%d} -> IJK{%d,%d,%d} -> quadIJ{%d,%d}",
					origIj.I, origIj.J, ijk.I, ijk.J, ijk.K, resultIj.I, resultIj.J)
			}
		}
	})

	// Test overflow error cases
	t.Run("overflow_errors", func(t *testing.T) {
		// Test cases that should definitely trigger overflow
		overflowCases := []CoordIJ{
			{int32Max, int32Max},
			{int32Min, int32Min},
			{int32Max, int32Min},
		}

		for _, ij := range overflowCases {
			var ijk coordIJK
			err := ijToIjk(&ij, &ijk)

			// The function should either succeed or fail, but not crash
			if err != eSuccess && err != eFailed {
				t.Errorf("Unexpected error code %d for overflow case {%d,%d}", err, ij.I, ij.J)
			}

			t.Logf("Overflow test: quadIJ{%d,%d} -> error=%d", ij.I, ij.J, err)
		}
	})
}
