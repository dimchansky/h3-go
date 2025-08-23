package c2go

// _neighbor applies a direction vector to IJK coordinates.
// Mirrors H3's coordijk.c::_neighbor behavior.
// Ported from H3 C: coordijk.c::neighbor
func _neighbor(ijk *CoordIJK, digit Direction) {
	if digit > CENTER_DIGIT && digit < NUM_DIGITS {
		_ijkAdd(ijk, &UNIT_VECS[digit], ijk)
		_ijkNormalize(ijk)
	}
}
