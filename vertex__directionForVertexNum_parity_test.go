//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_directionForVertexNum_parity(t *testing.T) {
	tests := []struct {
		name      string
		origin    H3Index
		vertexNum int32
	}{
		// Test with regular hexagon cells at various resolutions
		{"res0 hex", 0x8001fffffffffff, 0},
		{"res0 hex", 0x8001fffffffffff, 1},
		{"res0 hex", 0x8001fffffffffff, 2},
		{"res0 hex", 0x8001fffffffffff, 3},
		{"res0 hex", 0x8001fffffffffff, 4},
		{"res0 hex", 0x8001fffffffffff, 5},

		// Test with a resolution 1 cell
		{"res1 hex", 0x8101fffffffffff, 0},
		{"res1 hex", 0x8101fffffffffff, 3},
		{"res1 hex", 0x8101fffffffffff, 5},

		// Test with a resolution 5 cell
		{"res5 hex", 0x85283473fffffff, 0},
		{"res5 hex", 0x85283473fffffff, 1},
		{"res5 hex", 0x85283473fffffff, 2},

		// Test with pentagon cells at resolution 0 (base cell 4 is a pentagon)
		{"res0 pent", 0x8004fffffffffff, 0},
		{"res0 pent", 0x8004fffffffffff, 1},
		{"res0 pent", 0x8004fffffffffff, 2},
		{"res0 pent", 0x8004fffffffffff, 3},
		{"res0 pent", 0x8004fffffffffff, 4},

		// Test with pentagon at higher resolution
		{"res1 pent", 0x8104fffffffffff, 0},
		{"res1 pent", 0x8104fffffffffff, 2},
		{"res1 pent", 0x8104fffffffffff, 4},

		// Test invalid vertex numbers (should return INVALID_DIGIT)
		{"invalid vertex -1", 0x8001fffffffffff, -1},
		{"invalid vertex 6 for hex", 0x8001fffffffffff, 6},
		{"invalid vertex 5 for pent", 0x8004fffffffffff, 5},

		// Test with some generated cells at various resolutions
		{"res2 cell", 0x8201fffffffffff, 0},
		{"res3 cell", 0x8301fffffffffff, 1},
		{"res4 cell", 0x8401fffffffffff, 2},
		{"res6 cell", 0x8601fffffffffff, 3},
		{"res7 cell", 0x8701fffffffffff, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test parity between Go and C implementations
			goResult := directionForVertexNum(tt.origin, tt.vertexNum)
			cResult := directionForVertexNumC(tt.origin, tt.vertexNum)

			if goResult != cResult {
				t.Errorf("directionForVertexNum(%#x, %d): Go=%d, C=%d",
					tt.origin, tt.vertexNum, goResult, cResult)
			}

			// Additional validation - if the result is not INVALID_DIGIT,
			// it should be a valid direction
			if goResult != INVALID_DIGIT {
				if goResult < 0 || goResult >= NUM_DIGITS {
					t.Errorf("directionForVertexNum(%#x, %d) returned invalid direction %d",
						tt.origin, tt.vertexNum, goResult)
				}
			}
		})
	}
}

func Test_directionForVertexNum_invalidCells_parity(t *testing.T) {
	// Test with invalid H3 index
	invalidCells := []H3Index{
		0x0000000000000000, // Invalid: mode 0
		0x1000000000000000, // Invalid: mode 1
		0x2000000000000000, // Invalid: mode 2
		0xffffffffffffffff, // Invalid: all bits set
	}

	for _, cell := range invalidCells {
		for vertexNum := int32(0); vertexNum < 6; vertexNum++ {
			goResult := directionForVertexNum(cell, vertexNum)
			cResult := directionForVertexNumC(cell, vertexNum)

			if goResult != cResult {
				t.Errorf("directionForVertexNum(%#x, %d): Go=%d, C=%d",
					cell, vertexNum, goResult, cResult)
			}
		}
	}
}
