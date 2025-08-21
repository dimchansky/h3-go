package indexbits

import (
	"testing"
)

func TestIndexPack(t *testing.T) {
	tests := []struct {
		name     string
		mode     uint64
		res      int
		baseCell int
		digits   []int
		want     uint64
	}{
		{
			name:     "resolution 0 cell",
			mode:     1,
			res:      0,
			baseCell: 10,
			digits:   []int{},
			want:     0x08015fffffffffff, // mode=1, res=0, baseCell=10, all digits=7
		},
		{
			name:     "resolution 1 cell",
			mode:     1,
			res:      1,
			baseCell: 20,
			digits:   []int{3},
			want:     0x08128fffffffffff, // mode=1, res=1, baseCell=20, digit[0]=3, rest=7
		},
		{
			name:     "resolution 5 cell",
			mode:     1,
			res:      5,
			baseCell: 42,
			digits:   []int{0, 1, 2, 3, 4},
			want:     0x085540a73fffffff, // mode=1, res=5, baseCell=42, digits[0-4]=0,1,2,3,4, rest=7
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pack(tt.mode, tt.res, tt.baseCell, tt.digits)
			if got != tt.want {
				t.Errorf("Pack(%d, %d, %d, %v) = %016x, want %016x",
					tt.mode, tt.res, tt.baseCell, tt.digits, got, tt.want)
			}
		})
	}
}

func TestIndexUnpack(t *testing.T) {
	tests := []struct {
		name         string
		h            uint64
		wantMode     uint64
		wantRes      int
		wantBaseCell int
		wantDigits   []int
	}{
		{
			name:         "resolution 0 cell",
			h:            0x08015fffffffffff,
			wantMode:     1,
			wantRes:      0,
			wantBaseCell: 10,
			wantDigits:   []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
		},
		{
			name:         "resolution 1 cell",
			h:            0x08128fffffffffff,
			wantMode:     1,
			wantRes:      1,
			wantBaseCell: 20,
			wantDigits:   []int{3, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
		},
		{
			name:         "resolution 5 cell",
			h:            0x085540a73fffffff,
			wantMode:     1,
			wantRes:      5,
			wantBaseCell: 42,
			wantDigits:   []int{0, 1, 2, 3, 4, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, res, baseCell, digits := Unpack(tt.h, nil)

			if mode != tt.wantMode {
				t.Errorf("Unpack() mode = %d, want %d", mode, tt.wantMode)
			}
			if res != tt.wantRes {
				t.Errorf("Unpack() res = %d, want %d", res, tt.wantRes)
			}
			if baseCell != tt.wantBaseCell {
				t.Errorf("Unpack() baseCell = %d, want %d", baseCell, tt.wantBaseCell)
			}

			if len(digits) != len(tt.wantDigits) {
				t.Errorf("Unpack() digits length = %d, want %d", len(digits), len(tt.wantDigits))
			}
			for i := range digits {
				if digits[i] != tt.wantDigits[i] {
					t.Errorf("Unpack() digits[%d] = %d, want %d", i, digits[i], tt.wantDigits[i])
				}
			}
		})
	}

	// Test buffer reuse
	t.Run("buffer reuse", func(t *testing.T) {
		h := uint64(0x085540a73fffffff)

		// Test with nil buffer (should allocate)
		_, _, _, digits1 := Unpack(h, nil)
		if len(digits1) != MaxH3Resolution {
			t.Errorf("Expected %d digits, got %d", MaxH3Resolution, len(digits1))
		}

		// Test with sufficient capacity buffer (should reuse)
		buffer := make([]int, 0, MaxH3Resolution)
		_, _, _, digits2 := Unpack(h, buffer)
		if cap(digits2) != cap(buffer) {
			t.Errorf("Expected buffer reuse, got new allocation")
		}
		if len(digits2) != MaxH3Resolution {
			t.Errorf("Expected %d digits, got %d", MaxH3Resolution, len(digits2))
		}

		// Verify digits are the same
		for i := range digits1 {
			if digits1[i] != digits2[i] {
				t.Errorf("Digit mismatch at index %d: %d vs %d", i, digits1[i], digits2[i])
			}
		}
	})
}

func TestIsValidCell(t *testing.T) {
	tests := []struct {
		name string
		h    uint64
		want bool
	}{
		{
			name: "valid resolution 0 cell",
			h:    0x08015fffffffffff,
			want: true,
		},
		{
			name: "valid resolution 5 cell",
			h:    0x085540a73fffffff,
			want: true,
		},
		{
			name: "invalid mode (0)",
			h:    0x000a000000000000,
			want: false,
		},
		{
			name: "invalid mode (2)",
			h:    0x100a000000000000,
			want: false,
		},
		{
			name: "reserved bits set",
			h:    0x090a000000000000,
			want: false,
		},
		{
			name: "base cell too high",
			h:    0x08fe1fffffffffff, // mode=1, res=0, baseCell=127 > 121
			want: false,
		},
		{
			name: "digit after resolution not invalid",
			h:    0x08015ffffffffffe, // res=0 but last digit is 6, not 7
			want: false,
		},
		{
			name: "digit before resolution invalid",
			h:    0x08128f7fffffffff, // res=1, digit[0]=3 (valid), digit[1]=7 (invalid, should be 0-6)
			want: false,
		},
		{
			name: "resolution too high (16)",
			h:    0x08001fffffffffff | (uint64(16) << 52), // Mode=1, res=16 (invalid), baseCell=0
			want: false,
		},
		{
			name: "base cell negative (impossible with uint but testing boundary)",
			h:    Pack(1, 0, 122, nil), // Base cell > 121
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidCell(tt.h)
			if got != tt.want {
				t.Errorf("IsValidCell(%016x) = %v, want %v", tt.h, got, tt.want)
			}
		})
	}
}

func TestGetSetMode(t *testing.T) {
	h := uint64(0)
	// Test all possible mode values (4 bits = 0-15)
	for mode := uint64(0); mode <= 15; mode++ {
		h = SetMode(h, mode)
		if got := GetMode(h); got != mode {
			t.Errorf("GetMode() after SetMode(%d) = %d, want %d", mode, got, mode)
		}
	}
}

func TestGetSetResolution(t *testing.T) {
	h := uint64(0)
	for res := 0; res <= 15; res++ {
		h = SetResolution(h, res)
		if got := GetResolution(h); got != res {
			t.Errorf("GetResolution() after SetResolution(%d) = %d, want %d", res, got, res)
		}
	}
}

func TestGetSetBaseCell(t *testing.T) {
	h := uint64(0)
	for baseCell := 0; baseCell <= 127; baseCell++ {
		h = SetBaseCell(h, baseCell)
		if got := GetBaseCell(h); got != baseCell {
			t.Errorf("GetBaseCell() after SetBaseCell(%d) = %d, want %d", baseCell, got, baseCell)
		}
	}
}

func TestGetSetDigit(t *testing.T) {
	h := uint64(0)

	// Set different digits at different resolutions (1-15)
	for res := 1; res <= 15; res++ {
		for digit := 0; digit <= 7; digit++ {
			h = SetDigit(h, res, digit)
			if got := GetDigit(h, res); got != digit {
				t.Errorf("GetDigit(%d) after SetDigit(%d, %d) = %d, want %d",
					res, res, digit, got, digit)
			}
		}
	}

	// Test out of bounds resolution
	if got := GetDigit(h, 0); got != InvalidDigit {
		t.Errorf("GetDigit(0) = %d, want %d", got, InvalidDigit)
	}
	if got := GetDigit(h, 16); got != InvalidDigit {
		t.Errorf("GetDigit(16) = %d, want %d", got, InvalidDigit)
	}

	// Test SetDigit with out of bounds resolution (should return unchanged)
	originalH := uint64(0x085540a73fffffff)
	unchangedH := SetDigit(originalH, 0, 5) // resolution 0 is invalid
	if unchangedH != originalH {
		t.Errorf("SetDigit with invalid resolution should return unchanged: got %016x, want %016x", unchangedH, originalH)
	}

	unchangedH = SetDigit(originalH, 16, 5) // resolution 16 is invalid
	if unchangedH != originalH {
		t.Errorf("SetDigit with invalid resolution should return unchanged: got %016x, want %016x", unchangedH, originalH)
	}
}

func BenchmarkIndexUnpack(b *testing.B) {
	h := uint64(0x08752ad184927fff)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = Unpack(h, nil)
	}
}

func BenchmarkUnpackWithBuffer(b *testing.B) {
	h := uint64(0x08752ad184927fff)
	buffer := make([]int, 0, MaxH3Resolution)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = Unpack(h, buffer)
	}
}

// Test H3_INIT constant - should have mode=0, all digits=7.
func TestH3InitConstant(t *testing.T) {
	h := H3_INIT

	if GetMode(h) != 0 {
		t.Errorf("H3_INIT mode = %d, want 0", GetMode(h))
	}

	if GetResolution(h) != 0 {
		t.Errorf("H3_INIT resolution = %d, want 0", GetResolution(h))
	}

	if GetBaseCell(h) != 0 {
		t.Errorf("H3_INIT base cell = %d, want 0", GetBaseCell(h))
	}

	// All digits should be 7 (invalid)
	for res := 1; res <= MaxH3Resolution; res++ {
		digit := GetDigit(h, res)
		if digit != InvalidDigit {
			t.Errorf("H3_INIT digit[%d] = %d, want %d", res, digit, InvalidDigit)
		}
	}

	// Verify the actual value matches H3 C constant
	expectedH3Init := uint64(0x00001fffffffffff) // 35184372088831 in decimal
	if h != expectedH3Init {
		t.Errorf("H3_INIT = %016x, want %016x", h, expectedH3Init)
	}
}

// Test setH3Index equivalent - creating index with specific resolution, base cell, and digit.
func TestCreateIndexWithDigit(t *testing.T) {
	tests := []struct {
		name      string
		res       int
		baseCell  int
		initDigit int
		expected  uint64
	}{
		{
			name:      "resolution 5, base cell 12, digit 1",
			res:       5,
			baseCell:  12,
			initDigit: 1,
			expected:  0x85184927fffffff, // From H3 C test
		},
		{
			name:      "resolution 0, base cell 0, digit 0",
			res:       0,
			baseCell:  0,
			initDigit: 0,
			expected:  0x08001fffffffffff,
		},
		{
			name:      "resolution 1, base cell 4, digit 0",
			res:       1,
			baseCell:  4,
			initDigit: 0,
			expected:  0x081083ffffffffff, // Calculated from actual Pack result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create digits array filled with initDigit up to resolution
			digits := make([]int, tt.res)
			for i := 0; i < tt.res; i++ {
				digits[i] = tt.initDigit
			}

			h := Pack(1, tt.res, tt.baseCell, digits)

			if h != tt.expected {
				t.Errorf("Pack(1, %d, %d, %v) = %016x, want %016x",
					tt.res, tt.baseCell, digits, h, tt.expected)
			}

			// Verify components can be extracted correctly
			mode, res, baseCell, extractedDigits := Unpack(h, nil)

			if mode != 1 {
				t.Errorf("Unpacked mode = %d, want 1", mode)
			}
			if res != tt.res {
				t.Errorf("Unpacked resolution = %d, want %d", res, tt.res)
			}
			if baseCell != tt.baseCell {
				t.Errorf("Unpacked base cell = %d, want %d", baseCell, tt.baseCell)
			}

			// Check digits up to resolution
			for i := 0; i < tt.res; i++ {
				if extractedDigits[i] != tt.initDigit {
					t.Errorf("Unpacked digit[%d] = %d, want %d", i, extractedDigits[i], tt.initDigit)
				}
			}

			// Check digits beyond resolution are 7
			for i := tt.res; i < MaxH3Resolution; i++ {
				if extractedDigits[i] != InvalidDigit {
					t.Errorf("Unpacked digit[%d] = %d, want %d (invalid)", i, extractedDigits[i], InvalidDigit)
				}
			}
		})
	}
}

// Test all valid base cells (0-121).
func TestAllValidBaseCells(t *testing.T) {
	for baseCell := 0; baseCell <= MaxBaseCell; baseCell++ {
		h := Pack(1, 0, baseCell, nil)

		if !IsValidCell(h) {
			t.Errorf("IsValidCell failed on base cell %d", baseCell)
		}

		extractedBaseCell := GetBaseCell(h)
		if extractedBaseCell != baseCell {
			t.Errorf("Failed to recover base cell %d, got %d", baseCell, extractedBaseCell)
		}
	}
}

// Test all valid resolutions (0-15).
func TestAllValidResolutions(t *testing.T) {
	for res := 0; res <= MaxH3Resolution; res++ {
		h := Pack(1, res, 0, make([]int, res))

		if !IsValidCell(h) {
			t.Errorf("IsValidCell failed on resolution %d", res)
		}

		extractedRes := GetResolution(h)
		if extractedRes != res {
			t.Errorf("Failed to recover resolution %d, got %d", res, extractedRes)
		}
	}
}
