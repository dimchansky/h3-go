package c2go

// cellToChildPos returns the position of child within an ordered list of all
// children of its parent at parentRes. Ports H3_EXPORT(cellToChildPos).
func cellToChildPos(child H3Index, parentRes int) (int64, H3Error) {
	childRes := getResolution(child)
	// Get parent at res to validate
	originalParent, perr := cellToParent(child, parentRes)
	if perr != E_SUCCESS {
		return 0, perr
	}

	parent := originalParent
	parentIsPentagon := isPentagon(parent) != 0

	var out int64 = 0
	if parentIsPentagon {
		for res := childRes; res > parentRes; res-- {
			// update parent one level up
			p, perr := cellToParent(child, res-1)
			if perr != E_SUCCESS {
				return 0, perr
			}
			parent = p
			parentIsPentagon = isPentagon(parent) != 0

			rawDigit := getIndexDigit(child, res)
			if rawDigit == INVALID_DIGIT || (parentIsPentagon && rawDigit == K_AXES_DIGIT) {
				return 0, H3Error(5) // E_CELL_INVALID = 5 per H3ErrorDescriptions
			}
			digit := rawDigit
			if parentIsPentagon && rawDigit > 0 {
				digit = rawDigit - 1
			}
			if digit != CENTER_DIGIT {
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
			digit := getIndexDigit(child, res)
			if digit == INVALID_DIGIT {
				return 0, H3Error(5) // E_CELL_INVALID
			}
			out += int64(digit) * _ipow(7, int64(childRes-res))
		}
	}

	if validateChildPos(out, originalParent, childRes) != E_SUCCESS {
		return 0, E_FAILED
	}
	return out, E_SUCCESS
}
