//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_cellToBBox_parity(t *testing.T) {
	tests := []struct {
		name          string
		cell          H3Index
		coverChildren bool
		expectError   H3Error
	}{
		{
			name:          "Valid res 0 cell",
			cell:          0x8001fffffffffff, // Base cell 1 at res 0
			coverChildren: false,
			expectError:   E_SUCCESS,
		},
		{
			name:          "Valid res 0 cell with children coverage",
			cell:          0x8001fffffffffff,
			coverChildren: true,
			expectError:   E_SUCCESS,
		},
		{
			name:          "Valid res 5 cell",
			cell:          0x851fb46622dffff, // Some res 5 cell
			coverChildren: false,
			expectError:   E_SUCCESS,
		},
		{
			name:          "Valid res 5 cell with children coverage",
			cell:          0x851fb46622dffff,
			coverChildren: true,
			expectError:   E_SUCCESS,
		},
		{
			name:          "North pole cell res 1",
			cell:          NORTH_POLE_CELLS[1],
			coverChildren: false,
			expectError:   E_SUCCESS,
		},
		{
			name:          "South pole cell res 1",
			cell:          SOUTH_POLE_CELLS[1],
			coverChildren: false,
			expectError:   E_SUCCESS,
		},
		{
			name:          "Pentagon base cell",
			cell:          0x804dfffffffffff, // Base cell 4 (pentagon)
			coverChildren: false,
			expectError:   E_SUCCESS,
		},
		{
			name:          "Invalid cell",
			cell:          0x0,
			coverChildren: false,
			expectError:   E_CELL_INVALID,
		},
		{
			name:          "Another invalid cell",
			cell:          H3_NULL,
			coverChildren: false,
			expectError:   E_CELL_INVALID,
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
			if goErr == E_SUCCESS && cErr == E_SUCCESS {
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
	for baseCell := int32(0); baseCell < NUM_BASE_CELLS; baseCell++ {
		t.Run("BaseCell_"+string(rune(baseCell+'0')), func(t *testing.T) {
			cell := baseCellNumToCell(baseCell)
			if cell == H3_NULL {
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
				if goErr == E_SUCCESS && cErr == E_SUCCESS {
					tolerance := 1e-12

					if math.Abs(float64(goBBox.North-cBBox.North)) > tolerance ||
						math.Abs(float64(goBBox.South-cBBox.South)) > tolerance ||
						math.Abs(float64(goBBox.East-cBBox.East)) > tolerance ||
						math.Abs(float64(goBBox.West-cBBox.West)) > tolerance {
						t.Errorf("BaseCell %d, coverChildren=%v: BBox mismatch\n"+
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
			northCell := NORTH_POLE_CELLS[res]
			goBBoxN, goErrN := cellToBBox(northCell, false)
			cBBoxN, cErrN := cellToBBoxC(northCell, false)

			if goErrN != cErrN {
				t.Errorf("North pole res %d: Error mismatch: Go=%v, C=%v", res, goErrN, cErrN)
			} else if goErrN == E_SUCCESS {
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
			southCell := SOUTH_POLE_CELLS[res]
			goBBoxS, goErrS := cellToBBox(southCell, false)
			cBBoxS, cErrS := cellToBBoxC(southCell, false)

			if goErrS != cErrS {
				t.Errorf("South pole res %d: Error mismatch: Go=%v, C=%v", res, goErrS, cErrS)
			} else if goErrS == E_SUCCESS {
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
