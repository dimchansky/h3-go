package h3

// _downAp3 finds the center of the hex at the next finer aperture 3
// counter-clockwise resolution. Works in place.
// Mirrors H3's coordijk.c::_downAp3 behavior.
// Ported from H3 C: coordijk.c::_downAp3
func _downAp3(ijk *CoordIJK) {
	// res r unit vectors in res r+1
	iVec := CoordIJK{2, 0, 1}
	jVec := CoordIJK{1, 2, 0}
	kVec := CoordIJK{0, 1, 2}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
