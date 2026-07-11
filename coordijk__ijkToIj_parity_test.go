//go:build cgo && c2go

package h3

import "testing"

func Test_ijkToIj_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord coordIJK
	}{
		{"origin", coordIJK{0, 0, 0}},
		{"unit i", coordIJK{1, 0, 0}},
		{"unit j", coordIJK{0, 1, 0}},
		{"unit k", coordIJK{0, 0, 1}},
		{"positive coords", coordIJK{1, 2, 3}},
		{"mixed coords", coordIJK{3, -2, 5}},
		{"large coords", coordIJK{15, 10, 8}},
		{"negative coords", coordIJK{-1, -2, -3}},
		{"asymmetric", coordIJK{5, 3, 2}},
		{"small values", coordIJK{2, 1, 1}},
		{"normalized", coordIJK{2, 1, 0}},
		{"equal components", coordIJK{3, 3, 3}},
		{"zero i", coordIJK{0, 5, 2}},
		{"zero j", coordIJK{5, 0, 2}},
		{"zero k", coordIJK{5, 3, 0}},
		{"large k", coordIJK{2, 3, 10}},
		{"negative k", coordIJK{2, 3, -5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := ijkToIjC(&tt.coord)

			// Call Go implementation
			var gotGo CoordIJ
			ijkToIj(&tt.coord, &gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J {
				t.Errorf("ijkToIj() mismatch: Go{%d,%d} != C{%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotC.I, gotC.J,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test the mathematical relationship: quadIJ = IJK - K
	t.Run("mathematical_property", func(t *testing.T) {
		testCoords := []coordIJK{
			{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1},
			{5, 3, 2}, {-2, 4, -1}, {10, -5, 3},
		}

		for _, coord := range testCoords {
			var result CoordIJ
			ijkToIj(&coord, &result)

			expectedI := coord.I - coord.K
			expectedJ := coord.J - coord.K

			if result.I != expectedI || result.J != expectedJ {
				t.Errorf("Mathematical relationship failed for {%d,%d,%d}: got{%d,%d} expected{%d,%d}",
					coord.I, coord.J, coord.K, result.I, result.J, expectedI, expectedJ)
			}
		}
	})

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := coordIJK{7, 3, 2}

		// Apply transformation twice
		var result1, result2 CoordIJ
		ijkToIj(&coord, &result1)
		ijkToIj(&coord, &result2)

		if result1.I != result2.I || result1.J != result2.J {
			t.Errorf("ijkToIj should be deterministic: first{%d,%d} != second{%d,%d}",
				result1.I, result1.J, result2.I, result2.J)
		}
	})

	// Test edge cases
	t.Run("edge_cases", func(t *testing.T) {
		edgeCases := []struct {
			name     string
			coord    coordIJK
			expected CoordIJ
		}{
			{"zero_k_no_change", coordIJK{5, 3, 0}, CoordIJ{5, 3}},
			{"equal_k_zeros_result", coordIJK{2, 2, 2}, CoordIJ{0, 0}},
			{"large_k_negative_result", coordIJK{1, 1, 5}, CoordIJ{-4, -4}},
		}

		for _, tc := range edgeCases {
			t.Run(tc.name, func(t *testing.T) {
				var result CoordIJ
				ijkToIj(&tc.coord, &result)

				if result.I != tc.expected.I || result.J != tc.expected.J {
					t.Errorf("Edge case %s: got{%d,%d} expected{%d,%d}",
						tc.name, result.I, result.J, tc.expected.I, tc.expected.J)
				}
			})
		}
	})
}
