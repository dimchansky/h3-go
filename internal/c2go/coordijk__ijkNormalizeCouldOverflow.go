package c2go

// _ijkNormalizeCouldOverflow returns true if _ijkNormalize with the given input
// could have a signed integer overflow. Assumes k is set to 0.
// Mirrors H3's coordijk.c::_ijkNormalizeCouldOverflow behavior.
// Ported from H3 C: coordijk.c::ijkNormalizeCouldOverflow
func _ijkNormalizeCouldOverflow(ijk *CoordIJK) bool {
	// Check for the possibility of overflow
	var max, min int32
	if ijk.I > ijk.J {
		max = ijk.I
		min = ijk.J
	} else {
		max = ijk.J
		min = ijk.I
	}

	if min < 0 {
		// Only if the min is less than 0 will the resulting number be larger
		// than max. If min is positive, then max is also positive, and a
		// positive signed integer minus another positive signed integer will
		// not overflow.
		if addInt32sOverflows(max, min) {
			// max + min would overflow
			return true
		}
		if subInt32sOverflows(0, min) {
			// 0 - INT32_MIN would overflow
			return true
		}
		// Also check for max - min overflow (which happens when min is negative)
		if subInt32sOverflows(max, min) {
			return true
		}
	}
	return false
}
