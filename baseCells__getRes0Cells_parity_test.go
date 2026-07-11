//go:build cgo && c2go

package h3

import "testing"

func TestGetRes0CellsParity(t *testing.T) {
	// Test the Go implementation against C reference
	goResult := make([]h3Index, numBaseCells)
	goErr := getRes0Cells(goResult)

	cResult := make([]h3Index, numBaseCells)
	cErr := getRes0CellsC(cResult)

	// Check error codes match
	if goErr != cErr {
		t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
	}

	// Check all results match
	if goErr == eSuccess {
		for i := 0; i < numBaseCells; i++ {
			if goResult[i] != cResult[i] {
				t.Errorf("Cell %d mismatch: Go=0x%x, C=0x%x", i, goResult[i], cResult[i])
			}
		}
	}
}

func TestGetRes0CellsInvalidSize(t *testing.T) {
	// Test with wrong buffer size
	wrongSize := make([]h3Index, 10)
	err := getRes0Cells(wrongSize)
	if err != eFailed {
		t.Errorf("Expected eFailed for wrong buffer size, got %v", err)
	}
}

func TestSetBaseCellParity(t *testing.T) {
	// Test setBaseCell helper function
	for bc := int32(0); bc < numBaseCells; bc++ {
		h := h3Index(h3Init)
		h = setMode(h, h3CellMode)
		h = setBaseCell(h, bc)

		// Verify the base cell was set correctly by extracting it
		extractedBC := getBaseCell(h)
		if extractedBC != bc {
			t.Errorf("Base cell %d: expected %d, got %d", bc, bc, extractedBC)
		}
	}
}
