package h3

// _hasGoodTopBits checks if the top 8 bits of an H3Index are valid.
// The expected pattern is 0_0001_000 (bit 59 set, others clear).
// Ported from H3 C: h3Index.c::_hasGoodTopBits.
func _hasGoodTopBits(h H3Index) bool {
	// Extract top 8 bits by right-shifting 56 bits (64 - 8)
	topBits := h >> (64 - 8)
	// Check if it matches the expected pattern: 0b00001000 (decimal 8)
	return topBits == 0b00001000
}
