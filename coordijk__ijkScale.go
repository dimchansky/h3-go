package h3

// _ijkScale scales IJK coordinates by a factor.
// Mirrors H3's coordijk.c::_ijkScale behavior.
// Ported from H3 C: coordijk.c::_ijkScale
func _ijkScale(c *CoordIJK, factor int32) {
	c.I *= factor
	c.J *= factor
	c.K *= factor
}
