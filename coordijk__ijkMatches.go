package h3

// _ijkMatches reports whether two IJK coordinates are equal.
// Mirrors H3's coordijk.h::_ijkMatches behavior.
// Ported from H3 C: coordijk.h::_ijkMatches.
func _ijkMatches(c1, c2 *coordIJK) bool {
	return c1.I == c2.I && c1.J == c2.J && c1.K == c2.K
}
