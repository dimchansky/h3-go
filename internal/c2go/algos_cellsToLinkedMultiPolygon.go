package c2go

// cellsToLinkedMultiPolygon creates a LinkedGeoPolygon describing the outline(s) of a set of hexagons.
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
// Ported from H3 C: algos.c::cellsToLinkedMultiPolygon
func cellsToLinkedMultiPolygon(h3Set []H3Index, numHexes int32, out *LinkedGeoPolygon) H3Error {
	var graph VertexGraph

	// Create vertex graph from the hexagon set
	err := h3SetToVertexGraph(h3Set, numHexes, &graph)
	if err != E_SUCCESS {
		return err
	}

	// Convert vertex graph to linked geometry
	_vertexGraphToLinkedGeo(&graph, out)

	// Clean up the vertex graph
	destroyVertexGraph(&graph)

	// Normalize the multi-polygon structure
	normalizeResult := normalizeMultiPolygon(out)
	if normalizeResult != E_SUCCESS {
		// If normalization fails, clean up the output structure
		destroyLinkedMultiPolygon(out)
	}

	return normalizeResult
}
