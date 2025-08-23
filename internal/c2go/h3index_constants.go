package c2go

// Bit manipulation constants for H3Index validation and digit checking.
// These constants are used in _hasAny7UptoRes and related functions.
const (
	// MHI is a bit mask with alternating 100 pattern for each 3-bit digit position.
	// Used to efficiently check for invalid digits (7) without looping.
	// Pattern: 100100100100100100100100100100100100100100100 (binary)
	H3_DIGIT_CHECK_MHI = 0b100100100100100100100100100100100100100100100
	
	// MLO is MHI shifted right by 2 bits, used in the digit validation algorithm.
	H3_DIGIT_CHECK_MLO = H3_DIGIT_CHECK_MHI >> 2
)