package c2go

// _ijkScale scales IJK coordinates by a factor.
// Mirrors H3's coordijk.c::_ijkScale behavior.
func _ijkScale(c *CoordIJK, factor int) {
	c.I *= factor
	c.J *= factor
	c.K *= factor
}
