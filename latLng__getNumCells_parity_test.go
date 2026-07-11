//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_getNumCells_parity(t *testing.T) {
	tests := []struct {
		name string
		res  int32
	}{
		// Valid resolution range
		{"res0", 0},
		{"res1", 1},
		{"res2", 2},
		{"res5", 5},
		{"res7", 7},
		{"res10", 10},
		{"res15", 15}, // MAX_H3_RES

		// Invalid resolution values (should return E_RES_DOMAIN)
		{"negative_res", -1},
		{"res_too_high", 16},
		{"very_negative", -10},
		{"much_too_high", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cNumCells, cErr := getNumCellsC(tt.res)

			// Get Go implementation result
			goNumCells, goErr := getNumCells(tt.res)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch: C=%v, Go=%v", cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != E_SUCCESS {
				return
			}

			// Compare exact integer values (should be identical)
			if cNumCells != goNumCells {
				t.Errorf("NumCells mismatch for res %d: C=%d, Go=%d",
					tt.res, cNumCells, goNumCells)
			}
		})
	}
}
