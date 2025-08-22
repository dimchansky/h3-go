package c2go

// _downAp7 finds the center of the hex at the next finer aperture 7
// counter-clockwise resolution. Works in place.
// Mirrors H3's coordijk.c::_downAp7 behavior.
func _downAp7(ijk *CoordIJK) {
	// res r unit vectors in res r+1
	iVec := CoordIJK{3, 0, 1}
	jVec := CoordIJK{1, 3, 0}
	kVec := CoordIJK{0, 1, 3}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
