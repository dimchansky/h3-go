//go:build cgo && c2go

package h3

import "testing"

func Test_validatePolygonFlags_parity(t *testing.T) {
	tests := []struct {
		name  string
		flags uint32
	}{
		{"valid_center", uint32(CONTAINMENT_CENTER)},
		{"valid_full", uint32(CONTAINMENT_FULL)},
		{"valid_overlapping", uint32(CONTAINMENT_OVERLAPPING)},
		{"valid_overlapping_bbox", uint32(CONTAINMENT_OVERLAPPING_BBOX)},
		{"invalid_containment", uint32(CONTAINMENT_INVALID)},
		{"invalid_flag_bit_4", 16},                            // Bit 4 set (outside mask)
		{"invalid_flag_bit_5", 32},                            // Bit 5 set
		{"invalid_flag_bit_31", 1 << 31},                      // High bit set
		{"invalid_combined", uint32(CONTAINMENT_CENTER) | 16}, // Valid containment + invalid flag
		{"invalid_high_containment", 8},                       // Containment value above INVALID
		{"all_mask_bits", FLAG_CONTAINMENT_MODE_MASK},         // All containment bits set
		{"zero_flags", 0},                                     // No flags set (CONTAINMENT_CENTER)
		{"max_uint32", ^uint32(0)},                            // All bits set
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			cResult := validatePolygonFlagsC(tt.flags)

			// Call Go implementation
			goResult := uint32(validatePolygonFlags(tt.flags))

			if goResult != cResult {
				t.Errorf("validatePolygonFlags() parity mismatch: Go=%d, C=%d for flags=0x%x",
					goResult, cResult, tt.flags)
			}
		})
	}

	// Test specific flag combinations that are important
	t.Run("comprehensive_flag_tests", func(t *testing.T) {
		flagCombinations := []uint32{
			0,  // CONTAINMENT_CENTER (valid)
			1,  // CONTAINMENT_FULL (valid)
			2,  // CONTAINMENT_OVERLAPPING (valid)
			3,  // CONTAINMENT_OVERLAPPING_BBOX (valid)
			4,  // CONTAINMENT_INVALID (invalid)
			5,  // Above CONTAINMENT_INVALID (invalid)
			15, // All containment bits set (invalid)
			16, // First bit outside mask (invalid)
			17, // CONTAINMENT_FULL + invalid bit (invalid)
			31, // High bits within range (invalid)
		}

		for _, flags := range flagCombinations {
			cResult := validatePolygonFlagsC(flags)
			goResult := uint32(validatePolygonFlags(flags))

			if goResult != cResult {
				t.Errorf("Flag combination 0x%x: Go=%d, C=%d", flags, goResult, cResult)
			}
		}
	})

	// Test edge cases
	t.Run("edge_cases", func(t *testing.T) {
		edgeCases := []struct {
			name        string
			flags       uint32
			expectValid bool
		}{
			{"boundary_valid_max", uint32(CONTAINMENT_OVERLAPPING_BBOX), true},
			{"boundary_invalid_min", uint32(CONTAINMENT_INVALID), false},
			{"mask_boundary", FLAG_CONTAINMENT_MODE_MASK, false},
			{"just_above_mask", FLAG_CONTAINMENT_MODE_MASK + 1, false},
		}

		for _, tc := range edgeCases {
			t.Run(tc.name, func(t *testing.T) {
				cResult := validatePolygonFlagsC(tc.flags)
				goResult := uint32(validatePolygonFlags(tc.flags))

				if goResult != cResult {
					t.Errorf("Edge case %s: Go=%d, C=%d for flags=0x%x",
						tc.name, goResult, cResult, tc.flags)
				}

				expectedResult := uint32(0)
				if !tc.expectValid {
					expectedResult = uint32(E_OPTION_INVALID)
				}

				if goResult != expectedResult {
					t.Errorf("Edge case %s: expected %d, got %d for flags=0x%x",
						tc.name, expectedResult, goResult, tc.flags)
				}
			})
		}
	})
}
