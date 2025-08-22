package c2go

const (
	_maxH3Res                = 15
	_perDigitBitWidth        = 3
	_digitMask        uint64 = 7
)

// getIndexDigit returns the direction digit at res (port of H3_GET_INDEX_DIGIT).
func getIndexDigit(h H3Index, res int) int {
	shift := (_maxH3Res - res) * _perDigitBitWidth
	return int((uint64(h) >> shift) & _digitMask)
}

// setIndexDigit sets the direction digit at res (port of H3_SET_INDEX_DIGIT).
func setIndexDigit(h H3Index, res int, digit int) H3Index {
	shift := (_maxH3Res - res) * _perDigitBitWidth
	mask := _digitMask << shift
	x := uint64(h)
	x &^= mask
	x |= (uint64(digit) & _digitMask) << shift
	return H3Index(x)
}
