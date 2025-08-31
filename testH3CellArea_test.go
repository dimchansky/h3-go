// Tests ported from testH3CellArea.c
package h3

import (
	"math"
	"testing"
)

var areasKm2 = []float64{
	2.562182162955496e+06, 4.476842017201860e+05, 6.596162242711056e+04,
	9.228872919002590e+03, 1.318694490797110e+03, 1.879593512281298e+02,
	2.687164354763186e+01, 3.840848847060638e+00, 5.486939641329893e-01,
	7.838600808637444e-02, 1.119834221989390e-02, 1.599777169186614e-03,
	2.285390931423380e-04, 3.264850232091780e-05, 4.664070326136774e-06,
	6.662957615868888e-07,
}

func TestSpecificCellArea(t *testing.T) {
	t.Parallel()
	
	gc := LatLng{0.0, 0.0}
	for res := int32(0); res <= MAX_H3_RES-1; res++ {
		var cell H3Index
		err := latLngToCell(&gc, res, &cell)
		if err != E_SUCCESS {
			t.Fatalf("latLngToCell failed for resolution %d: %v", res, err)
		}
		
		area, err := cellAreaKm2(cell)
		if err != E_SUCCESS {
			t.Fatalf("cellAreaKm2 failed for resolution %d: %v", res, err)
		}
		
		if math.Abs(area-areasKm2[res]) >= 1e-8 {
			t.Errorf("cell area should match expectation for resolution %d: got %e, want %e", res, area, areasKm2[res])
		}
	}
}

func TestCellAreaInvalid(t *testing.T) {
	t.Parallel()
	
	invalid := H3Index(0xFFFFFFFFFFFFFFFF)
	
	// Test cellAreaRads2 with invalid input
	_, err := cellAreaRads2(invalid)
	if err != E_CELL_INVALID {
		t.Errorf("cellAreaRads2 invalid input: got %v, want %v", err, E_CELL_INVALID)
	}
	
	// Test cellAreaKm2 with invalid input
	_, err = cellAreaKm2(invalid)
	if err != E_CELL_INVALID {
		t.Errorf("cellAreaKm2 invalid input: got %v, want %v", err, E_CELL_INVALID)
	}
	
	// Test cellAreaM2 with invalid input
	_, err = cellAreaM2(invalid)
	if err != E_CELL_INVALID {
		t.Errorf("cellAreaM2 invalid input: got %v, want %v", err, E_CELL_INVALID)
	}
}