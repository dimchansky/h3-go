package h3

// _unitIjkToDigit converts a unit IJK coordinate to a direction digit.
// Mirrors H3's coordijk.c::_unitIjkToDigit behavior.
// Ported from H3 C: coordijk.c::_unitIjkToDigit.
func _unitIjkToDigit(ijk *CoordIJK) Direction {
	c := *ijk
	_ijkNormalize(&c)
	digit := INVALID_DIGIT
	for i := CENTER_DIGIT; i < NUM_DIGITS; i++ {
		if _ijkMatches(&c, &UNIT_VECS[i]) {
			digit = i
			break
		}
	}
	return digit
}
