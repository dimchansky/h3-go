package h3

// _ijkRotate60ccw rotates IJK coordinates 60 degrees counter-clockwise.
// Mirrors H3's coordijk.h::_ijkRotate60ccw behavior.
// Ported from H3 C: coordijk.h::_ijkRotate60ccw.
func _ijkRotate60ccw(ijk *coordIJK) {
	// unit vector rotations
	iVec := coordIJK{I: 1, J: 1, K: 0}
	jVec := coordIJK{I: 0, J: 1, K: 1}
	kVec := coordIJK{I: 1, J: 0, K: 1}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
