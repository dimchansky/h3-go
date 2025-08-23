package c2go

// _hasAll7AfterRes checks that all unused digits after `res` are set to 7 (INVALID_DIGIT).
// Uses bit shifts to avoid looping through digits.
// Ported from H3 C: h3Index.c::_hasAll7AfterRes
func _hasAll7AfterRes(h H3Index, res int) bool {
	// NOTE: res check is needed because we can't shift by 64
	if res < 15 {
		shift := 19 + 3*res
		if shift >= 64 {
			return true // No digits to check beyond resolution 15
		}
		
		h = ^h
		h <<= shift
		h >>= shift
		
		return h == 0
	}
	return true
}