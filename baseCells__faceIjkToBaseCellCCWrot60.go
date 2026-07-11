package h3

// _faceIjkToBaseCellCCWrot60 finds the CCW rotation for a base cell given faceIJK.
// Given the face number and a resolution 0 ijk+ coordinate in that face's
// face-centered ijk coordinate system, return the number of 60' ccw rotations
// to rotate into the coordinate system of the base cell at that coordinates.
// Valid ijk+ lookup coordinates are from (0, 0, 0) to (2, 2, 2).
// Ported from H3 C: baseCells.c::_faceIjkToBaseCellCCWrot60.
func _faceIjkToBaseCellCCWrot60(h *faceIJK) int32 {
	return faceIjkBaseCells[h.Face][h.Coord.I][h.Coord.J][h.Coord.K].CcwRot60
}
