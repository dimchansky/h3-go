package c2go

// _hasDeletedSubsequence validates pentagon cells for invalid subsequences.
// Pentagon cells start with a sequence of 0's (CENTER_DIGIT's).
// The first nonzero digit can't be a 1 (i.e., "deleted subsequence",
// PENTAGON_SKIPPED_DIGIT, or K_AXES_DIGIT).
// Ported from H3 C: h3Index.c::_hasDeletedSubsequence
func _hasDeletedSubsequence(h H3Index, baseCell int) bool {
	if baseCell >= 0 && baseCell < len(isBaseCellPentagonArr) && isBaseCellPentagonArr[baseCell] {
		// Keep only the lower 45 bits (15 digits × 3 bits each)
		h <<= 19
		h >>= 19
		
		if h == 0 {
			return false // all zeros: res 15 pentagon
		}
		return _firstOneIndex(h)%3 == 0
	}
	return false
}