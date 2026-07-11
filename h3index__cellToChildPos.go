package h3

// cellToChildPos returns the position of child within an ordered list of all
// children of its parent at parentRes. Ports H3_EXPORT(cellToChildPos).
// Ported from H3 C: h3Index.c::cellToChildPos.
func cellToChildPos(child h3Index, parentRes int32) (int64, h3Error) {
	childRes := getResolution(child)
	// Get parent at res to validate
	originalParent, perr := cellToParent(child, parentRes)
	if perr != eSuccess {
		return 0, perr
	}

	parent := originalParent
	parentIsPentagon := isPentagon(parent)

	var out int64 = 0
	if parentIsPentagon {
		for res := childRes; res > parentRes; res-- {
			// update parent one level up
			p, perr := cellToParent(child, res-1)
			if perr != eSuccess {
				return 0, perr
			}
			parent = p
			parentIsPentagon = isPentagon(parent)

			rawDigit := direction(getIndexDigit(child, res))
			if rawDigit == invalidDigit || (parentIsPentagon && rawDigit == kAxesDigit) {
				return 0, h3Error(5) // eCellInvalid = 5 per H3ErrorDescriptions
			}
			digit := rawDigit
			if parentIsPentagon && rawDigit > 0 {
				digit = rawDigit - 1
			}
			if digit != centerDigit {
				hexChildCount := _ipow(7, int64(childRes-res))
				var offset int64
				if parentIsPentagon {
					offset = 1 + (5*(hexChildCount-1))/6
				} else {
					offset = hexChildCount
				}
				out += offset + int64(digit-1)*hexChildCount
			}
		}
	} else {
		for res := childRes; res > parentRes; res-- {
			digit := direction(getIndexDigit(child, res))
			if digit == invalidDigit {
				return 0, h3Error(5) // eCellInvalid
			}
			out += int64(digit) * _ipow(7, int64(childRes-res))
		}
	}

	if validateChildPos(out, originalParent, childRes) != eSuccess {
		return 0, eFailed
	}
	return out, eSuccess
}
