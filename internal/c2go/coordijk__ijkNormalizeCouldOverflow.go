package c2go

// _ijkNormalizeCouldOverflow returns true if _ijkNormalize with the given input
// could have a signed integer overflow. Assumes k is set to 0.
// Mirrors H3's coordijk.c::_ijkNormalizeCouldOverflow behavior.
func _ijkNormalizeCouldOverflow(ijk *CoordIJK) bool {
	// Check for the possibility of overflow
	var max, min int
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
		if addInt32sOverflows(int32(max), int32(min)) {
			// max + min would overflow
			return true
		}
		if subInt32sOverflows(0, int32(min)) {
			// 0 - INT32_MIN would overflow
			return true
		}
		// Also check for max - min overflow (which happens when min is negative)
		if subInt32sOverflows(int32(max), int32(min)) {
			return true
		}
	}
	return false
}
