package c2go

// _ijkRotate60cw rotates IJK coordinates 60 degrees clockwise.
// Mirrors H3's coordijk.c::_ijkRotate60cw behavior.
// Ported from H3 C: coordijk.c::ijkRotate60cw
func _ijkRotate60cw(ijk *CoordIJK) {
	// unit vector rotations
	iVec := CoordIJK{I: 1, J: 0, K: 1}
	jVec := CoordIJK{I: 1, J: 1, K: 0}
	kVec := CoordIJK{I: 0, J: 1, K: 1}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
