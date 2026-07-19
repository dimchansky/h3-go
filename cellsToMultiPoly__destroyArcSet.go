package h3

// destroyArcSet is the helper function to free memory allocated for an
// ArcSet. Safe to call with partially initialized ArcSet (NULL pointers
// are skipped). In Go the garbage collector frees the backing arrays;
// mirroring C's pointer-nulling keeps the control flow and observable
// state C-shaped.
// Ported from H3 C: cellsToMultiPoly.h::destroyArcSet.
func destroyArcSet(arcset *arcSet) {
	if arcset.arcs != nil {
		arcset.arcs = nil
	}
	if arcset.buckets != nil {
		arcset.buckets = nil
	}
}
