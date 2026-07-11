//go:build cgo && c2go

package h3

import "testing"

func Test_res0CellCount_ParityWithC(t *testing.T) {
	// This function takes no parameters, so we just compare the results
	goResult := res0CellCount()
	cResult := res0CellCountC()

	if goResult != cResult {
		t.Fatalf("res0CellCount mismatch: go=%d c=%d", goResult, cResult)
	}

	// Verify the expected value
	expectedCount := int32(122) // numBaseCells
	if goResult != expectedCount {
		t.Errorf("Expected res0CellCount to return %d but got %d", expectedCount, goResult)
	}
}
