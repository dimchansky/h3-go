//go:build cgo

package c2go

import (
	"fmt"
	"testing"
)

func TestIsValidCellParity(t *testing.T) {
	testCases := []struct {
		name string
		h3   H3Index
	}{
		// Valid cases
		{"res_0_base_0", 0x8001fffffffffff},
		{"res_0_base_7_pentagon", 0x8007fffffffffff},
		{"res_1_regular", 0x81283ffffffffff},
		{"res_2_regular", 0x8228bffffffffff},
		{"res_15_max_res", 0x8f283082803ffff},

		// Invalid cases - bad top bits
		{"bad_high_bit", 0x0001fffffffffff},  // High bit should be 0
		{"bad_mode_bits", 0x8401fffffffffff}, // Mode should be 0001
		{"bad_reserved", 0x8081fffffffffff},  // Reserved should be 000

		// Invalid cases - invalid base cell (manually constructed)
		{"invalid_base_cell_122", 0x8f40fffffffffff}, // Invalid base cell
		{"invalid_base_cell_127", 0x8fe0fffffffffff}, // Invalid base cell

		// Invalid cases - digit 7 in wrong place
		{"digit_7_at_res_1", 0x81783ffffffffff}, // 7 at resolution 1
		{"digit_7_at_res_5", 0x852837fffffffff}, // 7 at resolution 5

		// Invalid cases - not all 7s after resolution
		{"missing_trailing_7s", 0x82283bfffffffff}, // Should have all 7s after res 2

		// Invalid cases - deleted subsequences (simplified)
		{"potential_deleted_subseq_1", 0x810e5ffffffffff}, // Base cell with leading 5
		{"potential_deleted_subseq_2", 0x81185ffffffffff}, // Base cell with leading 5

		// Edge cases
		{"zero_index", 0x0000000000000000},
		{"max_uint64", 0xffffffffffffffff},
		{"valid_pentagon_res_3", 0x8370bffffffffff},
		{"valid_high_res", 0x8a283082803ffff},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goResult := isValidCell(tc.h3)

			// Test C implementation
			cResult := isValidCellC(tc.h3)

			// Compare results (convert bool to int for comparison)
			if goResult != cResult {
				t.Errorf("Result mismatch for 0x%x: Go=%t (%v), C=%v",
					uint64(tc.h3), goResult, goResult, cResult)
			}

			// Log results for visibility
			t.Logf("H3Index 0x%x -> valid=%t", uint64(tc.h3), goResult)
		})
	}
}

func TestIsValidCellSystematic(t *testing.T) {
	// Test a range of systematic cases
	testCases := []H3Index{
		// All base resolution cells
		0x8001fffffffffff, // base 1
		0x8003fffffffffff, // base 3
		0x8005fffffffffff, // base 5
		0x8007fffffffffff, // base 7 (pentagon)
		0x800dfffffffffff, // base 13 (pentagon)
		0x800ffffffffffff, // base 15 (pentagon)

		// Various resolutions
		0x81283ffffffffff, // res 1
		0x8228bffffffffff, // res 2
		0x832830fffffffff, // res 3
		0x8427cffffffffff, // res 4
		0x85283473fffffff, // res 5

		// High resolution cases
		0x8a283082803ffff, // res 10
		0x8f283082800bfff, // res 15 (max)

		// Known valid H3 indexes from real locations
		0x891ea6d6533ffff, // San Francisco
		0x891f0ab9207ffff, // New York
	}

	for i, h3 := range testCases {
		t.Run(fmt.Sprintf("systematic_case_%d", i), func(t *testing.T) {
			// Test Go implementation
			goResult := isValidCell(h3)

			// Test C implementation
			cResult := isValidCellC(h3)

			// Compare results (convert bool to int for comparison)
			if goResult != cResult {
				t.Errorf("Result mismatch for 0x%x: Go=%t (%v), C=%v",
					uint64(h3), goResult, goResult, cResult)
			}

			// Log the result without assuming validity
			t.Logf("H3Index 0x%x -> Go=%t, C=%v", uint64(h3), goResult, cResult)
		})
	}
}
