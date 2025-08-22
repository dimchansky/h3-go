package c2go

// _unitIjkToDigit converts a unit IJK coordinate to a direction digit.
// Mirrors H3's coordijk.c::_unitIjkToDigit behavior.
func _unitIjkToDigit(ijk *CoordIJK) int {
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
