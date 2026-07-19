package h3

// _unitIjkToDigit converts a unit IJK coordinate to a direction digit.
// Mirrors H3's coordijk.h::_unitIjkToDigit behavior.
// Ported from H3 C: coordijk.h::_unitIjkToDigit.
func _unitIjkToDigit(ijk *coordIJK) direction {
	c := *ijk
	_ijkNormalize(&c)
	digit := invalidDigit
	for i := centerDigit; i < numDigits; i++ {
		if _ijkMatches(&c, &unitVecs[i]) {
			digit = i
			break
		}
	}
	return digit
}
