package c2go

// getReservedBits returns the 3 reserved bits (port of H3_GET_RESERVED_BITS).
func getReservedBits(h H3Index) int {
	const H3_RESERVED_OFFSET = 56
	const H3_RESERVED_MASK = uint64(7) << H3_RESERVED_OFFSET
	return int((uint64(h) & H3_RESERVED_MASK) >> H3_RESERVED_OFFSET)
}

// setReservedBits sets the 3 reserved bits (port of H3_SET_RESERVED_BITS).
// Note: setting to non-zero may produce invalid indexes.
func setReservedBits(h H3Index, v int) H3Index {
	const H3_RESERVED_OFFSET = 56
	const H3_RESERVED_MASK = uint64(7) << H3_RESERVED_OFFSET
	x := uint64(h)
	x &^= H3_RESERVED_MASK
	x |= (uint64(v) & 7) << H3_RESERVED_OFFSET
	return H3Index(x)
}
