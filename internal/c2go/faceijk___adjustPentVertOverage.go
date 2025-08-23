package c2go

// _adjustPentVertOverage adjusts a pentagon vertex FaceIJK address for overage across an
// icosahedral face. This is a specialized wrapper around _adjustOverageClassII
// that handles pentagon vertices specifically.
// Ported from H3 C: faceijk.c::_adjustPentVertOverage
func _adjustPentVertOverage(fijk *FaceIJK, res int32) Overage {
	var overage Overage

	// do-while loop: execute at least once, continue while overage == NEW_FACE
	for {
		overage = _adjustOverageClassII(fijk, res, false, true)
		if overage != NEW_FACE {
			break
		}
	}
	return overage
}
