package h3

// _ijkMatches reports whether two IJK coordinates are equal.
// Mirrors H3's coordijk.c::_ijkMatches behavior.
// Ported from H3 C: coordijk.c::_ijkMatches.
func _ijkMatches(c1, c2 *CoordIJK) bool {
	return c1.I == c2.I && c1.J == c2.J && c1.K == c2.K
}
