//go:build cgo

package h3

import "testing"

func TestGetRes0CellsParity(t *testing.T) {
	// Test the Go implementation against C reference
	goResult := make([]H3Index, NUM_BASE_CELLS)
	goErr := getRes0Cells(goResult)

	cResult := make([]H3Index, NUM_BASE_CELLS)
	cErr := getRes0CellsC(cResult)

	// Check error codes match
	if goErr != cErr {
		t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
	}

	// Check all results match
	if goErr == E_SUCCESS {
		for i := 0; i < NUM_BASE_CELLS; i++ {
			if goResult[i] != cResult[i] {
				t.Errorf("Cell %d mismatch: Go=0x%x, C=0x%x", i, goResult[i], cResult[i])
			}
		}
	}
}

func TestGetRes0CellsInvalidSize(t *testing.T) {
	// Test with wrong buffer size
	wrongSize := make([]H3Index, 10)
	err := getRes0Cells(wrongSize)
	if err != E_FAILED {
		t.Errorf("Expected E_FAILED for wrong buffer size, got %v", err)
	}
}

func TestSetBaseCellParity(t *testing.T) {
	// Test setBaseCell helper function
	for bc := int32(0); bc < NUM_BASE_CELLS; bc++ {
		h := H3Index(H3_INIT)
		h = setMode(h, H3_CELL_MODE)
		h = setBaseCell(h, bc)

		// Verify the base cell was set correctly by extracting it
		extractedBC := getBaseCell(h)
		if extractedBC != bc {
			t.Errorf("Base cell %d: expected %d, got %d", bc, bc, extractedBC)
		}
	}
}
