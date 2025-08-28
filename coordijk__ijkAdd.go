package h3

// _ijkAdd adds two IJK coordinates.
// Mirrors H3's coordijk.c::_ijkAdd behavior.
// Ported from H3 C: coordijk.c::_ijkAdd
func _ijkAdd(h1, h2, sum *CoordIJK) {
	sum.I = h1.I + h2.I
	sum.J = h1.J + h2.J
	sum.K = h1.K + h2.K
}
