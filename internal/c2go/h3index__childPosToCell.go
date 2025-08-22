package c2go

// childPosToCell returns the child cell at a given position under parent at childRes.
// Ports H3_EXPORT(childPosToCell).
func childPosToCell(childPos int64, parent H3Index, childRes int) (H3Index, uint32) {
    if childRes < 0 || childRes > MAX_H3_RES {
        return 0, _eResDomain
    }
    parentRes := getResolution(parent)
    if childRes < parentRes {
        return 0, _eResMismatch
    }
    if err := validateChildPos(childPos, parent, childRes); err != _eSuccess {
        return 0, err
    }

    resOffset := childRes - parentRes
    child := parent
    idx := childPos

    // Set resolution to childRes
    x := uint64(child)
    x &^= H3_RES_MASK
    x |= (uint64(childRes) & 15) << H3_RES_OFFSET
    child = H3Index(x)

    if isPentagon(parent) != 0 {
        inPent := true
        for res := 1; res <= resOffset; res++ {
            resWidth := _ipow(7, int64(resOffset-res))
            if inPent {
                pentWidth := 1 + (5*(resWidth-1))/6
                if idx < pentWidth {
                    child = setIndexDigit(child, parentRes+res, 0)
                } else {
                    idx -= pentWidth
                    inPent = false
                    child = setIndexDigit(child, parentRes+res, int((idx/resWidth)+2))
                    idx %= resWidth
                }
            } else {
                child = setIndexDigit(child, parentRes+res, int(idx/resWidth))
                idx %= resWidth
            }
        }
    } else {
        for res := 1; res <= resOffset; res++ {
            resWidth := _ipow(7, int64(resOffset-res))
            child = setIndexDigit(child, parentRes+res, int(idx/resWidth))
            idx %= resWidth
        }
    }

    return child, _eSuccess
}

