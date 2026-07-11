package h3

// _ijkNormalizeCouldOverflow returns true if _ijkNormalize with the given input
// could have a signed integer overflow. Assumes k is set to 0.
// Mirrors H3's coordijk.c::_ijkNormalizeCouldOverflow behavior.
// Ported from H3 C: coordijk.c::_ijkNormalizeCouldOverflow.
func _ijkNormalizeCouldOverflow(ijk *coordIJK) bool {
	// Check for the possibility of overflow
	var maxVal, minVal int32
	if ijk.I > ijk.J {
		maxVal = ijk.I
		minVal = ijk.J
	} else {
		maxVal = ijk.J
		minVal = ijk.I
	}

	if minVal < 0 {
		// Only if the minVal is less than 0 will the resulting number be larger
		// than maxVal. If minVal is positive, then maxVal is also positive, and a
		// positive signed integer minus another positive signed integer will
		// not overflow.
		if addInt32sOverflows(maxVal, minVal) {
			// maxVal + minVal would overflow
			return true
		}
		if subInt32sOverflows(0, minVal) {
			// 0 - int32Min would overflow
			return true
		}
		// Also check for maxVal - minVal overflow (which happens when minVal is negative)
		if subInt32sOverflows(maxVal, minVal) {
			return true
		}
	}
	return false
}
