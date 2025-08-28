package h3

// H3 index bit layout (from h3Index.h)
const (
	H3_NUM_BITS         = 64
	H3_MAX_OFFSET       = 63
	H3_MODE_OFFSET      = 59
	H3_BC_OFFSET        = 45
	H3_RES_OFFSET       = 52
	H3_RESERVED_OFFSET  = 56
	H3_PER_DIGIT_OFFSET = 3

	H3_HIGH_BIT_MASK          uint64 = 1 << H3_MAX_OFFSET
	H3_HIGH_BIT_MASK_NEGATIVE uint64 = ^H3_HIGH_BIT_MASK

	H3_MODE_MASK          uint64 = 15 << H3_MODE_OFFSET
	H3_MODE_MASK_NEGATIVE uint64 = ^H3_MODE_MASK

	H3_BC_MASK          uint64 = 127 << H3_BC_OFFSET
	H3_BC_MASK_NEGATIVE uint64 = ^H3_BC_MASK

	H3_RES_MASK          uint64 = 15 << H3_RES_OFFSET
	H3_RES_MASK_NEGATIVE uint64 = ^H3_RES_MASK

	H3_RESERVED_MASK          uint64 = 7 << H3_RESERVED_OFFSET
	H3_RESERVED_MASK_NEGATIVE uint64 = ^H3_RESERVED_MASK

	H3_DIGIT_MASK          uint64 = 7
	H3_DIGIT_MASK_NEGATIVE uint64 = ^H3_DIGIT_MASK

	// H3_INIT: mode=cell, res=0, base cell=0, digits all 7
	H3_INIT uint64 = 35184372088831

	// H3_NULL represents the null H3 index
	H3_NULL H3Index = 0
)

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

// isBaseCellPentagonArr mirrors the compact array used in h3Index.c for pentagon base cells.
// Size 128 for safe indexing; only first 122 are valid base cells.
var isBaseCellPentagonArr = [128]bool{
	/* 0-3 */ false, false, false, false,
	/* 4 */ true,
	/* 5-13 */ false, false, false, false, false, false, false, false, false,
	/* 14 */ true,
	/* 15-23 */ false, false, false, false, false, false, false, false, false,
	/* 24 */ true,
	/* 25-37 */ false, false, false, false, false, false, false, false, false, false, false, false, false,
	/* 38 */ true,
	/* 39-48 */ false, false, false, false, false, false, false, false, false, false,
	/* 49 */ true,
	/* 50-57 */ false, false, false, false, false, false, false, false,
	/* 58 */ true,
	/* 59-62 */ false, false, false, false,
	/* 63 */ true,
	/* 64-71 */ false, false, false, false, false, false, false, false,
	/* 72 */ true,
	/* 73-82 */ false, false, false, false, false, false, false, false, false, false,
	/* 83 */ true,
	/* 84-96 */ false, false, false, false, false, false, false, false, false, false, false, false, false,
	/* 97 */ true,
	/* 98-106 */ false, false, false, false, false, false, false, false, false,
	/* 107 */ true,
	/* 108-116 */ false, false, false, false, false, false, false, false, false,
	/* 117 */ true,
}

// Note: C implementation lives in baseCells.c as _isBaseCellPentagon.

// h3ErrorDescriptions contains error message strings for each H3Error code.
// Mirrored from H3 C: h3Index.c::H3ErrorDescriptions
var h3ErrorDescriptions = [16]string{
	/* E_SUCCESS */ "Success",
	/* E_FAILED */ "The operation failed but a more specific error is not available",
	/* E_DOMAIN */ "Argument was outside of acceptable range",
	/* E_LATLNG_DOMAIN */ "Latitude or longitude arguments were outside of acceptable range",
	/* E_RES_DOMAIN */ "Resolution argument was outside of acceptable range",
	/* E_CELL_INVALID */ "Cell argument was not valid",
	/* E_DIR_EDGE_INVALID */ "Directed edge argument was not valid",
	/* E_UNDIR_EDGE_INVALID */ "Undirected edge argument was not valid",
	/* E_VERTEX_INVALID */ "Vertex argument was not valid",
	/* E_PENTAGON */ "Pentagon distortion was encountered",
	/* E_DUPLICATE_INPUT */ "Duplicate input",
	/* E_NOT_NEIGHBORS */ "Cell arguments were not neighbors",
	/* E_RES_MISMATCH */ "Cell arguments had incompatible resolutions",
	/* E_MEMORY_ALLOC */ "Memory allocation failed",
	/* E_MEMORY_BOUNDS */ "Bounds of provided memory were insufficient",
	/* E_OPTION_INVALID */ "Mode or flags argument was not valid",
}
