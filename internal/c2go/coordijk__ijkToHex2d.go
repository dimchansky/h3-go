package c2go

// _ijkToHex2d converts IJK coordinates to 2D hex coordinates.
// Mirrors H3's coordijk.c::_ijkToHex2d behavior.
func _ijkToHex2d(h *CoordIJK, v *Vec2d) {
	i := h.I - h.K
	j := h.J - h.K
	v.X = float64(i) - 0.5*float64(j)
	v.Y = float64(j) * M_SQRT3_2
}
