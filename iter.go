package h3

import "iter"

// ChildrenSeq returns an iterator over all descendants of the cell at the
// given finer resolution, in canonical child order (see Children), without
// materializing them. The sequence is re-runnable — each range restarts
// from the first child, unlike C's single-pass iterator structs — and
// breaking early is safe. Iterating allocates nothing. If the cell or
// resolution is invalid the sequence is empty (matching the C iterator's
// null-iterator contract).
//
// H3 C API: iterInitParent / iterStepChild (iterators.h).
func (c Cell) ChildrenSeq(res int) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		if res < 0 || res > MaxResolution {
			return
		}
		var it iterCellsChildren
		iterInitParent(c, int32(res), &it)
		for ; it.H != h3Null; iterStepChild(&it) {
			if !yield(it.H) {
				return
			}
		}
	}
}

// CellsAtRes returns an iterator over all cells at the given resolution,
// without materializing them (there are NumCells(res) of them — over 500
// trillion at resolution 15). The iteration order is not part of the
// contract. The sequence is re-runnable — each range restarts from the
// first cell — and breaking early is safe. Iterating allocates nothing. If
// the resolution is invalid the sequence is empty.
//
// H3 C API: iterInitRes / iterStepRes (iterators.h).
func CellsAtRes(res int) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		if res < 0 || res > MaxResolution {
			return
		}
		it := iterInitRes(int32(res))
		for ; it.H != h3Null; iterStepRes(&it) {
			if !yield(it.H) {
				return
			}
		}
	}
}

// PolygonToCellsExperimentalSeq returns an iterator over the cells matching
// the polygon under the given containment mode (see ContainmentMode),
// without materializing the full result, in no particular order. Input
// validation happens before the sequence is returned (res outside
// 0..MaxResolution fails with ErrResolutionDomain, an invalid mode with
// ErrOptionInvalid); iteration itself cannot fail. The sequence is re-runnable — each range restarts the
// underlying iterator — and breaking early is safe (internal iterator
// state is released). Iteration may allocate internal iterator state; it
// is not allocation-free.
//
// The polygon is captured as a struct of slice headers: its coordinate
// slices are shared, not copied. Mutations made before a later iteration
// are observed by that iteration and may affect its results. Callers must
// not mutate the polygon's backing slices during an active iteration;
// concurrent mutation is unsafe.
//
// Like its C counterpart, this API is experimental and may change in minor
// versions.
//
// H3 C API: iterInitPolygon / iterStepPolygon (polyfill.h).
func PolygonToCellsExperimentalSeq(p GeoPolygon, res int, mode ContainmentMode) (iter.Seq[Cell], error) {
	if err := checkRes(res); err != nil {
		return nil, err
	}
	if errC := validatePolygonFlags(uint32(mode)); errC != eSuccess {
		return nil, toErr(errC)
	}
	return func(yield func(Cell) bool) {
		it := iterInitPolygon(&p, int32(res), uint32(mode))
		for ; it.Cell != h3Null; iterStepPolygon(&it) {
			if !yield(it.Cell) {
				iterDestroyPolygon(&it)
				return
			}
		}
	}, nil
}
