package h3

// _downAp3r finds the center of the hex at the next finer aperture 3
// clockwise resolution. Works in place.
// Mirrors H3's coordijk.h::_downAp3r behavior.
// Ported from H3 C: coordijk.h::_downAp3r.
func _downAp3r(ijk *coordIJK) {
	// res r unit vectors in res r+1
	iVec := coordIJK{2, 1, 0}
	jVec := coordIJK{0, 2, 1}
	kVec := coordIJK{1, 0, 2}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
