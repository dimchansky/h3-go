//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_cellToBBox_parity(t *testing.T) {
	tests := []struct {
		name          string
		cell          h3Index
		coverChildren bool
		expectError   h3Error
	}{
		{
			name:          "Valid res 0 cell",
			cell:          0x8001fffffffffff, // Base cell 1 at res 0
			coverChildren: false,
			expectError:   eSuccess,
		},
		{
			name:          "Valid res 0 cell with children coverage",
			cell:          0x8001fffffffffff,
			coverChildren: true,
			expectError:   eSuccess,
		},
		{
			name:          "Valid res 5 cell",
			cell:          0x851fb46622dffff, // Some res 5 cell
			coverChildren: false,
			expectError:   eSuccess,
		},
		{
			name:          "Valid res 5 cell with children coverage",
			cell:          0x851fb46622dffff,
			coverChildren: true,
			expectError:   eSuccess,
		},
		{
			name:          "North pole cell res 1",
			cell:          northPoleCells[1],
			coverChildren: false,
			expectError:   eSuccess,
		},
		{
			name:          "South pole cell res 1",
			cell:          southPoleCells[1],
			coverChildren: false,
			expectError:   eSuccess,
		},
		{
			name:          "Pentagon base cell",
			cell:          0x804dfffffffffff, // Base cell 4 (pentagon)
			coverChildren: false,
			expectError:   eSuccess,
		},
		{
			name:          "Invalid cell",
			cell:          0x0,
			coverChildren: false,
			expectError:   eCellInvalid,
		},
		{
			name:          "Another invalid cell",
			cell:          h3Null,
			coverChildren: false,
			expectError:   eCellInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Go implementation
			goBBox, goErr := cellToBBox(tt.cell, tt.coverChildren)

			// Call C implementation
			cBBox, cErr := cellToBBoxC(tt.cell, tt.coverChildren)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If both succeeded, compare the bounding boxes
			if goErr == eSuccess && cErr == eSuccess {
				tolerance := 1e-12 // Tight tolerance for coordinate comparisons

				if math.Abs(float64(goBBox.North-cBBox.North)) > tolerance {
					t.Errorf("North mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						float64(goBBox.North), float64(cBBox.North), float64(goBBox.North-cBBox.North))
				}
				if math.Abs(float64(goBBox.South-cBBox.South)) > tolerance {
					t.Errorf("South mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						float64(goBBox.South), float64(cBBox.South), float64(goBBox.South-cBBox.South))
				}
				if math.Abs(float64(goBBox.East-cBBox.East)) > tolerance {
					t.Errorf("East mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						float64(goBBox.East), float64(cBBox.East), float64(goBBox.East-cBBox.East))
				}
				if math.Abs(float64(goBBox.West-cBBox.West)) > tolerance {
					t.Errorf("West mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						float64(goBBox.West), float64(cBBox.West), float64(goBBox.West-cBBox.West))
				}
			}
		})
	}
}

func Test_cellToBBox_comprehensive_parity(t *testing.T) {
	// Test all base cells at resolution 0
	for baseCell := int32(0); baseCell < numBaseCells; baseCell++ {
		t.Run("BaseCell_"+string(rune(baseCell+'0')), func(t *testing.T) {
			cell := baseCellNumToCell(baseCell)
			if cell == h3Null {
				return // Skip invalid base cells
			}

			// Test both coverage modes
			for _, coverChildren := range []bool{false, true} {
				// Call Go implementation
				goBBox, goErr := cellToBBox(cell, coverChildren)

				// Call C implementation
				cBBox, cErr := cellToBBoxC(cell, coverChildren)

				// Compare errors
				if goErr != cErr {
					t.Errorf("BaseCell %d, coverChildren=%v: Error mismatch: Go=%v, C=%v",
						baseCell, coverChildren, goErr, cErr)
					continue
				}

				// If both succeeded, compare the bounding boxes
				if goErr == eSuccess && cErr == eSuccess {
					tolerance := 1e-12

					if math.Abs(float64(goBBox.North-cBBox.North)) > tolerance ||
						math.Abs(float64(goBBox.South-cBBox.South)) > tolerance ||
						math.Abs(float64(goBBox.East-cBBox.East)) > tolerance ||
						math.Abs(float64(goBBox.West-cBBox.West)) > tolerance {
						t.Errorf("BaseCell %d, coverChildren=%v: bbox mismatch\n"+
							"  Go:  N=%.15f, S=%.15f, E=%.15f, W=%.15f\n"+
							"  C:   N=%.15f, S=%.15f, E=%.15f, W=%.15f",
							baseCell, coverChildren,
							float64(goBBox.North), float64(goBBox.South), float64(goBBox.East), float64(goBBox.West),
							float64(cBBox.North), float64(cBBox.South), float64(cBBox.East), float64(cBBox.West))
					}
				}
			}
		})
	}
}

func Test_cellToBBox_pole_cells_parity(t *testing.T) {
	// Test pole cells at various resolutions
	for res := int32(0); res <= 5; res++ { // Test first few resolutions
		t.Run("Poles_res_"+string(rune(res+'0')), func(t *testing.T) {
			// Test north pole cell
			northCell := northPoleCells[res]
			goBBoxN, goErrN := cellToBBox(northCell, false)
			cBBoxN, cErrN := cellToBBoxC(northCell, false)

			if goErrN != cErrN {
				t.Errorf("North pole res %d: Error mismatch: Go=%v, C=%v", res, goErrN, cErrN)
			} else if goErrN == eSuccess {
				// For north pole cells, the north boundary should be π/2
				if goBBoxN.North != PiOver2 {
					t.Errorf("North pole res %d: Go north boundary should be π/2, got %.15f", res, float64(goBBoxN.North))
				}
				if cBBoxN.North != PiOver2 {
					t.Errorf("North pole res %d: C north boundary should be π/2, got %.15f", res, float64(cBBoxN.North))
				}
				// East/West should span full longitude
				if goBBoxN.East != Pi || goBBoxN.West != -Pi {
					t.Errorf("North pole res %d: Go should span full longitude, got E=%.15f, W=%.15f",
						res, float64(goBBoxN.East), float64(goBBoxN.West))
				}
			}

			// Test south pole cell
			southCell := southPoleCells[res]
			goBBoxS, goErrS := cellToBBox(southCell, false)
			cBBoxS, cErrS := cellToBBoxC(southCell, false)

			if goErrS != cErrS {
				t.Errorf("South pole res %d: Error mismatch: Go=%v, C=%v", res, goErrS, cErrS)
			} else if goErrS == eSuccess {
				// For south pole cells, the south boundary should be -π/2
				if goBBoxS.South != -PiOver2 {
					t.Errorf("South pole res %d: Go south boundary should be -π/2, got %.15f", res, float64(goBBoxS.South))
				}
				if cBBoxS.South != -PiOver2 {
					t.Errorf("South pole res %d: C south boundary should be -π/2, got %.15f", res, float64(cBBoxS.South))
				}
				// East/West should span full longitude
				if goBBoxS.East != Pi || goBBoxS.West != -Pi {
					t.Errorf("South pole res %d: Go should span full longitude, got E=%.15f, W=%.15f",
						res, float64(goBBoxS.East), float64(goBBoxS.West))
				}
			}
		})
	}
}
