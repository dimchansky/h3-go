package h3

// cellToEdgeArcs fills in edge arcs for a single cell:
//
//   - create one Arc for each edge of the cell
//   - set prev/next arcs in linked loop. ensures edges in
//     counter-clockwise order
//   - initialize parent and rank for union_find (each loop of cell edges
//     starts as its own separate connected component)
//
// arcs is the sub-slice of the ArcSet's arcs array reserved for this
// cell (C: an Arc* into that array).
// Ported from H3 C: cellsToMultiPoly.c::cellToEdgeArcs.
func cellToEdgeArcs(h h3Index, arcs []arc, numEdgesOut *int64) h3Error {
	var numEdges int64
	var _edges [6]h3Index
	var edges []h3Index

	idxh := [6]uint8{0, 4, 3, 5, 1, 2}
	idxp := [5]uint8{0, 1, 3, 2, 4}
	var idx []uint8

	err := originToDirectedEdges(h, _edges[:])
	if err != eSuccess {
		// NEVER in C: since already checked with validateCellSet, this
		// should never error.
		return err
	}

	// Set `edges` to contain the indices of cell edges in counter-clockwise
	// order the first directed edge of a pentagon is H3_NULL
	if _edges[0] == h3Null {
		numEdges = 5
		idx = idxp[:]
		edges = _edges[1:]
	} else {
		numEdges = 6
		idx = idxh[:]
		edges = _edges[:]
	}

	for i := int64(0); i < numEdges; i++ {
		// Arcs stay in same order as output of originToDirectedEdges.
		// That is, they are not in CCW order in the `arcs` array, but they
		// are in CCW in the linked loop.
		arcs[i].id = edges[i]
		arcs[i].isRemoved = false
		arcs[i].isVisited = false

		// initialize union-find datastructure
		// all edges in loop have same parent: first edge
		arcs[i].parent = &arcs[0]
		arcs[i].rank = 1

		// Connect so prev/next point to neighboring edges that share a vertex.
		// Edges/vertexes should follow right-hand rule as a result (CCW order).
		cur := int64(idx[i])
		prev := int64(idx[(i-1+numEdges)%numEdges])
		next := int64(idx[(i+1)%numEdges])
		arcs[cur].prev = &arcs[prev]
		arcs[cur].next = &arcs[next]
	}

	*numEdgesOut = numEdges
	return eSuccess
}
