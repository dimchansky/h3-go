package h3

// DirectedEdge is an H3 directed edge index: a directed connection from an
// origin cell to one of its adjacent cells.
//
// The zero value is not a valid directed edge; IsValid reports false for it.
// A DirectedEdge may be constructed by conversion from a raw uint64 index
// (unchecked — use IsValid to verify), parsed with ParseDirectedEdge, or
// produced by operations such as Cell.DirectedEdgeTo.
type DirectedEdge uint64

// IsNeighbor reports whether the two cells are adjacent. A cell is never
// its own neighbor: comparing a cell with itself returns false with a nil
// error. It fails with ErrCellInvalid when either index is not in cell mode
// (a lightweight mode check, not full validation) and with
// ErrResolutionMismatch for cells of different resolutions.
//
// H3 C API: areNeighborCells.
func (c Cell) IsNeighbor(other Cell) (bool, error) {
	ok, errC := areNeighborCells(c, other)
	if errC != eSuccess {
		return false, toErr(errC)
	}
	return ok, nil
}

// DirectedEdgeTo returns the directed edge from c to the neighboring
// destination cell. It fails with ErrNotNeighbors when the cells are not
// adjacent.
//
// H3 C API: cellsToDirectedEdge.
func (c Cell) DirectedEdgeTo(destination Cell) (DirectedEdge, error) {
	e, errC := cellsToDirectedEdge(c, destination)
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return DirectedEdge(e), nil
}

// DirectedEdges returns all directed edges originating at c: 6 for hexagons,
// 5 for pentagons.
//
// H3 C API: originToDirectedEdges.
func (c Cell) DirectedEdges() ([]DirectedEdge, error) {
	var raw [6]h3Index
	if errC := originToDirectedEdges(c, raw[:]); errC != eSuccess {
		return nil, toErr(errC)
	}
	out := make([]DirectedEdge, 0, 6)
	for _, e := range raw {
		if e != h3Null { // pentagons leave the deleted k-axis slot empty
			out = append(out, DirectedEdge(e))
		}
	}
	return out, nil
}

// IsValid reports whether the index is a valid H3 directed edge index.
//
// H3 C API: isValidDirectedEdge.
func (e DirectedEdge) IsValid() bool { return isValidDirectedEdge(h3Index(e)) }

// Resolution returns the edge's resolution. It is a pure bit accessor and
// does not validate the index.
//
// H3 C API: getResolution.
func (e DirectedEdge) Resolution() int { return int(getResolution(h3Index(e))) }

// IndexDigit returns the indexing digit of the edge's origin cell at the
// given resolution (1..MaxResolution; resolution 0 is the base cell number,
// not a digit). res may exceed the edge's actual resolution, in which case
// the stored digit (7 for valid edges) is returned.
//
// H3 C API: getIndexDigit.
func (e DirectedEdge) IndexDigit(res int) (int, error) {
	return indexDigit(e, res)
}

// Origin returns the origin cell of the directed edge.
//
// H3 C API: getDirectedEdgeOrigin.
func (e DirectedEdge) Origin() (Cell, error) {
	c, errC := getDirectedEdgeOrigin(h3Index(e))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return c, nil
}

// Destination returns the destination cell of the directed edge.
//
// H3 C API: getDirectedEdgeDestination.
func (e DirectedEdge) Destination() (Cell, error) {
	var out h3Index
	if errC := getDirectedEdgeDestination(h3Index(e), &out); errC != eSuccess {
		return 0, toErr(errC)
	}
	return out, nil
}

// Cells returns the origin and destination cells of the directed edge, in
// that order.
//
// H3 C API: directedEdgeToCells.
func (e DirectedEdge) Cells() (origin, destination Cell, err error) {
	var pair [2]h3Index
	if errC := directedEdgeToCells(h3Index(e), pair[:]); errC != eSuccess {
		return 0, 0, toErr(errC)
	}
	return pair[0], pair[1], nil
}

// Boundary returns the geographic boundary of the directed edge: the
// vertices shared by its origin and destination. The boundary contains the
// edge's two topological endpoints and may contain one additional
// distortion vertex when the edge crosses an icosahedron face. The returned
// value involves no heap allocation.
//
// H3 C API: directedEdgeToBoundary.
func (e DirectedEdge) Boundary() (CellBoundary, error) {
	var cb CellBoundary
	if errC := directedEdgeToBoundary(h3Index(e), &cb); errC != eSuccess {
		return CellBoundary{}, toErr(errC)
	}
	return cb, nil
}
