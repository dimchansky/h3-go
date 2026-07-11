package h3

// cellsToLinkedMultiPolygon creates a linkedGeoPolygon describing the outline(s) of a set of hexagons.
// Polygon outlines will follow GeoJSON MultiPolygon order: Each polygon will
// have one outer loop, which is first in the list, followed by any holes.
//
// It is the responsibility of the caller to call destroyLinkedMultiPolygon on
// the populated linked geo structure, or the memory for that structure will not
// be freed.
//
// It is expected that all hexagons in the set have the same resolution and
// that the set contains no duplicates. Behavior is undefined if duplicates
// or multiple resolutions are present, and the algorithm may produce
// unexpected or invalid output.
//
// Technical implementation details:
// 1. Creates a vertex graph from the hexagon set using h3SetToVertexGraph
// 2. Converts the vertex graph to linked geometry using _vertexGraphToLinkedGeo
// 3. Normalizes the result to follow GeoJSON MultiPolygon structure
// 4. Cleans up intermediate vertex graph memory
//
// Ported from H3 C: algos.c::cellsToLinkedMultiPolygon.
func cellsToLinkedMultiPolygon(h3Set []h3Index, numHexes int32, out *linkedGeoPolygon) h3Error {
	var graph vertexGraph

	// Create vertex graph from the hexagon set
	err := h3SetToVertexGraph(h3Set, numHexes, &graph)
	if err != eSuccess {
		return err
	}

	// Convert vertex graph to linked geometry
	_vertexGraphToLinkedGeo(&graph, out)

	// Clean up the vertex graph
	destroyVertexGraph(&graph)

	// Normalize the multi-polygon structure
	normalizeResult := normalizeMultiPolygon(out)
	if normalizeResult != eSuccess {
		// If normalization fails, clean up the output structure
		destroyLinkedMultiPolygon(out)
	}

	return normalizeResult
}
