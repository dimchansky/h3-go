//go:build cgo

package h3

import (
	"testing"
)

func Test_gridPathCells_parity(t *testing.T) {
	tests := []struct {
		name      string
		start     H3Index
		end       H3Index
		expectErr H3Error
	}{
		{
			name:      "same cell",
			start:     0x8a1fb46622dffff,
			end:       0x8a1fb46622dffff,
			expectErr: E_SUCCESS,
		},
		{
			name:      "adjacent cells",
			start:     0x8a1fb46622dffff,
			end:       0x8a1fb46622d7fff,
			expectErr: E_SUCCESS,
		},
		{
			name:      "adjacent cells res 5",
			start:     0x85283473fffffff,
			end:       0x85283477fffffff,
			expectErr: E_SUCCESS,
		},
		{
			name:      "distant cells res 7",
			start:     0x87283470fffffff,
			end:       0x87283471fffffff,
			expectErr: E_SUCCESS,
		},
		{
			name:      "resolution 0 cells",
			start:     0x8001fffffffffff,
			end:       0x8007fffffffffff,
			expectErr: E_SUCCESS,
		},
		{
			name:      "resolution 1 cells",
			start:     0x81083ffffffffff,
			end:       0x81093ffffffffff,
			expectErr: E_SUCCESS,
		},
		{
			name:      "different resolutions",
			start:     0x8a1fb46622dffff,
			end:       0x891fb46622dffff,
			expectErr: E_RES_MISMATCH,
		},
		{
			name:      "invalid start cell",
			start:     0x0,
			end:       0x85283473fffffff,
			expectErr: E_RES_MISMATCH, // Invalid cell causes resolution mismatch error
		},
		{
			name:      "invalid end cell",
			start:     0x85283473fffffff,
			end:       0x0,
			expectErr: E_RES_MISMATCH, // Invalid cell causes resolution mismatch error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First, get the expected size to allocate the buffer
			var expectedSize int64
			cSizeErr := _gridPathCellsSizeC(tt.start, tt.end, &expectedSize)

			var goSize int64
			goSizeErr := gridPathCellsSize(tt.start, tt.end, &goSize)

			// Check size calculation parity
			if cSizeErr != goSizeErr {
				t.Errorf("gridPathCellsSize error mismatch: C got %v, Go got %v", cSizeErr, goSizeErr)
				return
			}

			if cSizeErr != tt.expectErr {
				t.Errorf("Expected error %v, got %v", tt.expectErr, cSizeErr)
				return
			}

			if tt.expectErr != E_SUCCESS {
				// Expected error case, no need to test further
				return
			}

			// Compare sizes
			if expectedSize != goSize {
				t.Errorf("gridPathCellsSize mismatch: C got %d, Go got %d", expectedSize, goSize)
				return
			}

			// Now test the actual gridPathCells function
			cOut := make([]H3Index, expectedSize)
			cErr := _gridPathCellsC(tt.start, tt.end, cOut)

			goOut, goErr := gridPathCells(nil, tt.start, tt.end)

			// Check error parity
			if cErr != goErr {
				t.Errorf("gridPathCells error mismatch: C got %v, Go got %v", cErr, goErr)
				return
			}

			if cErr != tt.expectErr {
				t.Errorf("Expected error %v, got %v", tt.expectErr, cErr)
				return
			}

			if tt.expectErr != E_SUCCESS {
				return
			}

			// Check output length
			if int64(len(goOut)) != expectedSize {
				t.Errorf("Go output length mismatch: expected %d, got %d", expectedSize, len(goOut))
				return
			}

			// Check each cell in the path
			for i := 0; i < len(cOut); i++ {
				if cOut[i] != goOut[i] {
					t.Errorf("Cell at index %d mismatch: C got %016x, Go got %016x", i, cOut[i], goOut[i])
				}
			}

			// Note: The H3 documentation states that gridPathCells output is not guaranteed to
			// start/end with the exact input cells. The algorithm draws lines in grid space
			// which may not correspond exactly to the input coordinates.
			// We only verify parity between C and Go implementations.

			// Test with pre-allocated buffer (dst-buffer pattern)
			dst := make([]H3Index, expectedSize)
			goOutDst, goErrDst := gridPathCells(dst, tt.start, tt.end)

			if goErrDst != goErr {
				t.Errorf("Error mismatch with dst buffer: got %v, expected %v", goErrDst, goErr)
			}

			if len(goOutDst) != len(goOut) {
				t.Errorf("Length mismatch with dst buffer: got %d, expected %d", len(goOutDst), len(goOut))
			}

			for i := 0; i < len(goOut); i++ {
				if goOutDst[i] != goOut[i] {
					t.Errorf("Cell at index %d mismatch with dst buffer: got %016x, expected %016x", i, goOutDst[i], goOut[i])
				}
			}
		})
	}
}

func Test_gridPathCells_pentagon_distortion(t *testing.T) {
	// Test cases that may involve pentagon distortion
	tests := []struct {
		name  string
		start H3Index
		end   H3Index
	}{
		{
			name:  "pentagon cell base",
			start: 0x81283ffffffffff, // Resolution 1 pentagon cell
			end:   0x8128bffffffffff,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cSize int64
			cSizeErr := _gridPathCellsSizeC(tt.start, tt.end, &cSize)

			var goSize int64
			goSizeErr := gridPathCellsSize(tt.start, tt.end, &goSize)

			// Check size calculation parity
			if cSizeErr != goSizeErr {
				t.Errorf("gridPathCellsSize error mismatch: C got %v, Go got %v", cSizeErr, goSizeErr)
				return
			}

			if cSizeErr != E_SUCCESS {
				// If size calculation fails, gridPathCells should also fail
				goOut, goErr := gridPathCells(nil, tt.start, tt.end)
				if goErr == E_SUCCESS {
					t.Errorf("Expected gridPathCells to fail when gridPathCellsSize failed, but got success with %d cells", len(goOut))
				}
				return
			}

			// Test gridPathCells
			cOut := make([]H3Index, cSize)
			cErr := _gridPathCellsC(tt.start, tt.end, cOut)

			goOut, goErr := gridPathCells(nil, tt.start, tt.end)

			// Check error parity
			if cErr != goErr {
				t.Errorf("gridPathCells error mismatch: C got %v, Go got %v", cErr, goErr)
				return
			}

			if cErr == E_SUCCESS {
				// Compare results if successful
				if int64(len(goOut)) != cSize {
					t.Errorf("Output length mismatch: expected %d, got %d", cSize, len(goOut))
					return
				}

				for i := 0; i < len(cOut); i++ {
					if cOut[i] != goOut[i] {
						t.Errorf("Cell at index %d mismatch: C got %016x, Go got %016x", i, cOut[i], goOut[i])
					}
				}
			}
		})
	}
}
