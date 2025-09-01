// Tests ported from testCellToLocalIj.c
package h3

import (
	"math"
	"testing"
)

func Test_ijBaseCells(t *testing.T) {
	t.Parallel()
	
	var ij CoordIJ
	origin := H3Index(0x8029fffffffffff)
	
	// Test getting origin back
	ij = CoordIJ{I: 0, J: 0}
	var retrieved H3Index
	err := localIjToCell(origin, &ij, 0, &retrieved)
	if err != E_SUCCESS {
		t.Fatalf("got origin back failed: %v", err)
	}
	if retrieved != 0x8029fffffffffff {
		t.Errorf("origin matches self: expected 0x8029fffffffffff, got 0x%x", retrieved)
	}
	
	// Test offset index
	ij.I = 1
	err = localIjToCell(origin, &ij, 0, &retrieved)
	if err != E_SUCCESS {
		t.Fatalf("got offset index failed: %v", err)
	}
	if retrieved != 0x8051fffffffffff {
		t.Errorf("modified index matches expected: expected 0x8051fffffffffff, got 0x%x", retrieved)
	}
	
	// Test out of range base cell (1)
	ij.I = 2
	err = localIjToCell(origin, &ij, 0, &retrieved)
	if err == E_SUCCESS {
		t.Error("out of range base cell (1) should fail")
	}
	
	// Test out of range base cell (2)
	ij.I = 0
	ij.J = 2
	err = localIjToCell(origin, &ij, 0, &retrieved)
	if err == E_SUCCESS {
		t.Error("out of range base cell (2) should fail")
	}
	
	// Test out of range base cell (3)
	ij.I = -2
	ij.J = -2
	err = localIjToCell(origin, &ij, 0, &retrieved)
	if err == E_SUCCESS {
		t.Error("out of range base cell (3) should fail")
	}
}

func Test_ijOutOfRange(t *testing.T) {
	t.Parallel()
	
	coords := []CoordIJ{
		{0, 0}, {1, 0}, {2, 0}, {3, 0},
		{4, 0}, {-4, 0}, {0, 4},
	}
	expected := []H3Index{
		0x81283ffffffffff,
		0x81293ffffffffff,
		0x8150bffffffffff,
		0x8151bffffffffff,
		H3_NULL,
		H3_NULL,
		H3_NULL,
	}
	
	for i, coord := range coords {
		var result H3Index
		err := localIjToCell(expected[0], &coord, 0, &result)
		if expected[i] == H3_NULL {
			if err == E_SUCCESS {
				t.Errorf("coordinates out of range at index %d", i)
			}
		} else {
			if err != E_SUCCESS {
				t.Errorf("coordinates in range at index %d failed: %v", i, err)
			}
			if result != expected[i] {
				t.Errorf("result matches expectation at index %d: expected 0x%x, got 0x%x", i, expected[i], result)
			}
		}
	}
}

func Test_cellToLocalIjFailed(t *testing.T) {
	t.Parallel()
	
	// Some indexes that represent base cells. Base cells
	// are hexagons except for `pent1`.
	var bc1 H3Index
	setH3Index(&bc1, 0, 15, 0)
	
	var bc2 H3Index
	setH3Index(&bc2, 0, 8, 0)
	
	var bc3 H3Index
	setH3Index(&bc3, 0, 31, 0)
	
	var pent1 H3Index
	setH3Index(&pent1, 0, 4, 0)
	
	var ij CoordIJ
	
	// Test bc1 to bc1
	err := cellToLocalIj(bc1, bc1, 0, &ij)
	if err != E_SUCCESS {
		t.Errorf("found IJ (1) failed: %v", err)
	}
	if ij.I != 0 || ij.J != 0 {
		t.Errorf("ij correct (1): expected (0,0), got (%d,%d)", ij.I, ij.J)
	}
	
	// Test bc1 to pent1
	err = cellToLocalIj(bc1, pent1, 0, &ij)
	if err != E_SUCCESS {
		t.Errorf("found IJ (2) failed: %v", err)
	}
	if ij.I != 1 || ij.J != 0 {
		t.Errorf("ij correct (2): expected (1,0), got (%d,%d)", ij.I, ij.J)
	}
	
	// Test bc1 to bc2
	err = cellToLocalIj(bc1, bc2, 0, &ij)
	if err != E_SUCCESS {
		t.Errorf("found IJ (3) failed: %v", err)
	}
	if ij.I != 0 || ij.J != -1 {
		t.Errorf("ij correct (3): expected (0,-1), got (%d,%d)", ij.I, ij.J)
	}
	
	// Test bc1 to bc3
	err = cellToLocalIj(bc1, bc3, 0, &ij)
	if err != E_SUCCESS {
		t.Errorf("found IJ (4) failed: %v", err)
	}
	if ij.I != -1 || ij.J != 0 {
		t.Errorf("ij correct (4): expected (-1,0), got (%d,%d)", ij.I, ij.J)
	}
	
	// Test pent1 to bc3 - should fail
	err = cellToLocalIj(pent1, bc3, 0, &ij)
	if err == E_SUCCESS {
		t.Error("found IJ (5) should fail")
	}
}

func Test_cellToLocalIjInvalid(t *testing.T) {
	t.Parallel()
	
	var bc1 H3Index
	setH3Index(&bc1, 0, 15, 0)
	
	var ij CoordIJ
	
	// Test invalid index
	invalidIndex := setResolution(H3Index(0x7fffffffffffffff), getResolution(bc1))
	err := cellToLocalIj(bc1, invalidIndex, 0, &ij)
	if err != E_CELL_INVALID {
		t.Errorf("invalid index: expected E_CELL_INVALID, got %v", err)
	}
	
	// Test invalid origin
	err = cellToLocalIj(0x7fffffffffffffff, bc1, 0, &ij)
	if err != E_RES_MISMATCH {
		t.Errorf("invalid origin: expected E_RES_MISMATCH, got %v", err)
	}
	
	// Test invalid origin and index
	err = cellToLocalIj(0x7fffffffffffffff, 0x7fffffffffffffff, 0, &ij)
	if err != E_CELL_INVALID {
		t.Errorf("invalid origin and index: expected E_CELL_INVALID, got %v", err)
	}
}

func Test_localIjToCellInvalid(t *testing.T) {
	t.Parallel()
	
	ij := CoordIJ{0, 0}
	var index H3Index
	err := localIjToCell(0x7fffffffffffffff, &ij, 0, &index)
	if err != E_CELL_INVALID {
		t.Errorf("invalid origin for ijToH3: expected E_CELL_INVALID, got %v", err)
	}
}

func Test_indexOnPentInvalid(t *testing.T) {
	t.Parallel()
	
	// Tests for INVALID_DIGIT being detected and failed on in various cases.
	var onPentInvalid H3Index
	setH3Index(&onPentInvalid, 1, 4, int32(INVALID_DIGIT))
	
	var offPent H3Index
	setH3Index(&offPent, 1, 3, int32(CENTER_DIGIT))
	
	var ij CoordIJ
	err := cellToLocalIj(offPent, onPentInvalid, 0, &ij)
	if err != E_CELL_INVALID {
		t.Errorf("invalid index on pentagon: expected E_CELL_INVALID, got %v", err)
	}
	
	var onPentValid H3Index
	setH3Index(&onPentValid, 1, 4, int32(CENTER_DIGIT))
	
	err = cellToLocalIj(onPentInvalid, onPentValid, 0, &ij)
	if err != E_CELL_INVALID {
		t.Errorf("invalid both on pentagon (1): expected E_CELL_INVALID, got %v", err)
	}
	
	err = cellToLocalIj(onPentValid, onPentInvalid, 0, &ij)
	if err != E_CELL_INVALID {
		t.Errorf("invalid both on pentagon (2): expected E_CELL_INVALID, got %v", err)
	}
	
	ij.I = 0
	ij.J = 0
	var out H3Index
	err = localIjToCell(onPentInvalid, &ij, 0, &out)
	if err != E_CELL_INVALID {
		t.Errorf("invalid both on pentagon (3): expected E_CELL_INVALID, got %v", err)
	}
	
	ij.I = 3
	ij.J = 3
	err = localIjToCell(onPentInvalid, &ij, 0, &out)
	if err != E_CELL_INVALID {
		t.Errorf("invalid origin on pentagon: expected E_CELL_INVALID, got %v", err)
	}
}

func Test_invalidMode(t *testing.T) {
	t.Parallel()
	
	var ij CoordIJ
	cell := H3Index(0x85283473fffffff)
	
	// Test valid mode first
	err := cellToLocalIj(cell, cell, 0, &ij)
	if err != E_SUCCESS {
		t.Fatalf("valid mode should succeed: %v", err)
	}
	
	// Test invalid modes 1-32
	for i := uint32(1); i <= 32; i++ {
		var ij2 CoordIJ
		err := cellToLocalIj(cell, cell, i, &ij2)
		if err != E_OPTION_INVALID {
			t.Errorf("Invalid mode fail for cellToLocalIj at mode %d: expected E_OPTION_INVALID, got %v", i, err)
		}
		
		var cell2 H3Index
		err = localIjToCell(cell, &ij2, i, &cell2)
		if err != E_OPTION_INVALID {
			t.Errorf("Invalid mode fail for localIjToCell at mode %d: expected E_OPTION_INVALID, got %v", i, err)
		}
	}
}

func Test_invalid_negativeIj(t *testing.T) {
	t.Parallel()
	
	index := H3Index(0x200f202020202020)
	ij := CoordIJ{I: -14671840, J: math.MinInt32}
	var out H3Index
	err := localIjToCell(index, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("Negative I and J components fail")
	}
}

func Test_localIjToCell_overflow_i(t *testing.T) {
	t.Parallel()
	
	var origin H3Index
	setH3Index(&origin, 2, 2, int32(CENTER_DIGIT))
	ij := CoordIJ{I: math.MinInt32, J: math.MaxInt32}
	var out H3Index
	err := localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("High magnitude I and J components fail")
	}
}

func Test_localIjToCell_overflow_j(t *testing.T) {
	t.Parallel()
	
	var origin H3Index
	setH3Index(&origin, 2, 2, int32(CENTER_DIGIT))
	ij := CoordIJ{I: math.MaxInt32, J: math.MinInt32}
	var out H3Index
	err := localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("High magnitude J and I components fail")
	}
}

func Test_localIjToCell_overflow_ij(t *testing.T) {
	t.Parallel()
	
	var origin H3Index
	setH3Index(&origin, 2, 2, int32(CENTER_DIGIT))
	ij := CoordIJ{I: math.MinInt32, J: math.MinInt32}
	var out H3Index
	err := localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("High magnitude J and I components fail")
	}
}

func Test_localIjToCell_overflow_particularCases(t *testing.T) {
	t.Parallel()
	
	var origin H3Index
	setH3Index(&origin, 2, 2, int32(CENTER_DIGIT))
	
	var originRes3 H3Index
	setH3Index(&originRes3, 2, 2, int32(CENTER_DIGIT))
	
	var out H3Index
	
	// Test case 1
	ij := CoordIJ{I: 553648127, J: -2145378272}
	err := localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("Particular high magnitude J and I components fail (1)")
	}
	
	// Test case 2
	ij.I = math.MaxInt32 - 10
	ij.J = -11
	err = localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("Particular high magnitude J and I components fail (2)")
	}
	
	// Test case 3
	ij.I = 553648127
	ij.J = -2145378272
	err = localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("Particular high magnitude J and I components fail (3)")
	}
	
	// Test case 4
	ij.I = math.MaxInt32 - 10
	ij.J = -10
	err = localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("Particular high magnitude J and I components fail (4)")
	}
	
	// Test case 5
	ij.I = math.MaxInt32 - 10
	ij.J = -9
	err = localIjToCell(origin, &ij, 0, &out)
	if err == E_SUCCESS {
		t.Error("Particular high magnitude J and I components fail (5)")
	}
}