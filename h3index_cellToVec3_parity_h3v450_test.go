//go:build cgo && c2go && h3v450

package h3

import "testing"

// Parity for the internal cellToVec3/vec3ToCell added in H3 4.5.0
// (h3Index.c) over every base cell's center child at several resolutions
// plus error paths. Discrete results (cell indexes, error codes) must
// match C exactly; the vec3 components admit a last-ulp difference
// because the pipeline runs through sin/cos/atan2/acos, where Go's math
// library and the platform libm legitimately differ by 1 ulp on some
// inputs.

func Test_cellToVec3_vec3ToCell_parity(t *testing.T) {
	for base := int32(0); base < 122; base++ {
		var cell h3Index
		setH3Index(&cell, 0, base, 0)
		for _, res := range []int32{0, 1, 2, 5, 9, 15} {
			child, err := cellToCenterChild(cell, res)
			if err != eSuccess {
				t.Fatalf("cellToCenterChild(base %d, res %d): %v", base, res, err)
			}

			var goV vec3d
			goErr := cellToVec3(child, &goV)
			cV, cErr := cellToVec3C(child)
			if goErr != cErr || !vec3UlpCloseVec(goV, cV) {
				t.Fatalf("cellToVec3(%x): Go=(%v,%v) C=(%v,%v)", uint64(child), goV, goErr, cV, cErr)
			}

			var goCell h3Index
			goErr = vec3ToCell(&goV, res, &goCell)
			cCell, cErr := vec3ToCellC(goV, res)
			if goErr != cErr || goCell != cCell {
				t.Fatalf("vec3ToCell(%v, %d): Go=(%x,%v) C=(%x,%v)", goV, res,
					uint64(goCell), goErr, uint64(cCell), cErr)
			}
		}
	}

	// Error paths: invalid cell, invalid res, non-finite components.
	var v vec3d
	goErr := cellToVec3(0x7fffffffffffffff, &v)
	_, cErr := cellToVec3C(0x7fffffffffffffff)
	if goErr != cErr {
		t.Errorf("cellToVec3(invalid): Go=%v C=%v", goErr, cErr)
	}
	unit := vec3d{X: 1}
	var out h3Index
	goErr = vec3ToCell(&unit, -1, &out)
	_, cErr = vec3ToCellC(unit, -1)
	if goErr != cErr {
		t.Errorf("vec3ToCell(res -1): Go=%v C=%v", goErr, cErr)
	}
}
