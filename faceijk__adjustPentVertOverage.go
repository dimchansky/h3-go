package h3

// _adjustPentVertOverage adjusts a pentagon vertex faceIJK address for overage across an
// icosahedral face. This is a specialized wrapper around _adjustOverageClassII
// that handles pentagon vertices specifically.
// Ported from H3 C: faceijk.c::_adjustPentVertOverage.
//
//nolint:unparam // return value mirrors H3 C _adjustPentVertOverage
func _adjustPentVertOverage(fijk *faceIJK, res int32) overage {
	var overage overage

	// do-while loop: execute at least once, continue while overage == newFace
	for {
		overage = _adjustOverageClassII(fijk, res, false, true)
		if overage != newFace {
			break
		}
	}
	return overage
}
