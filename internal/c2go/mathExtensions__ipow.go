package c2go

// _ipow performs integer exponentiation using exponentiation by squaring.
// Signature preserved where possible.
// Ported from H3 C: mathExtensions.c::_ipow
func _ipow(base int64, exp int64) int64 {
	var result int64 = 1
	for exp != 0 {
		if (exp & 1) != 0 {
			result *= base
		}
		exp >>= 1
		base *= base
	}
	return result
}
