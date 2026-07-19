package h3

// destroySortablePolyVerts is the helper function to free outer loop
// vertices from an array of SortablePoly. Frees the verts arrays from
// each polygon's geoloop, then the polygon array (in Go: nils the
// references). Used during partial cleanup when constructing the
// polygon array fails. numPolys specifies how many polygons to clean
// up.
// Ported from H3 C: cellsToMultiPoly.h::destroySortablePolyVerts.
func destroySortablePolyVerts(spolys []sortablePoly, numPolys int64) {
	if spolys != nil {
		for i := int64(0); i < numPolys; i++ {
			if spolys[i].poly.GeoLoop != nil {
				spolys[i].poly.GeoLoop = nil
			}
		}
	}
}
