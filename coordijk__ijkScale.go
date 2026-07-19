package h3

// _ijkScale scales IJK coordinates by a factor.
// Mirrors H3's coordijk.h::_ijkScale behavior.
// Ported from H3 C: coordijk.h::_ijkScale.
func _ijkScale(c *coordIJK, factor int32) {
	c.I *= factor
	c.J *= factor
	c.K *= factor
}
