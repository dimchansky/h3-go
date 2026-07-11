package h3

// H3 index bit layout (from h3Index.h).
const (
	h3NumBits        = 64
	h3MaxOffset      = 63
	h3ModeOffset     = 59
	h3BcOffset       = 45
	h3ResOffset      = 52
	h3ReservedOffset = 56
	h3PerDigitOffset = 3

	h3HighBitMask         uint64 = 1 << h3MaxOffset
	h3HighBitMaskNegative uint64 = ^h3HighBitMask

	h3ModeMask         uint64 = 15 << h3ModeOffset
	h3ModeMaskNegative uint64 = ^h3ModeMask

	h3BcMask         uint64 = 127 << h3BcOffset
	h3BcMaskNegative uint64 = ^h3BcMask

	h3ResMask         uint64 = 15 << h3ResOffset
	h3ResMaskNegative uint64 = ^h3ResMask

	h3ReservedMask         uint64 = 7 << h3ReservedOffset
	h3ReservedMaskNegative uint64 = ^h3ReservedMask

	h3DigitMask         uint64 = 7
	h3DigitMaskNegative uint64 = ^h3DigitMask

	// h3Init: mode=cell, res=0, base cell=0, digits all 7.
	h3Init uint64 = 35184372088831

	// h3Null represents the null H3 index.
	h3Null h3Index = 0
)

// Bit manipulation constants for h3Index validation and digit checking.
// These constants are used in _hasAny7UptoRes and related functions.
const (
	// MHI is a bit mask with alternating 100 pattern for each 3-bit digit position.
	// Used to efficiently check for invalid digits (7) without looping.
	// Pattern: 100100100100100100100100100100100100100100100 (binary).
	h3DigitCheckMhi = 0b100100100100100100100100100100100100100100100

	// MLO is MHI shifted right by 2 bits, used in the digit validation algorithm.
	h3DigitCheckMlo = h3DigitCheckMhi >> 2
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

// h3ErrorDescriptions contains error message strings for each h3Error code.
// Mirrored from H3 C: h3Index.c::H3ErrorDescriptions.
var h3ErrorDescriptions = [h3ErrorEnd]string{
	/* eSuccess */ "Success",
	/* eFailed */ "The operation failed but a more specific error is not available",
	/* eDomain */ "Argument was outside of acceptable range",
	/* eLatlngDomain */ "Latitude or longitude arguments were outside of acceptable range",
	/* eResDomain */ "Resolution argument was outside of acceptable range",
	/* eCellInvalid */ "Cell argument was not valid",
	/* eDirEdgeInvalid */ "Directed edge argument was not valid",
	/* eUndirEdgeInvalid */ "Undirected edge argument was not valid",
	/* eVertexInvalid */ "Vertex argument was not valid",
	/* ePentagon */ "Pentagon distortion was encountered",
	/* eDuplicateInput */ "Duplicate input",
	/* eNotNeighbors */ "Cell arguments were not neighbors",
	/* eResMismatch */ "Cell arguments had incompatible resolutions",
	/* eMemoryAlloc */ "Memory allocation failed",
	/* eMemoryBounds */ "Bounds of provided memory were insufficient",
	/* eOptionInvalid */ "Mode or flags argument was not valid",
	/* eIndexInvalid */ "Index argument was not valid",
	/* eBaseCellDomain */ "Base cell number was outside of acceptable range",
	/* eDigitDomain */ "Child digits invalid",
	/* eDeletedDigit */ "Deleted subsequence indicates invalid index",
}
