package h3

//lint:file-ignore U1000 Ignore all unused code

// _ipow does integer exponentiation efficiently. Taken from StackOverflow.
// `base` is the integer base.
// `exp` is the integer exponent.
// Returns the exponentiated value.
func _ipow(base int, exp int) int {
	result := 1
	for exp != 0 {
		if exp&1 != 0 {
			result *= base
		}
		exp >>= 1
		base *= base
	}

	return result
}
