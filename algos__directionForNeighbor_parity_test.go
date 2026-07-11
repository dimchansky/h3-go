//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_directionForNeighbor_parity(t *testing.T) {
	testCases := []struct {
		name        string
		origin      h3Index
		destination h3Index
	}{
		{
			name:        "hex neighbor",
			origin:      0x8930062838bffff, // Valid hex cell
			destination: 0x8930062838affff, // Adjacent hex cell
		},
		{
			name:        "pentagon neighbor",
			origin:      0x893006283807fff, // Pentagon cell
			destination: 0x8930062838b7fff, // Adjacent to pentagon
		},
		{
			name:        "invalid destination",
			origin:      0x8930062838bffff,
			destination: 0x8930062838dffff, // Not a direct neighbor
		},
		{
			name:        "resolution 0 hex",
			origin:      0x8001fffffffffff, // Base cell 0
			destination: 0x8007fffffffffff, // Base cell 1 (neighbor)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call Go implementation
			goResult := _directionForNeighbor(tc.origin, tc.destination)

			// Call C implementation
			cResult := directionForNeighborC(tc.origin, tc.destination)

			if goResult != cResult {
				t.Errorf("_directionForNeighbor(%#016x, %#016x) = %d, want %d",
					tc.origin, tc.destination, goResult, cResult)
			}
		})
	}
}
