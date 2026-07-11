package h3

// getReservedBits returns the 3 reserved bits (port of H3_GET_RESERVED_BITS).
// Ported from H3 C: h3Index.h::H3_GET_RESERVED_BITS.
func getReservedBits(h h3Index) int32 {
	const h3ReservedOffset = 56
	const h3ReservedMask = uint64(7) << h3ReservedOffset
	return int32((uint64(h) & h3ReservedMask) >> h3ReservedOffset)
}

// setReservedBits sets the 3 reserved bits (port of H3_SET_RESERVED_BITS).
// Note: setting to non-zero may produce invalid indexes.
// Ported from H3 C: h3Index.h::H3_SET_RESERVED_BITS.
func setReservedBits(h h3Index, v int32) h3Index {
	const h3ReservedOffset = 56
	const h3ReservedMask = uint64(7) << h3ReservedOffset
	x := uint64(h)
	x &^= h3ReservedMask
	x |= (uint64(v) & 7) << h3ReservedOffset
	return h3Index(x)
}
