package h3

//lint:file-ignore U1000 Ignore all unused code

// ipow does integer exponentiation efficiently. Taken from StackOverflow.
// `base` is the integer base.
// `exp` is the integer exponent.
// Returns the exponentiated value.
func ipow(base int64, exp int) int64 {
	result := int64(1)
	for exp != 0 {
		if exp&1 != 0 {
			result *= base
		}
		exp >>= 1
		base *= base
	}

	return result
}
