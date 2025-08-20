// Package indexbits provides utilities for packing and unpacking H3 64-bit indices.
// H3 indices encode resolution, base cell, and up to 15 directional digits in a uint64.
package indexbits

// H3 index bit layout constants
const (
	// Mode bits (3 bits at positions 59-61)
	ModeBitOffset = 59
	ModeBitMask   = uint64(0x7) << ModeBitOffset
	ModeHexagon   = uint64(1) << ModeBitOffset

	// Reserved bits (4 bits at positions 55-58)
	ReservedBitOffset = 55
	ReservedBitMask   = uint64(0xF) << ReservedBitOffset

	// Resolution bits (4 bits at positions 51-54)
	ResolutionBitOffset = 51
	ResolutionBitMask   = uint64(0xF) << ResolutionBitOffset

	// Base cell bits (7 bits at positions 44-50)
	BaseCellBitOffset = 44
	BaseCellBitMask   = uint64(0x7F) << BaseCellBitOffset

	// 15 directional digits (3 bits each from positions 0-44)
	DirectionBitOffset = 0
	DirectionBitMask   = uint64(0x7)
	NumDigits          = 15
	DigitBits          = 3

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

// GetDigit extracts a directional digit at a given position (0-14).
func GetDigit(h uint64, position int) int {
	if position < 0 || position >= NumDigits {
		return InvalidDigit
	}
	offset := (NumDigits - 1 - position) * DigitBits
	return int((h >> offset) & DirectionBitMask)
}

// SetDigit sets a directional digit at a given position (0-14).
func SetDigit(h uint64, position int, digit int) uint64 {
	if position < 0 || position >= NumDigits {
		return h
	}
	offset := (NumDigits - 1 - position) * DigitBits
	mask := DirectionBitMask << offset
	return (h &^ mask) | (uint64(digit&0x7) << offset)
}

// Pack creates an H3 index from components.
func Pack(mode uint64, res int, baseCell int, digits []int) uint64 {
	h := uint64(0)
	h = SetMode(h, mode)
	h = SetResolution(h, res)
	h = SetBaseCell(h, baseCell)
	
	// Set directional digits
	for i := 0; i < len(digits) && i < NumDigits; i++ {
		h = SetDigit(h, i, digits[i])
	}
	
	// Fill remaining with invalid digits
	for i := len(digits); i < NumDigits; i++ {
		h = SetDigit(h, i, InvalidDigit)
	}
	
	return h
}

// Unpack extracts all components from an H3 index.
func Unpack(h uint64) (mode uint64, res int, baseCell int, digits []int) {
	mode = GetMode(h)
	res = GetResolution(h)
	baseCell = GetBaseCell(h)
	
	digits = make([]int, NumDigits)
	for i := 0; i < NumDigits; i++ {
		digits[i] = GetDigit(h, i)
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
	
	// Check that digits after resolution are all invalid (7)
	for i := res; i < NumDigits; i++ {
		if GetDigit(h, i) != InvalidDigit {
			return false
		}
	}
	
	// Check that digits before resolution are valid (0-6)
	for i := 0; i < res; i++ {
		digit := GetDigit(h, i)
		if digit < 0 || digit > 6 {
			return false
		}
	}
	
	return true
}