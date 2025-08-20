package h3

import (
	"testing"
	
	"github.com/dimchansky/h3-go/internal/indexbits"
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
			want:     0x080a000000000000, // mode=1, res=0, baseCell=10, all digits=7
		},
		{
			name:     "resolution 1 cell",
			mode:     1,
			res:      1,
			baseCell: 20,
			digits:   []int{3},
			want:     0x0814dfffffffffff, // mode=1, res=1, baseCell=20, digit[0]=3, rest=7
		},
		{
			name:     "resolution 5 cell",
			mode:     1,
			res:      5,
			baseCell: 42,
			digits:   []int{0, 1, 2, 3, 4},
			want:     0x085548e7ffffffff, // mode=1, res=5, baseCell=42, digits[0-4]=0,1,2,3,4, rest=7
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := indexbits.Pack(tt.mode, tt.res, tt.baseCell, tt.digits)
			if got != tt.want {
				t.Errorf("Pack() = 0x%016x, want 0x%016x", got, tt.want)
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
			h:            0x080a000000000000,
			wantMode:     1,
			wantRes:      0,
			wantBaseCell: 10,
			wantDigits:   []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
		},
		{
			name:         "resolution 1 cell",
			h:            0x0814dfffffffffff,
			wantMode:     1,
			wantRes:      1,
			wantBaseCell: 20,
			wantDigits:   []int{3, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
		},
		{
			name:         "resolution 5 cell",
			h:            0x085548e7ffffffff,
			wantMode:     1,
			wantRes:      5,
			wantBaseCell: 42,
			wantDigits:   []int{0, 1, 2, 3, 4, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, res, baseCell, digits := indexbits.Unpack(tt.h)
			
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
}

func TestIsValidCell(t *testing.T) {
	tests := []struct {
		name string
		h    uint64
		want bool
	}{
		{
			name: "valid resolution 0 cell",
			h:    0x080a000000000000,
			want: true,
		},
		{
			name: "valid resolution 5 cell",
			h:    0x085548e7ffffffff,
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
			h:    0x090a000000000000, // reserved bit set
			want: false,
		},
		{
			name: "resolution too high",
			h:    0x08fa000000000000, // res=16 (>15)
			want: false,
		},
		{
			name: "base cell too high",
			h:    0x08ff000000000000, // baseCell=127 (>121)
			want: false,
		},
		{
			name: "digit after resolution not invalid",
			h:    0x081400ffffffffff, // res=1, digit[1]=0 (should be 7)
			want: false,
		},
		{
			name: "digit before resolution invalid",
			h:    0x0854ffffffffff7f, // res=5, digit[4]=invalid (last digit before padding)
			want: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := indexbits.IsValidCell(tt.h)
			if got != tt.want {
				t.Errorf("IsValidCell(0x%016x) = %v, want %v", tt.h, got, tt.want)
			}
		})
	}
}

func TestGetSetMode(t *testing.T) {
	h := uint64(0)
	
	// Set mode to 1 (hexagon)
	h = indexbits.SetMode(h, 1)
	if got := indexbits.GetMode(h); got != 1 {
		t.Errorf("GetMode after SetMode(1) = %d, want 1", got)
	}
	
	// Set mode to 7 (max value for 3 bits)
	h = indexbits.SetMode(h, 7)
	if got := indexbits.GetMode(h); got != 7 {
		t.Errorf("GetMode after SetMode(7) = %d, want 7", got)
	}
}

func TestGetSetResolution(t *testing.T) {
	h := uint64(0)
	
	for res := 0; res <= 15; res++ {
		h = indexbits.SetResolution(h, res)
		if got := indexbits.GetResolution(h); got != res {
			t.Errorf("GetResolution after SetResolution(%d) = %d, want %d", res, got, res)
		}
	}
}

func TestGetSetBaseCell(t *testing.T) {
	h := uint64(0)
	
	testCases := []int{0, 1, 42, 121, 127}
	for _, bc := range testCases {
		h = indexbits.SetBaseCell(h, bc)
		if got := indexbits.GetBaseCell(h); got != bc {
			t.Errorf("GetBaseCell after SetBaseCell(%d) = %d, want %d", bc, got, bc)
		}
	}
}

func TestGetSetDigit(t *testing.T) {
	h := uint64(0)
	
	// Set different digits at different positions
	for pos := 0; pos < 15; pos++ {
		for digit := 0; digit <= 7; digit++ {
			h = indexbits.SetDigit(h, pos, digit)
			if got := indexbits.GetDigit(h, pos); got != digit {
				t.Errorf("GetDigit(%d) after SetDigit(%d, %d) = %d, want %d", 
					pos, pos, digit, got, digit)
			}
		}
	}
	
	// Test out of bounds
	if got := indexbits.GetDigit(h, -1); got != indexbits.InvalidDigit {
		t.Errorf("GetDigit(-1) = %d, want %d", got, indexbits.InvalidDigit)
	}
	if got := indexbits.GetDigit(h, 15); got != indexbits.InvalidDigit {
		t.Errorf("GetDigit(15) = %d, want %d", got, indexbits.InvalidDigit)
	}
}

func BenchmarkIndexPack(b *testing.B) {
	digits := []int{0, 1, 2, 3, 4, 5, 6}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = indexbits.Pack(1, 7, 42, digits)
	}
}

func BenchmarkIndexUnpack(b *testing.B) {
	h := uint64(0x08752ad184927fff)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = indexbits.Unpack(h)
	}
}