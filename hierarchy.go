package h3

import "slices"

// Parent returns the ancestor of the cell at the given coarser resolution.
//
// H3 C API: cellToParent.
func (c Cell) Parent(res int) (Cell, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	p, errC := cellToParent(c, int32(res))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return p, nil
}

// CenterChild returns the center descendant of the cell at the given finer
// resolution.
//
// H3 C API: cellToCenterChild.
func (c Cell) CenterChild(res int) (Cell, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	ch, errC := cellToCenterChild(c, int32(res))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return ch, nil
}

// NumChildren returns the number of descendants of the cell at the given
// finer resolution.
//
// H3 C API: cellToChildrenSize.
func (c Cell) NumChildren(res int) (int64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	n, errC := cellToChildrenSize(c, int32(res))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return n, nil
}

// Children returns all descendants of the cell at the given finer resolution,
// in canonical child order.
//
// H3 C API: cellToChildren.
func (c Cell) Children(res int) ([]Cell, error) { return c.AppendChildren(nil, res) }

// AppendChildren appends all descendants of the cell at the given finer
// resolution to dst and returns the extended slice, in canonical child order.
// Pass dst[:0] (or a nil slice) to reuse dst's capacity; when capacity
// suffices the call does not allocate.
//
// H3 C API: cellToChildren.
func (c Cell) AppendChildren(dst []Cell, res int) ([]Cell, error) {
	n, err := c.NumChildren(res)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst = slices.Grow(dst, int(n))[:start+int(n)]
	if errC := cellToChildren(c, int32(res), dst[start:]); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst, nil
}

// ChildPos returns the position of the cell within an ordered list of all
// children of the cell's ancestor at the given coarser resolution.
//
// H3 C API: cellToChildPos.
func (c Cell) ChildPos(parentRes int) (int64, error) {
	if err := checkRes(parentRes); err != nil {
		return 0, err
	}
	pos, errC := cellToChildPos(c, int32(parentRes))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return pos, nil
}

// ChildAtPos returns the child cell at the given position within an ordered
// list of all descendants of c at the given finer resolution; it is the
// inverse of ChildPos.
//
// H3 C API: childPosToCell.
func (c Cell) ChildAtPos(pos int64, res int) (Cell, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	ch, errC := childPosToCell(pos, c, int32(res))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return ch, nil
}
