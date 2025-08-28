//go:build cgo

package h3

import (
	"testing"
)

func Test_areNeighborCells_parity(t *testing.T) {
	tests := []struct {
		name        string
		origin      H3Index
		destination H3Index
	}{
		// Same cell cases
		{"same_cell", H3Index(0x8001fffffffffff), H3Index(0x8001fffffffffff)},
		{"same_cell_2", H3Index(0x8007fffffffffff), H3Index(0x8007fffffffffff)},

		// Potential neighbor pairs - these may or may not be neighbors
		// The key is that C and Go implementations should match
		{"potential_neighbors_1", H3Index(0x8001fffffffffff), H3Index(0x8002fffffffffff)},
		{"potential_neighbors_2", H3Index(0x8002fffffffffff), H3Index(0x8001fffffffffff)},
		{"potential_neighbors_3", H3Index(0x8007fffffffffff), H3Index(0x8009fffffffffff)},
		{"potential_neighbors_4", H3Index(0x8009fffffffffff), H3Index(0x800bfffffffffff)},

		// Different resolution cases - should return E_RES_MISMATCH
		{"different_res_1", H3Index(0x8001fffffffffff), H3Index(0x81283ffffffffff)},
		{"different_res_2", H3Index(0x81283ffffffffff), H3Index(0x8228bffffffffff)},
		{"different_res_3", H3Index(0x8007fffffffffff), H3Index(0x8228bffffffffff)},

		// Invalid cell mode cases - should return E_CELL_INVALID
		{"invalid_origin_mode", H3Index(0x2001fffffffffff), H3Index(0x8001fffffffffff)}, // directed edge mode
		{"invalid_dest_mode", H3Index(0x8001fffffffffff), H3Index(0x2001fffffffffff)},   // directed edge mode
		{"both_invalid_mode", H3Index(0x2001fffffffffff), H3Index(0x2002fffffffffff)},

		// Zero/invalid values
		{"zero_origin", H3Index(0), H3Index(0x8001fffffffffff)},
		{"zero_destination", H3Index(0x8001fffffffffff), H3Index(0)},
		{"both_zero", H3Index(0), H3Index(0)},

		// Some resolution 1 cells for variety (using more realistic values)
		{"res1_same", H3Index(0x81283ffffffffff), H3Index(0x81283ffffffffff)},
		{"res1_potential_neighbors", H3Index(0x81283ffffffffff), H3Index(0x81287ffffffffff)},

		// Some resolution 2 cells
		{"res2_potential_neighbors", H3Index(0x8228bffffffffff), H3Index(0x8228dffffffffff)},
		{"res2_non_neighbors", H3Index(0x8228bffffffffff), H3Index(0x82291ffffffffff)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cResult, cErr := areNeighborCellsC(tt.origin, tt.destination)

			// Get Go implementation result
			goResult, goErr := areNeighborCells(tt.origin, tt.destination)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch for origin=0x%x, destination=0x%x: C=%v, Go=%v",
					tt.origin, tt.destination, cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != E_SUCCESS {
				return
			}

			// Compare boolean results
			if cResult != goResult {
				t.Errorf("Result mismatch for origin=0x%x, destination=0x%x: C=%v, Go=%v",
					tt.origin, tt.destination, cResult, goResult)
			}
		})
	}
}

func Test_areNeighborCells_known_valid_neighbors_parity(t *testing.T) {
	// Test with some constructed H3 indexes that should be neighbors
	// This is a more targeted test using specific constructed values

	validNeighborTests := []struct {
		name        string
		origin      H3Index
		destination H3Index
	}{
		// These are constructed to test specific behavior patterns
		// Based on resolution 0 base cells which should have known relationships
		{"base_cell_test_1", H3Index(0x8001fffffffffff), H3Index(0x8003fffffffffff)},
		{"base_cell_test_2", H3Index(0x8003fffffffffff), H3Index(0x8005fffffffffff)},
		{"base_cell_test_3", H3Index(0x8005fffffffffff), H3Index(0x8007fffffffffff)},
		{"base_cell_test_4", H3Index(0x8007fffffffffff), H3Index(0x8009fffffffffff)},
		{"base_cell_test_5", H3Index(0x8009fffffffffff), H3Index(0x800bfffffffffff)},

		// Test reverse direction
		{"base_cell_reverse_1", H3Index(0x8003fffffffffff), H3Index(0x8001fffffffffff)},
		{"base_cell_reverse_2", H3Index(0x8005fffffffffff), H3Index(0x8003fffffffffff)},

		// Test some higher resolution cells
		{"res1_test_1", H3Index(0x81283ffffffffff), H3Index(0x81287ffffffffff)},
		{"res1_test_2", H3Index(0x81287ffffffffff), H3Index(0x8128bffffffffff)},

		// Test known non-neighbor pairs
		{"non_neighbor_1", H3Index(0x8001fffffffffff), H3Index(0x8009fffffffffff)},
		{"non_neighbor_2", H3Index(0x8003fffffffffff), H3Index(0x800bfffffffffff)},
	}

	for _, tt := range validNeighborTests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cResult, cErr := areNeighborCellsC(tt.origin, tt.destination)

			// Get Go implementation result
			goResult, goErr := areNeighborCells(tt.origin, tt.destination)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch for origin=0x%x, destination=0x%x: C=%v, Go=%v",
					tt.origin, tt.destination, cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != E_SUCCESS {
				return
			}

			// Compare boolean results
			if cResult != goResult {
				t.Errorf("Result mismatch for origin=0x%x, destination=0x%x: C=%v, Go=%v",
					tt.origin, tt.destination, cResult, goResult)
			}

			// Log the result for manual verification
			t.Logf("Origin=0x%x, Destination=0x%x, AreNeighbors=%v",
				tt.origin, tt.destination, goResult)
		})
	}
}

func Test_areNeighborCells_pentagon_edge_cases_parity(t *testing.T) {
	// Test pentagon-specific edge cases
	// These test cases focus on pentagon cells and their special behavior

	pentagonTests := []struct {
		name        string
		origin      H3Index
		destination H3Index
	}{
		// Pentagon base cells (base cells 4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117)
		// Test with some constructed pentagon cells
		{"pentagon_case_1", H3Index(0x8009fffffffffff), H3Index(0x800bfffffffffff)},
		{"pentagon_case_2", H3Index(0x800bfffffffffff), H3Index(0x8009fffffffffff)},
		{"pentagon_case_3", H3Index(0x8001fffffffffff), H3Index(0x8007fffffffffff)},

		// Test pentagon with itself
		{"pentagon_self_1", H3Index(0x8009fffffffffff), H3Index(0x8009fffffffffff)},
		{"pentagon_self_2", H3Index(0x800bfffffffffff), H3Index(0x800bfffffffffff)},

		// Test potential invalid K_AXES cases with pentagons
		{"pentagon_k_axes_test_1", H3Index(0x8001fffffffffff), H3Index(0x8009fffffffffff)},
		{"pentagon_k_axes_test_2", H3Index(0x8009fffffffffff), H3Index(0x8001fffffffffff)},
	}

	for _, tt := range pentagonTests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cResult, cErr := areNeighborCellsC(tt.origin, tt.destination)

			// Get Go implementation result
			goResult, goErr := areNeighborCells(tt.origin, tt.destination)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch for origin=0x%x, destination=0x%x: C=%v, Go=%v",
					tt.origin, tt.destination, cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != E_SUCCESS {
				return
			}

			// Compare boolean results
			if cResult != goResult {
				t.Errorf("Result mismatch for origin=0x%x, destination=0x%x: C=%v, Go=%v",
					tt.origin, tt.destination, cResult, goResult)
			}

			// Log pentagon test results
			t.Logf("Pentagon test: Origin=0x%x, Destination=0x%x, AreNeighbors=%v",
				tt.origin, tt.destination, goResult)
		})
	}
}
