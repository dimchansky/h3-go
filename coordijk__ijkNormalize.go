package h3

// _ijkNormalize normalizes IJK coordinates by removing negative values.
// Mirrors H3's coordijk.c::_ijkNormalize behavior.
// Ported from H3 C: coordijk.c::_ijkNormalize
func _ijkNormalize(c *CoordIJK) {
	// remove any negative values
	if c.I < 0 {
		c.J -= c.I
		c.K -= c.I
		c.I = 0
	}

	if c.J < 0 {
		c.I -= c.J
		c.K -= c.J
		c.J = 0
	}

	if c.K < 0 {
		c.I -= c.K
		c.J -= c.K
		c.K = 0
	}

	// remove the min value if needed
	min := c.I
	if c.J < min {
		min = c.J
	}
	if c.K < min {
		min = c.K
	}
	if min > 0 {
		c.I -= min
		c.J -= min
		c.K -= min
	}
}
