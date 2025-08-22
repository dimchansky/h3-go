package c2go

// _h3LeadingNonZeroDigit returns the highest resolution non-zero digit.
func _h3LeadingNonZeroDigit(h H3Index) int {
	res := getResolution(h)
	for r := 1; r <= res; r++ {
		d := getIndexDigit(h, r)
		if d != 0 {
			return d
		}
	}
	// All zeros => CENTER_DIGIT (0)
	return 0
}
