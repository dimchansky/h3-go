// Tests ported from testCellToLocalIjInternal.c
package h3

import (
	"testing"
)

func Test_ijkBaseCells(t *testing.T) {
	t.Parallel()
	
	// Some indexes that represent base cells. Base cells
	// are hexagons except for `pent1`.
	var bc1 H3Index
	setH3Index(&bc1, 0, 15, 0)
	
	var pent1 H3Index
	setH3Index(&pent1, 0, 4, 0)
	
	var ijk CoordIJK
	err := cellToLocalIjk(pent1, bc1, &ijk)
	if err != E_SUCCESS {
		t.Errorf("expected cellToLocalIjk to succeed for base cells 4 and 15, got error: %v", err)
	}
	
	if !_ijkMatches(&ijk, &UNIT_VECS[2]) {
		t.Errorf("expected neighboring base cell at 0,1,0, got %v", ijk)
	}
}

// Test that coming from the same direction outside the pentagon is handled
// the same as coming from the same direction inside the pentagon.
func Test_onOffPentagonSame(t *testing.T) {
	t.Parallel()
	
	for bc := int32(0); bc < NUM_BASE_CELLS; bc++ {
		for res := int32(1); res <= MAX_H3_RES; res++ {
			// K_AXES_DIGIT is the first internal direction, and it's also
			// invalid for pentagons, so skip to next.
			startDir := K_AXES_DIGIT
			if _isBaseCellPentagon(bc) {
				startDir++
			}
			
			for dir := startDir; dir < NUM_DIGITS; dir++ {
				var internalOrigin H3Index
				setH3Index(&internalOrigin, res, bc, int32(dir))
				
				var externalOrigin H3Index
				setH3Index(&externalOrigin, res, _getBaseCellNeighbor(bc, dir), int32(CENTER_DIGIT))
				
				for testDir := startDir; testDir < NUM_DIGITS; testDir++ {
					var testIndex H3Index
					setH3Index(&testIndex, res, bc, int32(testDir))
					
					var internalIj CoordIJ
					internalIjFailed := cellToLocalIj(internalOrigin, testIndex, 0, &internalIj)
					
					var externalIj CoordIJ
					externalIjFailed := cellToLocalIj(externalOrigin, testIndex, 0, &externalIj)
					
					// Check that both succeed or both fail
					if (internalIjFailed != E_SUCCESS) != (externalIjFailed != E_SUCCESS) {
						t.Errorf("internal/external failed mismatch when getting IJ: bc=%d, res=%d, dir=%d, testDir=%d, internal=%v, external=%v",
							bc, res, dir, testDir, internalIjFailed, externalIjFailed)
					}
					
					if internalIjFailed != E_SUCCESS {
						continue
					}
					
					var internalIndex H3Index
					internalIjFailed2 := localIjToCell(internalOrigin, &internalIj, 0, &internalIndex)
					
					var externalIndex H3Index
					externalIjFailed2 := localIjToCell(externalOrigin, &externalIj, 0, &externalIndex)
					
					// Check that both succeed or both fail
					if (internalIjFailed2 != E_SUCCESS) != (externalIjFailed2 != E_SUCCESS) {
						t.Errorf("internal/external failed mismatch when getting index: bc=%d, res=%d, dir=%d, testDir=%d, internal=%v, external=%v",
							bc, res, dir, testDir, internalIjFailed2, externalIjFailed2)
					}
					
					if internalIjFailed2 != E_SUCCESS {
						continue
					}
					
					if internalIndex != externalIndex {
						t.Errorf("internal/external index mismatch: bc=%d, res=%d, dir=%d, testDir=%d, internal=0x%x, external=0x%x",
							bc, res, dir, testDir, internalIndex, externalIndex)
					}
				}
			}
		}
	}
}