package c2go

// _setIJK sets IJK coordinate components.
// Mirrors H3's coordijk.c::_setIJK behavior.
// Ported from H3 C: coordijk.c::setIJK
func _setIJK(ijk *CoordIJK, i, j, k int) {
	ijk.I = i
	ijk.J = j
	ijk.K = k
}
