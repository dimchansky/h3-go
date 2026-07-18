package h3

// reverseDirectedEdge returns the directed edge with origin and
// destination cells reversed.
// Ported from H3 C: directedEdge.c::reverseDirectedEdge.
func reverseDirectedEdge(edge h3Index, out *h3Index) h3Error {
	origin, err := getDirectedEdgeOrigin(edge)
	if err != eSuccess {
		return err
	}
	var destination h3Index
	err = getDirectedEdgeDestination(edge, &destination)
	if err != eSuccess {
		return err
	}
	*out, err = cellsToDirectedEdge(destination, origin)
	if err != eSuccess {
		return err
	}
	return eSuccess
}
