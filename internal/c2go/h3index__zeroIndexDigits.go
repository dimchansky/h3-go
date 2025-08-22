package c2go

// _zeroIndexDigits clears index digits from start to end (inclusive).
// This mirrors the C helper _zeroIndexDigits.
func _zeroIndexDigits(h H3Index, start, end int) H3Index {
	if start > end {
		return h
	}
	// Allow start at 0 to mirror C behavior (affects bits overlapping base cell).
	if end > MAX_H3_RES {
		// Mirror C: out-of-range end results in no-op.
		return h
	}
	for r := start; r <= end; r++ {
		h = setIndexDigit(h, r, 0)
	}
	return h
}
