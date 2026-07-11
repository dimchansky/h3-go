package h3

// constructCell creates a cell from its components (resolution, base cell,
// children digits). Only allows for constructing valid H3 cells.
//
// digits is an array of child digits (0--6) of length res; nil is allowed
// for res=0.
// Ported from H3 C: h3Index.c::constructCell.
func constructCell(res int32, baseCellNumber int32, digits []int32, out *h3Index) h3Error {
	if res < 0 || res > maxH3Res {
		return eResDomain
	}
	if baseCellNumber < 0 || baseCellNumber >= numBaseCells {
		return eBaseCellDomain
	}

	h := h3Index(h3Init)
	h = setMode(h, h3CellMode)
	h = setResolution(h, res)
	h = setBaseCell(h, baseCellNumber)

	isPent := isBaseCellPentagonArr[baseCellNumber]

	for r := int32(1); r <= res; r++ {
		d := direction(digits[r-1])
		if d < centerDigit || d >= invalidDigit { // (d < 0 || d >= 7)
			return eDigitDomain
		}
		if isPent {
			// check for deleted subsequences of pentagons
			switch d {
			case centerDigit: // d == 0
				// do nothing; still a pentagon
			case kAxesDigit: // d == 1
				return eDeletedDigit
			default:
				isPent = false
			}
		}
		h = h3SetIndexDigit(h, r, int32(d))
	}

	*out = h
	return eSuccess
}
