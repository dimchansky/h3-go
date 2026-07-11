//go:build cgo && c2go

package h3

import "testing"

func Test_validatePolygonFlags_parity(t *testing.T) {
	tests := []struct {
		name  string
		flags uint32
	}{
		{"valid_center", uint32(ContainmentCenter)},
		{"valid_full", uint32(ContainmentFull)},
		{"valid_overlapping", uint32(ContainmentOverlapping)},
		{"valid_overlapping_bbox", uint32(ContainmentOverlappingBBox)},
		{"invalid_containment", uint32(ContainmentInvalid)},
		{"invalid_flag_bit_4", 16},                           // Bit 4 set (outside mask)
		{"invalid_flag_bit_5", 32},                           // Bit 5 set
		{"invalid_flag_bit_31", 1 << 31},                     // High bit set
		{"invalid_combined", uint32(ContainmentCenter) | 16}, // Valid containment + invalid flag
		{"invalid_high_containment", 8},                      // Containment value above INVALID
		{"all_mask_bits", flagContainmentModeMask},           // All containment bits set
		{"zero_flags", 0},                                    // No flags set (ContainmentCenter)
		{"max_uint32", ^uint32(0)},                           // All bits set
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
			0,  // ContainmentCenter (valid)
			1,  // ContainmentFull (valid)
			2,  // ContainmentOverlapping (valid)
			3,  // ContainmentOverlappingBBox (valid)
			4,  // ContainmentInvalid (invalid)
			5,  // Above ContainmentInvalid (invalid)
			15, // All containment bits set (invalid)
			16, // First bit outside mask (invalid)
			17, // ContainmentFull + invalid bit (invalid)
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
			{"boundary_valid_max", uint32(ContainmentOverlappingBBox), true},
			{"boundary_invalid_min", uint32(ContainmentInvalid), false},
			{"mask_boundary", flagContainmentModeMask, false},
			{"just_above_mask", flagContainmentModeMask + 1, false},
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
					expectedResult = uint32(eOptionInvalid)
				}

				if goResult != expectedResult {
					t.Errorf("Edge case %s: expected %d, got %d for flags=0x%x",
						tc.name, expectedResult, goResult, tc.flags)
				}
			})
		}
	})
}
