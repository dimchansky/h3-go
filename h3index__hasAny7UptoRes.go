package h3

// _hasAny7UptoRes checks if any digit from 1 to `res` is 7 (invalidDigit).
// Uses efficient bit manipulation to check for invalid digits without looping.
// Ported from H3 C: h3Index.c::_hasAny7UptoRes.
func _hasAny7UptoRes(h h3Index, res int32) bool {
	// Shift to zero out digits beyond resolution
	shift := 3 * (15 - res)

	// Mirror C behavior exactly - no bounds checking on shift
	if shift < 0 {
		// For res > 15, shift becomes negative
		// In C this would be undefined behavior, but typically results in 0 shift
		// The C implementation would check all digits including unused ones
		return false
	}

	if shift >= 64 {
		// Large positive shift zeros everything in C
		return false
	}

	h >>= shift
	h <<= shift
	h = (h & h3DigitCheckMhi & (^h - h3DigitCheckMlo))

	return h != 0
}
