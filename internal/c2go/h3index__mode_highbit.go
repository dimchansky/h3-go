package c2go

// getMode returns the H3 mode (port of H3_GET_MODE)
func getMode(h H3Index) int {
	const H3_MODE_OFFSET = 59
	const H3_MODE_MASK = uint64(15) << H3_MODE_OFFSET
	return int((uint64(h) & H3_MODE_MASK) >> H3_MODE_OFFSET)
}

// setMode sets the H3 mode (port of H3_SET_MODE)
func setMode(h H3Index, v int) H3Index {
	const H3_MODE_OFFSET = 59
	const H3_MODE_MASK = uint64(15) << H3_MODE_OFFSET
	x := uint64(h)
	x &^= H3_MODE_MASK
	x |= (uint64(v) & 15) << H3_MODE_OFFSET
	return H3Index(x)
}

// getHighBit returns the highest bit (port of H3_GET_HIGH_BIT)
func getHighBit(h H3Index) int {
	const H3_MAX_OFFSET = 63
	const H3_HIGH_BIT_MASK = uint64(1) << H3_MAX_OFFSET
	if (uint64(h)&H3_HIGH_BIT_MASK)>>H3_MAX_OFFSET != 0 {
		return 1
	}
	return 0
}

// setHighBit sets the highest bit (port of H3_SET_HIGH_BIT)
func setHighBit(h H3Index, v int) H3Index {
	const H3_MAX_OFFSET = 63
	const H3_HIGH_BIT_MASK = uint64(1) << H3_MAX_OFFSET
	x := uint64(h)
	x &^= H3_HIGH_BIT_MASK
	x |= (uint64(v) & 1) << H3_MAX_OFFSET
	return H3Index(x)
}
