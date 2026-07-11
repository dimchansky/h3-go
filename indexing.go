package h3

// checkRes validates a public resolution argument before it is narrowed to
// the C layer's int32 (guarding against wrap-around on narrowing).
func checkRes(res int) error {
	if res < 0 || res > MaxResolution {
		return ErrResolutionDomain
	}
	return nil
}

// LatLngToCell returns the cell containing the given coordinate at the given
// resolution.
//
// H3 C API: latLngToCell.
func LatLngToCell(g LatLng, res int) (Cell, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	var out h3Index
	if errC := latLngToCell(&g, int32(res), &out); errC != eSuccess {
		return 0, toErr(errC)
	}
	return out, nil
}

// LatLng returns the geographic center point of the cell.
//
// H3 C API: cellToLatLng.
func (c Cell) LatLng() (LatLng, error) {
	var g LatLng
	if errC := cellToLatLng(c, &g); errC != eSuccess {
		return LatLng{}, toErr(errC)
	}
	return g, nil
}

// Boundary returns the cell's boundary in counterclockwise order: 6 vertices
// for hexagons, 5 for pentagons (up to 10 for distorted pentagon boundaries).
// The returned value involves no heap allocation.
//
// H3 C API: cellToBoundary.
func (c Cell) Boundary() (CellBoundary, error) {
	var cb CellBoundary
	if errC := cellToBoundary(c, &cb); errC != eSuccess {
		return CellBoundary{}, toErr(errC)
	}
	return cb, nil
}

// Resolution returns the cell's resolution (0..MaxResolution). It is a pure
// bit accessor and does not validate the index.
//
// H3 C API: getResolution.
func (c Cell) Resolution() int { return int(getResolution(c)) }

// BaseCellNumber returns the number (0..NumBaseCells-1) of the resolution-0
// base cell the cell descends from. It is a pure bit accessor and does not
// validate the index.
//
// H3 C API: getBaseCellNumber.
func (c Cell) BaseCellNumber() int { return int(getBaseCellNumber(c)) }

// IsValid reports whether the index is a valid H3 cell index.
//
// H3 C API: isValidCell.
func (c Cell) IsValid() bool { return isValidCell(c) }

// IsPentagon reports whether the cell is one of the 12 pentagonal cells of
// its resolution.
//
// H3 C API: isPentagon.
func (c Cell) IsPentagon() bool { return isPentagon(c) }

// IsResClassIII reports whether the cell's resolution is a Class III
// resolution (rotated ~19.1° relative to Class II resolutions).
//
// H3 C API: isResClassIII.
func (c Cell) IsResClassIII() bool { return isResClassIII(c) }

// IcosahedronFaces returns the icosahedron face numbers (0..19) intersected
// by the cell, in ascending order. Cells intersect at most 5 faces.
//
// H3 C API: getIcosahedronFaces (sized via maxFaceCount).
func (c Cell) IcosahedronFaces() ([]int, error) {
	var cnt int32
	if errC := maxFaceCount(c, &cnt); errC != eSuccess {
		return nil, toErr(errC)
	}
	var buf [5]int32
	out := buf[:cnt]
	if errC := getIcosahedronFaces(c, out); errC != eSuccess {
		return nil, toErr(errC)
	}
	faces := make([]int, 0, cnt)
	for _, f := range out {
		if f != invalidFace {
			faces = append(faces, int(f))
		}
	}
	return faces, nil
}

// IndexDigit returns the indexing digit of the cell at the given resolution
// (1..MaxResolution; resolution 0 is the base cell number, not a digit). res
// may exceed the cell's actual resolution, in which case the stored digit
// (7 for valid cells) is returned.
//
// H3 C API: getIndexDigit (added in H3 4.4.0).
func (c Cell) IndexDigit(res int) (int, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	var d int32
	if errC := getIndexDigit(c, int32(res), &d); errC != eSuccess {
		return 0, toErr(errC)
	}
	return int(d), nil
}

// ConstructCell creates a cell from its components: a resolution, a base
// cell number (0..NumBaseCells-1), and res child digits (each 0..6). Only
// valid cells can be constructed; invalid digit sequences fail with
// ErrDigitDomain or ErrDeletedDigit.
//
// H3 C API: constructCell (added in H3 4.4.0).
func ConstructCell(res, baseCellNumber int, digits []int) (Cell, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	if len(digits) < res {
		return 0, ErrDigitDomain
	}
	var d32 [maxH3Res]int32
	for i := range res {
		d := digits[i]
		if d < 0 || d > int(invalidDigit) {
			return 0, ErrDigitDomain
		}
		d32[i] = int32(d)
	}
	var out h3Index
	if errC := constructCell(int32(res), int32(baseCellNumber), d32[:res], &out); errC != eSuccess {
		return 0, toErr(errC)
	}
	return out, nil
}

// IsValidIndex reports whether the raw 64-bit index is valid for any H3 mode
// (cell, directed edge, or vertex).
//
// H3 C API: isValidIndex (added in H3 4.4.0).
func IsValidIndex(raw uint64) bool { return isValidIndex(h3Index(raw)) }

// Pentagons returns the 12 pentagonal cells at the given resolution.
//
// H3 C API: getPentagons.
func Pentagons(res int) ([]Cell, error) {
	if err := checkRes(res); err != nil {
		return nil, err
	}
	out := make([]Cell, NumPentagons)
	if errC := getPentagons(int32(res), out); errC != eSuccess {
		return nil, toErr(errC)
	}
	return out, nil
}

// Res0Cells returns all 122 resolution-0 cells (the base cells), in ascending
// base-cell order.
//
// H3 C API: getRes0Cells.
func Res0Cells() []Cell {
	out := make([]Cell, NumRes0Cells)
	// getRes0Cells cannot fail with a correctly sized buffer.
	_ = getRes0Cells(out)
	return out
}

// NumCells returns the number of unique H3 cells at the given resolution.
//
// H3 C API: getNumCells.
func NumCells(res int) (int64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	n, errC := getNumCells(int32(res))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return n, nil
}
