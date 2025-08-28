package h3

// _faceIjkToBaseCell finds the base cell given a FaceIJK coordinate.
// Given the face number and a resolution 0 ijk+ coordinate in that face's
// face-centered ijk coordinate system, return the base cell number.
// Valid ijk+ lookup coordinates are from (0, 0, 0) to (2, 2, 2).
// Ported from H3 C: baseCells.c::_faceIjkToBaseCell
func _faceIjkToBaseCell(h *FaceIJK) int32 {
	return faceIjkBaseCells[h.Face][h.Coord.I][h.Coord.J][h.Coord.K].BaseCell
}
