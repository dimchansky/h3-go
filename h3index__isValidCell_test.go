// Tests ported from testH3Index.c
package h3

import (
	"fmt"
	"testing"
)

func TestLatLngToCellExtremeCoordinates(t *testing.T) {
	t.Parallel()
	// Check that none of these cause crashes.
	tests := []struct {
		lat float64
		lng float64
		res int32
	}{
		{0, 1e45, 14},
		{1e46, 1e45, 15},
		{degsToRads(2), degsToRads(-3e39), 0},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("lat_%v_lng_%v_res_%d", tc.lat, tc.lng, tc.res), func(t *testing.T) {
			g := LatLng{Lat: Rad(tc.lat), Lng: Rad(tc.lng)}
			var h H3Index
			err := latLngToCell(&g, tc.res, &h)
			if err != E_SUCCESS {
				t.Errorf("latLngToCell failed with extreme coordinates: %v", err)
			}
		})
	}
}

func TestIsValidCellAtResolution(t *testing.T) {
	t.Parallel()
	for i := int32(0); i <= MAX_H3_RES; i++ {
		t.Run(fmt.Sprintf("res_%d", i), func(t *testing.T) {
			g := LatLng{Lat: Rad(0), Lng: Rad(0)}
			var h3 H3Index
			err := latLngToCell(&g, i, &h3)
			if err != E_SUCCESS {
				t.Fatalf("latLngToCell failed: %v", err)
			}
			if !isValidCell(h3) {
				t.Errorf("isValidCell failed on resolution %d", i)
			}
		})
	}
}

func TestIsValidCellDigits(t *testing.T) {
	t.Parallel()
	g := LatLng{Lat: Rad(0), Lng: Rad(0)}
	var h3 H3Index
	err := latLngToCell(&g, 1, &h3)
	if err != E_SUCCESS {
		t.Fatalf("latLngToCell failed: %v", err)
	}
	// Set a bit for an unused digit to something else.
	h3 ^= 1
	if isValidCell(h3) {
		t.Error("isValidCell failed on invalid unused digits")
	}
}

func TestIsValidCellBaseCell(t *testing.T) {
	t.Parallel()
	for i := int32(0); i < NUM_BASE_CELLS; i++ {
		t.Run(fmt.Sprintf("base_cell_%d", i), func(t *testing.T) {
			h := H3Index(H3_INIT)
			h = setMode(h, H3_CELL_MODE)
			h = setBaseCell(h, i)

			if !isValidCell(h) {
				t.Errorf("isValidCell failed on base cell %d", i)
			}

			if getBaseCellNumber(h) != i {
				t.Errorf("failed to recover base cell: got %d, want %d", getBaseCellNumber(h), i)
			}
		})
	}
}

func TestIsValidCellBaseCellInvalid(t *testing.T) {
	t.Parallel()
	hWrongBaseCell := H3Index(0)
	hWrongBaseCell = setMode(hWrongBaseCell, H3_CELL_MODE)
	hWrongBaseCell = setBaseCell(hWrongBaseCell, NUM_BASE_CELLS)
	if isValidCell(hWrongBaseCell) {
		t.Error("isValidCell failed on invalid base cell")
	}
}

func TestIsValidCellWithMode(t *testing.T) {
	t.Parallel()
	for i := int32(0); i <= 0xf; i++ {
		t.Run(fmt.Sprintf("mode_%d", i), func(t *testing.T) {
			h := H3Index(H3_INIT)
			h = setMode(h, i)
			if i == H3_CELL_MODE {
				if !isValidCell(h) {
					t.Error("isValidCell should succeed on valid mode")
				}
			} else {
				if isValidCell(h) {
					t.Errorf("isValidCell failed on mode %d", i)
				}
			}
		})
	}
}

func TestIsValidCellReservedBits(t *testing.T) {
	t.Parallel()
	for i := int32(0); i < 8; i++ {
		t.Run(fmt.Sprintf("reserved_bits_%d", i), func(t *testing.T) {
			h := H3Index(H3_INIT)
			h = setMode(h, H3_CELL_MODE)
			h = setReservedBits(h, i)
			if i == 0 {
				if !isValidCell(h) {
					t.Error("isValidCell should succeed on valid reserved bits")
				}
			} else {
				if isValidCell(h) {
					t.Errorf("isValidCell failed on reserved bits %d", i)
				}
			}
		})
	}
}

func TestIsValidCellHighBit(t *testing.T) {
	t.Parallel()
	h := H3Index(H3_INIT)
	h = setMode(h, H3_CELL_MODE)
	h = setHighBit(h, 1)
	if isValidCell(h) {
		t.Error("isValidCell failed on high bit")
	}
}

func TestH3BadDigitInvalid(t *testing.T) {
	t.Parallel()
	h := H3Index(0)
	// By default the first index digit is out of range.
	h = setMode(h, H3_CELL_MODE)
	h = setResolution(h, 1)
	if isValidCell(h) {
		t.Error("isValidCell failed on too large digit")
	}
}

func TestH3DeletedSubsequenceInvalid(t *testing.T) {
	t.Parallel()
	// Create an index located in a deleted subsequence of a pentagon.
	var h H3Index
	setH3Index(&h, 1, 4, int32(K_AXES_DIGIT))
	if isValidCell(h) {
		t.Error("isValidCell failed on deleted subsequence")
	}
}

func TestMoreDeletedSubsequenceInvalid(t *testing.T) {
	t.Parallel()
	p := H3Index(0x80c3fffffffffff) // res 0 pentagon

	for res := int32(1); res <= 15; res++ {
		t.Run(fmt.Sprintf("res_%d", res), func(t *testing.T) {
			h, err := cellToCenterChild(p, res)
			if err != E_SUCCESS {
				t.Fatalf("cellToCenterChild failed: %v", err)
			}
			if !isValidCell(h) {
				t.Error("should be a valid pentagon")
			}

			for d := int32(0); d <= 6; d++ {
				hTest := setIndexDigit(h, res, d)
				if d == 1 {
					if isValidCell(hTest) {
						t.Error("fail on deleted subsequence")
					}
				} else {
					if !isValidCell(hTest) {
						t.Errorf("should be valid for digit %d", d)
					}
				}
			}
		})
	}
}

func TestH3ToString(t *testing.T) {
	t.Parallel()
	// Note: The Go implementation of h3ToString returns a string directly,
	// not using a buffer like the C version. Testing the actual behavior.

	tests := []struct {
		name     string
		h        H3Index
		expected string
	}{
		{
			name:     "base16_cafe",
			h:        0xcafe,
			expected: "cafe",
		},
		{
			name:     "large_input",
			h:        0xffffffffffffffff,
			expected: "ffffffffffffffff",
		},
		{
			name:     "small_value",
			h:        0x1234,
			expected: "1234",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, errCode := h3ToString(tc.h)

			if errCode != 0 {
				t.Errorf("h3ToString returned error code %d", errCode)
				return
			}

			if result != tc.expected {
				t.Errorf("h3ToString produced %q, expected %q", result, tc.expected)
			}
		})
	}
}

func TestStringToH3(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected H3Index
		wantErr  bool
	}{
		{
			name:    "empty_string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "junk_input",
			input:   "**",
			wantErr: true,
		},
		{
			name:     "large_input",
			input:    "ffffffffffffffff",
			expected: 0xffffffffffffffff,
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h3, err := stringToH3(tc.input)

			if tc.wantErr {
				if err == E_SUCCESS {
					t.Error("expected error but got success")
				}
				return
			}

			if err != E_SUCCESS {
				t.Errorf("stringToH3 failed with error %v", err)
				return
			}

			if h3 != tc.expected {
				t.Errorf("got %#x, expected %#x", h3, tc.expected)
			}
		})
	}
}

func TestSetH3Index(t *testing.T) {
	t.Parallel()
	var h H3Index
	setH3Index(&h, 5, 12, 1)

	if getResolution(h) != 5 {
		t.Errorf("resolution: got %d, expected 5", getResolution(h))
	}
	if getBaseCellNumber(h) != 12 {
		t.Errorf("base cell: got %d, expected 12", getBaseCellNumber(h))
	}
	if getMode(h) != H3_CELL_MODE {
		t.Errorf("mode: got %d, expected %d", getMode(h), H3_CELL_MODE)
	}

	for i := int32(1); i <= 5; i++ {
		if getIndexDigit(h, i) != 1 {
			t.Errorf("digit %d: got %d, expected 1", i, getIndexDigit(h, i))
		}
	}

	for i := int32(6); i <= MAX_H3_RES; i++ {
		if getIndexDigit(h, i) != int32(INVALID_DIGIT) {
			t.Errorf("blanked digit %d: got %d, expected %d", i, getIndexDigit(h, i), INVALID_DIGIT)
		}
	}

	if h != 0x85184927fffffff {
		t.Errorf("index: got %#x, expected %#x", h, 0x85184927fffffff)
	}
}

func TestIsResClassIII(t *testing.T) {
	t.Parallel()
	coord := LatLng{Lat: Rad(0), Lng: Rad(0)}
	for i := int32(0); i <= MAX_H3_RES; i++ {
		t.Run(fmt.Sprintf("res_%d", i), func(t *testing.T) {
			var h H3Index
			err := latLngToCell(&coord, i, &h)
			if err != E_SUCCESS {
				t.Fatalf("latLngToCell failed: %v", err)
			}

			got := isResClassIII(h)
			want := isResolutionClassIII(i)

			if got != want {
				t.Errorf("isResClassIII mismatch for res %d: got %v, want %v", i, got, want)
			}
		})
	}
}
