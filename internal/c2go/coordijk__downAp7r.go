package c2go

// _downAp7r finds the center of the hex at the next finer aperture 7
// clockwise resolution. Works in place.
// Mirrors H3's coordijk.c::_downAp7r behavior.
// Ported from H3 C: coordijk.c::downAp7r
func _downAp7r(ijk *CoordIJK) {
	// res r unit vectors in res r+1
	iVec := CoordIJK{3, 1, 0}
	jVec := CoordIJK{0, 3, 1}
	kVec := CoordIJK{1, 0, 3}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
