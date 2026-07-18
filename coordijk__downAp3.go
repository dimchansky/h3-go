package h3

// _downAp3 finds the center of the hex at the next finer aperture 3
// counter-clockwise resolution. Works in place.
// Mirrors H3's coordijk.h::_downAp3 behavior.
// Ported from H3 C: coordijk.h::_downAp3.
func _downAp3(ijk *coordIJK) {
	// res r unit vectors in res r+1
	iVec := coordIJK{2, 0, 1}
	jVec := coordIJK{1, 2, 0}
	kVec := coordIJK{0, 1, 2}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
