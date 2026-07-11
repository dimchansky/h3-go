package h3

// _ijkSub subtracts two IJK coordinates (h1 - h2).
// Mirrors H3's coordijk.c::_ijkSub behavior.
// Ported from H3 C: coordijk.c::_ijkSub.
func _ijkSub(h1, h2, diff *coordIJK) {
	diff.I = h1.I - h2.I
	diff.J = h1.J - h2.J
	diff.K = h1.K - h2.K
}
