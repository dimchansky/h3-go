package h3

// destroySortablePolyVerts is the helper function to free outer loop
// vertices from an array of SortablePoly. Frees the verts arrays from
// each polygon's geoloop, then the polygon array. In Go only the
// nested GeoLoop references are cleared: the slice is passed by value,
// so the caller-visible slice header cannot be nilled here — C's final
// free of the array has no Go equivalent; the garbage collector
// reclaims it when the caller drops its reference. Used during partial
// cleanup when constructing the polygon array fails. numPolys
// specifies how many polygons to clean up.
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
