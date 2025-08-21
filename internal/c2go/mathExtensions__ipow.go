package c2go

// _ipow performs integer exponentiation using exponentiation by squaring.
// Ported from H3 C v4.3.0: testref/h3-4.3.0/src/h3lib/lib/mathExtensions.c
// Signature preserved where possible.
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
