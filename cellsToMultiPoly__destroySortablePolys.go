package h3

// destroySortablePolys is the helper function to free memory allocated
// for an array of SortablePoly. Frees the holes arrays in each polygon,
// then the polygon array itself (in Go: nils the references). numPolys
// specifies how many polygons to clean up.
// Ported from H3 C: cellsToMultiPoly.h::destroySortablePolys.
func destroySortablePolys(spolys []sortablePoly, numPolys int64) {
	if spolys != nil {
		for i := int64(0); i < numPolys; i++ {
			if spolys[i].poly.Holes != nil {
				spolys[i].poly.Holes = nil
			}
		}
	}
}
