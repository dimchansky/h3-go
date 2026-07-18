package h3

// _ijkRotate60cw rotates IJK coordinates 60 degrees clockwise.
// Mirrors H3's coordijk.h::_ijkRotate60cw behavior.
// Ported from H3 C: coordijk.h::_ijkRotate60cw.
func _ijkRotate60cw(ijk *coordIJK) {
	// unit vector rotations
	iVec := coordIJK{I: 1, J: 0, K: 1}
	jVec := coordIJK{I: 1, J: 1, K: 0}
	kVec := coordIJK{I: 0, J: 1, K: 1}

	_ijkScale(&iVec, ijk.I)
	_ijkScale(&jVec, ijk.J)
	_ijkScale(&kVec, ijk.K)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}
