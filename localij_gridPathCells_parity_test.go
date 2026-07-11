//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_gridPathCells_parity(t *testing.T) {
	tests := []struct {
		name      string
		start     h3Index
		end       h3Index
		expectErr h3Error
	}{
		{
			name:      "same cell",
			start:     0x8a1fb46622dffff,
			end:       0x8a1fb46622dffff,
			expectErr: eSuccess,
		},
		{
			name:      "adjacent cells",
			start:     0x8a1fb46622dffff,
			end:       0x8a1fb46622d7fff,
			expectErr: eSuccess,
		},
		{
			name:      "adjacent cells res 5",
			start:     0x85283473fffffff,
			end:       0x85283477fffffff,
			expectErr: eSuccess,
		},
		{
			name:      "distant cells res 7",
			start:     0x87283470fffffff,
			end:       0x87283471fffffff,
			expectErr: eSuccess,
		},
		{
			name:      "resolution 0 cells",
			start:     0x8001fffffffffff,
			end:       0x8007fffffffffff,
			expectErr: eSuccess,
		},
		{
			name:      "resolution 1 cells",
			start:     0x81083ffffffffff,
			end:       0x81093ffffffffff,
			expectErr: eSuccess,
		},
		{
			name:      "different resolutions",
			start:     0x8a1fb46622dffff,
			end:       0x891fb46622dffff,
			expectErr: eResMismatch,
		},
		{
			name:      "invalid start cell",
			start:     0x0,
			end:       0x85283473fffffff,
			expectErr: eResMismatch, // Invalid cell causes resolution mismatch error
		},
		{
			name:      "invalid end cell",
			start:     0x85283473fffffff,
			end:       0x0,
			expectErr: eResMismatch, // Invalid cell causes resolution mismatch error
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

			if tt.expectErr != eSuccess {
				// Expected error case, no need to test further
				return
			}

			// Compare sizes
			if expectedSize != goSize {
				t.Errorf("gridPathCellsSize mismatch: C got %d, Go got %d", expectedSize, goSize)
				return
			}

			// Now test the actual gridPathCells function
			cOut := make([]h3Index, expectedSize)
			cErr := _gridPathCellsC(tt.start, tt.end, cOut)

			goOut := make([]h3Index, expectedSize)
			goErr := gridPathCells(goOut, tt.start, tt.end)

			// Check error parity
			if cErr != goErr {
				t.Errorf("gridPathCells error mismatch: C got %v, Go got %v", cErr, goErr)
				return
			}

			if cErr != tt.expectErr {
				t.Errorf("Expected error %v, got %v", tt.expectErr, cErr)
				return
			}

			if tt.expectErr != eSuccess {
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
			dst := make([]h3Index, expectedSize)
			goErrDst := gridPathCells(dst, tt.start, tt.end)

			if goErrDst != goErr {
				t.Errorf("Error mismatch with dst buffer: got %v, expected %v", goErrDst, goErr)
			}

			if len(dst) != len(goOut) {
				t.Errorf("Length mismatch with dst buffer: got %d, expected %d", len(dst), len(goOut))
			}

			for i := 0; i < len(goOut); i++ {
				if dst[i] != goOut[i] {
					t.Errorf("Cell at index %d mismatch with dst buffer: got %016x, expected %016x", i, dst[i], goOut[i])
				}
			}
		})
	}
}

func Test_gridPathCells_pentagon_distortion(t *testing.T) {
	// Test cases that may involve pentagon distortion
	tests := []struct {
		name  string
		start h3Index
		end   h3Index
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

			if cSizeErr != eSuccess {
				// If size calculation fails, gridPathCells should also fail
				goOut := make([]h3Index, 1) // Minimal buffer to test failure
				goErr := gridPathCells(goOut, tt.start, tt.end)
				if goErr == eSuccess {
					t.Errorf("Expected gridPathCells to fail when gridPathCellsSize failed, but got success")
				}
				return
			}

			// Test gridPathCells
			cOut := make([]h3Index, cSize)
			cErr := _gridPathCellsC(tt.start, tt.end, cOut)

			goOut := make([]h3Index, cSize)
			goErr := gridPathCells(goOut, tt.start, tt.end)

			// Check error parity
			if cErr != goErr {
				t.Errorf("gridPathCells error mismatch: C got %v, Go got %v", cErr, goErr)
				return
			}

			if cErr == eSuccess {
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
