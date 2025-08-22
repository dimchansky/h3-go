package c2go

// _ijkRotate60ccw rotates IJK coordinates 60 degrees counter-clockwise.
// Mirrors H3's coordijk.c::_ijkRotate60ccw behavior.
func _ijkRotate60ccw(ijk *CoordIJK) {
	// unit vector rotations
	iVec := CoordIJK{I: 1, J: 1, K: 0}
	jVec := CoordIJK{I: 0, J: 1, K: 1}
	kVec := CoordIJK{I: 1, J: 0, K: 1}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
