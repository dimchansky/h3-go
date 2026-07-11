package h3

// childPosToCell returns the child cell at a given position under parent at childRes.
// Ports H3_EXPORT(childPosToCell).
// Ported from H3 C: h3Index.c::childPosToCell.
func childPosToCell(childPos int64, parent h3Index, childRes int32) (h3Index, h3Error) {
	if childRes < 0 || childRes > maxH3Res {
		return 0, eResDomain
	}
	parentRes := getResolution(parent)
	if childRes < parentRes {
		return 0, eResMismatch
	}
	if err := validateChildPos(childPos, parent, childRes); err != eSuccess {
		return 0, err
	}

	resOffset := childRes - parentRes
	child := parent
	idx := childPos

	// Set resolution to childRes
	child = setResolution(child, childRes)

	if isPentagon(parent) {
		inPent := true
		for res := int32(1); res <= resOffset; res++ {
			resWidth := _ipow(7, int64(resOffset-res))
			if inPent {
				pentWidth := 1 + (5*(resWidth-1))/6
				if idx < pentWidth {
					child = setIndexDigit(child, parentRes+res, 0)
				} else {
					idx -= pentWidth
					inPent = false
					child = setIndexDigit(child, parentRes+res, int32((idx/resWidth)+2))
					idx %= resWidth
				}
			} else {
				child = setIndexDigit(child, parentRes+res, int32(idx/resWidth))
				idx %= resWidth
			}
		}
	} else {
		for res := int32(1); res <= resOffset; res++ {
			resWidth := _ipow(7, int64(resOffset-res))
			child = setIndexDigit(child, parentRes+res, int32(idx/resWidth))
			idx %= resWidth
		}
	}

	return child, eSuccess
}
