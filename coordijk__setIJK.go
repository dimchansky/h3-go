package h3

// _setIJK sets IJK coordinate components.
// Mirrors H3's coordijk.c::_setIJK behavior.
// Ported from H3 C: coordijk.c::_setIJK.
//
//nolint:unparam // signature mirrors H3 C _setIJK
func _setIJK(ijk *CoordIJK, i, j, k int32) {
	ijk.I = i
	ijk.J = j
	ijk.K = k
}
