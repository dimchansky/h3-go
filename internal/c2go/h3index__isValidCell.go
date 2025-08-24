package c2go

// isValidCell determines if an H3 index is valid by checking bit patterns
// and structural constraints.
// Ported from H3 C: h3Index.c::isValidCell
func isValidCell(h H3Index) bool {
	/*
	   Look for bit patterns that would disqualify an H3Index from
	   being valid. If identified, exit early.

	   For reference the H3 index bit layout:

	   |   Region   | # bits |
	   |------------|--------|
	   | High       |      1 |
	   | Mode       |      4 |
	   | Reserved   |      3 |
	   | Resolution |      4 |
	   | Base Cell  |      7 |
	   | Digit 1    |      3 |
	   | Digit 2    |      3 |
	   | ...        |    ... |
	   | Digit 15   |      3 |

	   Speed benefits come from using bit manipulation instead of loops,
	   whenever possible.
	*/
	if !_hasGoodTopBits(h) {
		return false
	}

	// No need to check resolution; any 4 bits give a valid resolution.
	res := getResolution(h)

	// Get base cell number and check that it is valid.
	bc := getBaseCellNumber(h)
	if bc >= NUM_BASE_CELLS {
		return false
	}

	if _hasAny7UptoRes(h, res) {
		return false
	}
	if !_hasAll7AfterRes(h, res) {
		return false
	}
	if _hasDeletedSubsequence(h, bc) {
		return false
	}

	// If no disqualifications were identified, the index is a valid H3 cell.
	return true
}
