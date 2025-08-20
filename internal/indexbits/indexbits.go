// Package indexbits provides utilities for packing and unpacking H3 64-bit indices.
// H3 indices encode resolution, base cell, and up to 15 directional digits in a uint64.
package indexbits

// H3 index bit layout constants (from H3 C v4.3.0)
const (
	// Mode bits (4 bits at positions 59-62)
	ModeBitOffset = 59
	ModeBitMask   = uint64(0xF) << ModeBitOffset
	ModeHexagon   = uint64(1) << ModeBitOffset

	// Reserved bits (3 bits at positions 56-58)
	ReservedBitOffset = 56
	ReservedBitMask   = uint64(0x7) << ReservedBitOffset

	// Resolution bits (4 bits at positions 52-55)
	ResolutionBitOffset = 52
	ResolutionBitMask   = uint64(0xF) << ResolutionBitOffset

	// Base cell bits (7 bits at positions 45-51)
	BaseCellBitOffset = 45
	BaseCellBitMask   = uint64(0x7F) << BaseCellBitOffset

	// 15 directional digits (3 bits each from positions 0-44)
	DirectionBitMask = uint64(0x7)
	NumDigits        = 15
	DigitBits        = 3
	MaxH3Resolution  = 15

	// Invalid digit marker
	InvalidDigit = 7

	// Maximum values
	MaxResolution = 15
	MaxBaseCell   = 121
)

// GetMode extracts the mode from an H3 index.
func GetMode(h uint64) uint64 {
	return (h & ModeBitMask) >> ModeBitOffset
}

// SetMode sets the mode in an H3 index.
func SetMode(h uint64, mode uint64) uint64 {
	return (h &^ ModeBitMask) | ((mode & 0x7) << ModeBitOffset)
}

// GetResolution extracts the resolution from an H3 index.
func GetResolution(h uint64) int {
	return int((h & ResolutionBitMask) >> ResolutionBitOffset)
}

// SetResolution sets the resolution in an H3 index.
func SetResolution(h uint64, res int) uint64 {
	return (h &^ ResolutionBitMask) | (uint64(res&0xF) << ResolutionBitOffset)
}

// GetBaseCell extracts the base cell from an H3 index.
func GetBaseCell(h uint64) int {
	return int((h & BaseCellBitMask) >> BaseCellBitOffset)
}

// SetBaseCell sets the base cell in an H3 index.
func SetBaseCell(h uint64, baseCell int) uint64 {
	return (h &^ BaseCellBitMask) | (uint64(baseCell&0x7F) << BaseCellBitOffset)
}

// GetDigit extracts a directional digit at a given resolution (1-15).
// This follows H3 C implementation where digits are 1-indexed by resolution.
func GetDigit(h uint64, resolution int) int {
	if resolution < 1 || resolution > MaxH3Resolution {
		return InvalidDigit
	}
	offset := (MaxH3Resolution - resolution) * DigitBits
	return int((h >> offset) & DirectionBitMask)
}

// SetDigit sets a directional digit at a given resolution (1-15).
// This follows H3 C implementation where digits are 1-indexed by resolution.
func SetDigit(h uint64, resolution int, digit int) uint64 {
	if resolution < 1 || resolution > MaxH3Resolution {
		return h
	}
	offset := (MaxH3Resolution - resolution) * DigitBits
	mask := DirectionBitMask << offset
	return (h &^ mask) | (uint64(digit&0x7) << offset)
}

// H3_INIT equivalent - all digits set to 7 (invalid)
const H3_INIT = uint64(0x00001fffffffffff)

// Pack creates an H3 index from components.
func Pack(mode uint64, res int, baseCell int, digits []int) uint64 {
	// Start with H3_INIT (all digits = 7, mode = 0)
	h := H3_INIT
	h = SetMode(h, mode)
	h = SetResolution(h, res)
	h = SetBaseCell(h, baseCell)
	
	// Set directional digits from resolution 1 to res
	for i := 0; i < len(digits) && i < res; i++ {
		h = SetDigit(h, i+1, digits[i])
	}
	
	return h
}

// Unpack extracts all components from an H3 index.
func Unpack(h uint64) (mode uint64, res int, baseCell int, digits []int) {
	mode = GetMode(h)
	res = GetResolution(h)
	baseCell = GetBaseCell(h)
	
	digits = make([]int, MaxH3Resolution)
	for r := 1; r <= MaxH3Resolution; r++ {
		digits[r-1] = GetDigit(h, r)
	}
	
	return
}

// IsValidCell performs basic bit-level validation of an H3 index.
func IsValidCell(h uint64) bool {
	// Check mode is hexagon (1)
	if GetMode(h) != 1 {
		return false
	}
	
	// Check reserved bits are zero
	if (h & ReservedBitMask) != 0 {
		return false
	}
	
	// Check resolution is valid
	res := GetResolution(h)
	if res < 0 || res > MaxResolution {
		return false
	}
	
	// Check base cell is valid
	baseCell := GetBaseCell(h)
	if baseCell < 0 || baseCell > MaxBaseCell {
		return false
	}
	
	// Check that digits beyond resolution are all invalid (7)
	for r := res + 1; r <= MaxH3Resolution; r++ {
		if GetDigit(h, r) != InvalidDigit {
			return false
		}
	}
	
	// Check that digits from resolution 1 to res are valid (0-6)
	for r := 1; r <= res; r++ {
		digit := GetDigit(h, r)
		if digit < 0 || digit > 6 {
			return false
		}
	}
	
	return true
}