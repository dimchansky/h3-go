//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_vertexRotations_parity(t *testing.T) {
	// Test with various cells including hexagons and pentagons
	testCells := []H3Index{
		0x85283473fffffff, // res 5 hexagon
		0x8928308280fffff, // res 9 hexagon
		0x8a2a1072b587fff, // res 10 hexagon
		0x8b2a1072b59ffff, // res 11 hexagon
		0x821c07fffffffff, // res 1 pentagon
		0x841c001ffffffff, // res 4 pentagon
		0x851c003ffffffff, // res 5 pentagon
		0x891c000003fffff, // res 9 pentagon

		// Additional test cases at various resolutions
		0x8001fffffffffff, // res 0 hexagon
		0x8101fffffffffff, // res 1 hexagon
		0x8201fffffffffff, // res 2 hexagon
		0x8301fffffffffff, // res 3 hexagon
		0x8401fffffffffff, // res 4 hexagon
		0x8501fffffffffff, // res 5 hexagon
		0x8601fffffffffff, // res 6 hexagon
		0x8701fffffffffff, // res 7 hexagon
		0x8801fffffffffff, // res 8 hexagon
		0x8901fffffffffff, // res 9 hexagon
		0x8a01fffffffffff, // res 10 hexagon
		0x8b01fffffffffff, // res 11 hexagon
		0x8c01fffffffffff, // res 12 hexagon
		0x8d01fffffffffff, // res 13 hexagon
		0x8e01fffffffffff, // res 14 hexagon
		0x8f01fffffffffff, // res 15 hexagon
	}

	for _, cell := range testCells {
		if !isValidCell(cell) {
			continue
		}

		var goOut, cOut int32
		goErr := vertexRotations(cell, &goOut)
		cErr := vertexRotationsC(cell, &cOut)

		if goErr != cErr {
			t.Errorf("vertexRotations(0x%x) error mismatch: Go=%v, C=%v", cell, goErr, cErr)
			continue
		}

		if goErr == E_SUCCESS && goOut != cOut {
			t.Errorf("vertexRotations(0x%x) output mismatch: Go=%d, C=%d", cell, goOut, cOut)
		}
	}
}

func Test_vertexRotations_parity_pentagons(t *testing.T) {
	// Test all pentagons at different resolutions
	for res := 0; res <= 15; res++ {
		pentagons := make([]H3Index, NUM_PENTAGONS)
		getPentagons(int32(res), pentagons)

		for _, pentagon := range pentagons {
			var goOut, cOut int32
			goErr := vertexRotations(pentagon, &goOut)
			cErr := vertexRotationsC(pentagon, &cOut)

			if goErr != cErr {
				t.Errorf("vertexRotations(pentagon 0x%x res=%d) error mismatch: Go=%v, C=%v",
					pentagon, res, goErr, cErr)
				continue
			}

			if goErr == E_SUCCESS && goOut != cOut {
				t.Errorf("vertexRotations(pentagon 0x%x res=%d) output mismatch: Go=%d, C=%d",
					pentagon, res, goOut, cOut)
			}
		}
	}
}

func Test_vertexRotations_parity_neighbors(t *testing.T) {
	// Test cells and their neighbors to exercise the rotation logic
	baseCells := []H3Index{
		0x85283473fffffff, // res 5 hexagon
		0x8528342bfffffff, // res 5 hexagon
		0x85283447fffffff, // res 5 hexagon
		0x85283463fffffff, // res 5 hexagon
		0x8528347bfffffff, // res 5 hexagon
	}

	for _, cell := range baseCells {
		if !isValidCell(cell) {
			continue
		}

		// Test the cell itself
		var goOut, cOut int32
		goErr := vertexRotations(cell, &goOut)
		cErr := vertexRotationsC(cell, &cOut)

		if goErr != cErr {
			t.Errorf("vertexRotations(0x%x) error mismatch: Go=%v, C=%v", cell, goErr, cErr)
			continue
		}

		if goErr == E_SUCCESS && goOut != cOut {
			t.Errorf("vertexRotations(0x%x) output mismatch: Go=%d, C=%d", cell, goOut, cOut)
			continue
		}

		// Test with some additional cells for coverage
		additionalCells := []H3Index{
			cell + 1, // Adjacent cell
			cell - 1, // Adjacent cell
		}

		for _, additionalCell := range additionalCells {
			if !isValidCell(additionalCell) {
				continue
			}

			var goOutAdd, cOutAdd int32
			goErrAdd := vertexRotations(additionalCell, &goOutAdd)
			cErrAdd := vertexRotationsC(additionalCell, &cOutAdd)

			if goErrAdd != cErrAdd {
				t.Errorf("vertexRotations(additional 0x%x) error mismatch: Go=%v, C=%v",
					additionalCell, goErrAdd, cErrAdd)
				continue
			}

			if goErrAdd == E_SUCCESS && goOutAdd != cOutAdd {
				t.Errorf("vertexRotations(additional 0x%x) output mismatch: Go=%d, C=%d",
					additionalCell, goOutAdd, cOutAdd)
			}
		}
	}
}
