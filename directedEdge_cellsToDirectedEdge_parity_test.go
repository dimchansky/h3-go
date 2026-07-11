//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_cellsToDirectedEdge_parity(t *testing.T) {
	tests := []struct {
		name        string
		origin      H3Index
		destination H3Index
	}{
		// Valid neighbor pairs - these should work
		// Using some constructed valid cell pairs that are neighbors
		{"valid_neighbors_1", H3Index(0x8001fffffffffff), H3Index(0x8002fffffffffff)},
		{"valid_neighbors_2", H3Index(0x8101fffffffffff), H3Index(0x8102fffffffffff)},
		{"valid_neighbors_3", H3Index(0x8201fffffffffff), H3Index(0x8202fffffffffff)},

		// Non-neighbor cases - these should return E_NOT_NEIGHBORS
		{"same_cell", H3Index(0x8001fffffffffff), H3Index(0x8001fffffffffff)},
		{"non_neighbors", H3Index(0x8001fffffffffff), H3Index(0x8003fffffffffff)},
		{"different_res", H3Index(0x8001fffffffffff), H3Index(0x8101fffffffffff)},

		// Invalid cell cases - these should fail with appropriate errors
		{"invalid_origin", H3Index(0x0), H3Index(0x8001fffffffffff)},
		{"invalid_destination", H3Index(0x8001fffffffffff), H3Index(0x0)},
		{"both_invalid", H3Index(0x0), H3Index(0x0)},

		// Edge cases
		{"origin_invalid_mode", H3Index(0x2001fffffffffff), H3Index(0x8001fffffffffff)}, // origin has directed edge mode
		{"dest_invalid_mode", H3Index(0x8001fffffffffff), H3Index(0x2001fffffffffff)},   // destination has directed edge mode

		// Zero values
		{"zero_origin", H3Index(0), H3Index(0x8001fffffffffff)},
		{"zero_destination", H3Index(0x8001fffffffffff), H3Index(0)},
		{"both_zero", H3Index(0), H3Index(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cEdge, cErr := cellsToDirectedEdgeC(tt.origin, tt.destination)

			// Get Go implementation result
			goEdge, goErr := cellsToDirectedEdge(tt.origin, tt.destination)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch: C=%v, Go=%v", cErr, goErr)
				return
			}

			// If there was an error, we're done (no need to compare outputs)
			if cErr != E_SUCCESS {
				return
			}

			// Compare edges (should be identical)
			if cEdge != goEdge {
				t.Errorf("Edge mismatch: C=0x%x, Go=0x%x", cEdge, goEdge)
			}

			// Additional validation - if successful, the result should be a directed edge
			if getMode(goEdge) != H3_DIRECTEDEDGE_MODE {
				t.Errorf("Result is not a directed edge: mode=%d", getMode(goEdge))
			}
		})
	}
}

func Test_cellsToDirectedEdge_known_neighbors_parity(t *testing.T) {
	// Test with some known neighbor pairs - simplified without gridDisk dependency
	// These are just test cases that may or may not be actual neighbors
	// The key is that both C and Go implementations should behave identically

	testPairs := []struct {
		name        string
		origin      H3Index
		destination H3Index
	}{
		{"pair_1", H3Index(0x8007fffffffffff), H3Index(0x8009fffffffffff)},
		{"pair_2", H3Index(0x8009fffffffffff), H3Index(0x800bfffffffffff)},
		{"pair_3", H3Index(0x800bfffffffffff), H3Index(0x8007fffffffffff)},
		{"pair_4", H3Index(0x8001fffffffffff), H3Index(0x8003fffffffffff)},
		{"pair_5", H3Index(0x8003fffffffffff), H3Index(0x8005fffffffffff)},
	}

	for _, tt := range testPairs {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cEdge, cErr := cellsToDirectedEdgeC(tt.origin, tt.destination)

			// Get Go implementation result
			goEdge, goErr := cellsToDirectedEdge(tt.origin, tt.destination)

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

			// Compare edges
			if cEdge != goEdge {
				t.Errorf("Edge mismatch for origin=0x%x, destination=0x%x: C=0x%x, Go=0x%x",
					tt.origin, tt.destination, cEdge, goEdge)
			}

			// Verify the edge is valid
			if getMode(goEdge) != H3_DIRECTEDEDGE_MODE {
				t.Errorf("Result is not a directed edge: mode=%d", getMode(goEdge))
			}

			// Verify the origin matches by extracting it from the edge
			extractedOrigin, extractErr := getDirectedEdgeOrigin(goEdge)
			if extractErr != E_SUCCESS {
				t.Errorf("Failed to extract origin from edge: %v", extractErr)
			} else if extractedOrigin != tt.origin {
				t.Errorf("Extracted origin doesn't match: expected=0x%x, got=0x%x", tt.origin, extractedOrigin)
			}
		})
	}
}

func Test_cellsToDirectedEdge_pentagon_parity(t *testing.T) {
	// Test with pentagon base cells and some potential neighbors
	// Simplified test - focus on ensuring C and Go implementations match

	pentagonPairs := []struct {
		name        string
		origin      H3Index
		destination H3Index
	}{
		// Pentagon base cells with various potential destinations
		{"pentagon_1_to_2", H3Index(0x8001fffffffffff), H3Index(0x8002fffffffffff)},
		{"pentagon_1_to_3", H3Index(0x8001fffffffffff), H3Index(0x8003fffffffffff)},
		{"pentagon_3_to_1", H3Index(0x8003fffffffffff), H3Index(0x8001fffffffffff)},
		{"pentagon_3_to_5", H3Index(0x8003fffffffffff), H3Index(0x8005fffffffffff)},
		{"pentagon_7_to_9", H3Index(0x8007fffffffffff), H3Index(0x8009fffffffffff)},
	}

	for _, tt := range pentagonPairs {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cEdge, cErr := cellsToDirectedEdgeC(tt.origin, tt.destination)

			// Get Go implementation result
			goEdge, goErr := cellsToDirectedEdge(tt.origin, tt.destination)

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

			// Compare edges
			if cEdge != goEdge {
				t.Errorf("Edge mismatch for origin=0x%x, destination=0x%x: C=0x%x, Go=0x%x",
					tt.origin, tt.destination, cEdge, goEdge)
			}
		})
	}
}
