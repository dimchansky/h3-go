package h3

// _neighbor applies a direction vector to IJK coordinates.
// Mirrors H3's coordijk.c::_neighbor behavior.
// Ported from H3 C: coordijk.c::_neighbor.
func _neighbor(ijk *coordIJK, digit direction) {
	if digit > centerDigit && digit < numDigits {
		_ijkAdd(ijk, &unitVecs[digit], ijk)
		_ijkNormalize(ijk)
	}
}
