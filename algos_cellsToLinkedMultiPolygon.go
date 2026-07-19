package h3

// cellsToLinkedMultiPolygon creates a LinkedGeoPolygon describing the
// outline(s) of a set of hexagons. Polygon outlines will follow GeoJSON
// MultiPolygon order: Each polygon will have one outer loop, which is
// first in the list, followed by any holes.
//
// All cells in the set must be valid, have the same resolution, and
// contain no duplicates; violations fail with the validateCellSet
// errors (H3 4.5.0 behavioral change, record §7.1).
//
// It is the responsibility of the caller to call
// destroyLinkedMultiPolygon on the populated linked geo structure, or
// the memory for that structure will not be freed (in Go, the garbage
// collector frees it).
//
// The H3 4.5.0 implementation delegates to the arc-cancellation
// cellsToMultiPolygon pipeline and converts its flat GeoMultiPolygon
// output to the linked form; the output is normalized by construction
// (no normalizeMultiPolygon call on this path).
//
// Ported from H3 C: algos.c::cellsToLinkedMultiPolygon.
func cellsToLinkedMultiPolygon(h3Set []h3Index, numHexes int32, out *linkedGeoPolygon) h3Error {
	var mpoly geoMultiPolygon
	err := cellsToMultiPolygon(h3Set, int64(numHexes), &mpoly)
	if err != eSuccess {
		return err
	}
	err = geoMultiPolygonToLinkedGeoPolygon(&mpoly, out)
	destroyGeoMultiPolygon(&mpoly)
	return err
}
