package c2go

// _getResDigit gets the digit at the specified resolution position.
// This function extracts a specific resolution digit from an H3 index.
// Ported from H3 C: h3Index.c::_getResDigit
func _getResDigit(h H3Index, res int) Direction {
	if res < 1 || res > MAX_H3_RES {
		return INVALID_DIGIT
	}
	return Direction(getIndexDigit(h, res))
}
