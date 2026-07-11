package h3

// _downAp7 finds the center of the hex at the next finer aperture 7
// counter-clockwise resolution. Works in place.
// Mirrors H3's coordijk.c::_downAp7 behavior.
// Ported from H3 C: coordijk.c::_downAp7.
func _downAp7(ijk *coordIJK) {
	// res r unit vectors in res r+1
	iVec := coordIJK{3, 0, 1}
	jVec := coordIJK{1, 3, 0}
	kVec := coordIJK{0, 1, 3}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
