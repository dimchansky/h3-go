package h3

import "slices"

// Parent returns the ancestor of the cell at the given coarser resolution.
// The parent/child relationship is logical (index-hierarchical), not
// geometric: a descendant's boundary is not required to lie within its
// ancestor's boundary.
//
// Requesting the cell's own resolution returns the cell itself; a
// resolution finer than the cell's fails with ErrResolutionMismatch, and
// res outside 0..MaxResolution fails with ErrResolutionDomain.
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

// ImmediateParent returns the cell's parent one resolution coarser. A
// resolution-0 cell has no parent and returns ErrResolutionDomain.
//
// H3 C API: cellToParent.
func (c Cell) ImmediateParent() (Cell, error) {
	return c.Parent(c.Resolution() - 1)
}

// CenterChild returns the center descendant of the cell at the given finer
// resolution (the hierarchy is logical, not geometric — see Parent).
// Requesting the cell's own resolution returns the cell itself; a
// resolution coarser than the cell's, or outside 0..MaxResolution, fails
// with ErrResolutionDomain.
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
// finer resolution. Pentagons have fewer descendants than hexagons at the
// same depth (each pentagon level contributes 6 rather than 7 children), so
// the count depends on the cell. A resolution coarser than the cell's, or
// outside 0..MaxResolution, fails with ErrResolutionDomain.
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

// Children returns all descendants of the cell at the given finer
// resolution, in canonical child order: the fixed deterministic order in
// which H3 C's cellToChildren produces children. It is the same order that
// ChildPos and ChildAtPos index, so Children(res)[i] is the cell
// ChildAtPos(i, res) returns; the concrete sequence is otherwise
// unspecified. Like Parent, the hierarchy is logical, not geometric — a
// child's boundary is not required to lie within this cell's boundary.
// Resolution errors are those of NumChildren (a res coarser than the
// cell's, or outside 0..MaxResolution, fails with ErrResolutionDomain).
//
// H3 C API: cellToChildren.
func (c Cell) Children(res int) ([]Cell, error) { return c.AppendChildren(nil, res) }

// ImmediateChildren returns the cell's children one resolution finer, in
// canonical child order (see Children): 7 for a hexagon and 6 for a
// pentagon. A resolution-MaxResolution cell has no children and returns
// ErrResolutionDomain.
//
// H3 C API: cellToChildren.
func (c Cell) ImmediateChildren() ([]Cell, error) {
	return c.AppendImmediateChildren(nil)
}

// AppendImmediateChildren appends the cell's children one resolution finer
// to dst and returns the extended slice, in canonical child order (see
// Children). A resolution-MaxResolution cell has no children and fails
// with ErrResolutionDomain. Pass dst[:0] (or nil) to reuse capacity; a
// capacity of 7 is always sufficient and makes the call allocation-free.
// On error the returned slice has dst's original length and elements.
//
// H3 C API: cellToChildren.
func (c Cell) AppendImmediateChildren(dst []Cell) ([]Cell, error) {
	if err := checkRes(c.Resolution() + 1); err != nil {
		return dst, err
	}
	n := 7
	if isPentagon(c) {
		n = 6
	}
	start := len(dst)
	dst = slices.Grow(dst, n)[:start+n]
	i := start
	for digit := int32(0); digit <= 6; digit++ {
		if n == 6 && digit == int32(kAxesDigit) {
			continue
		}
		dst[i] = makeDirectChild(c, digit)
		i++
	}
	return dst, nil
}

// AppendChildren appends all descendants of the cell at the given finer
// resolution to dst and returns the extended slice, in canonical child
// order (see Children). Pass dst[:0] (or a nil slice) to reuse dst's
// capacity; when capacity suffices the call does not allocate. On error the
// returned slice has dst's original length and elements.
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

// ChildPos returns the position of the cell within the canonical child
// order (see Children) of its ancestor at the given coarser resolution.
// Positions round-trip with ChildAtPos. A parentRes finer than the cell's
// resolution fails with ErrResolutionMismatch; parentRes outside
// 0..MaxResolution fails with ErrResolutionDomain.
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

// ChildAtPos returns the child cell at the given position within the
// canonical child order (see Children) of c's descendants at the given
// finer resolution; it is the inverse of ChildPos. A position outside
// 0..NumChildren(res)-1 fails with ErrDomain, and res outside
// 0..MaxResolution fails with ErrResolutionDomain.
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
