package c2go

// Local H3 constants mirrored from testref h3lib/include/constants.h
// Keep values in sync with the referenced H3 version used in tests.

const (
	// Mathematical constants
	M_PI           = 3.14159265358979323846
	M_PI_2         = 1.5707963267948966
	M_2PI          = 6.28318530717958647692528676655900576839433
	M_PI_180       = 0.0174532925199432957692369076848861271111
	M_180_PI       = 57.29577951308232087679815481410517033240547
	EPSILON        = 0.0000000000000001
	M_SQRT3_2      = 0.8660254037844386467637231707529361834714
	M_SIN60        = M_SQRT3_2
	M_RSIN60       = 1.1547005383792515290182975610039149112953
	M_ONETHIRD     = 0.333333333333333333333333333333333333333
	M_ONESEVENTH   = 0.14285714285714285714285714285714285
	M_AP7_ROT_RADS = 0.333473172251832115336090755351601070065900389
	M_SIN_AP7_ROT  = 0.3273268353539885718950318
	M_COS_AP7_ROT  = 0.9449111825230680680167902

	// Earth radius (km)
	EARTH_RADIUS_KM = 6371.007180918475

	// Resolution and topology
	MAX_H3_RES      = 15
	NUM_ICOSA_FACES = 20
	NUM_BASE_CELLS  = 122
	NUM_HEX_VERTS   = 6
	NUM_PENT_VERTS  = 5
	NUM_PENTAGONS   = 12

	// Modes
	H3_CELL_MODE         = 1
	H3_DIRECTEDEDGE_MODE = 2
	H3_EDGE_MODE         = 3
	H3_VERTEX_MODE       = 4
)

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
)

// Direction enum (from coordijk.h)
const (
	CENTER_DIGIT           = 0
	K_AXES_DIGIT           = 1
	J_AXES_DIGIT           = 2
	JK_AXES_DIGIT          = J_AXES_DIGIT | K_AXES_DIGIT // 3
	I_AXES_DIGIT           = 4
	IK_AXES_DIGIT          = I_AXES_DIGIT | K_AXES_DIGIT // 5
	IJ_AXES_DIGIT          = I_AXES_DIGIT | J_AXES_DIGIT // 6
	INVALID_DIGIT          = 7
	NUM_DIGITS             = INVALID_DIGIT
	PENTAGON_SKIPPED_DIGIT = K_AXES_DIGIT
)

// UNIT_VECS mirrors coordijk.h::UNIT_VECS - CoordIJK unit vectors corresponding to the 7 H3 digits.
var UNIT_VECS = [NUM_DIGITS]CoordIJK{
	{0, 0, 0}, // direction 0 (CENTER_DIGIT)
	{0, 0, 1}, // direction 1 (K_AXES_DIGIT)
	{0, 1, 0}, // direction 2 (J_AXES_DIGIT)
	{0, 1, 1}, // direction 3 (JK_AXES_DIGIT)
	{1, 0, 0}, // direction 4 (I_AXES_DIGIT)
	{1, 0, 1}, // direction 5 (IK_AXES_DIGIT)
	{1, 1, 0}, // direction 6 (IJ_AXES_DIGIT)
}

// H3Error represents H3 operation result codes.
type H3Error uint32

// H3Error codes (from h3api.h)
const (
	E_SUCCESS            H3Error = 0
	E_FAILED             H3Error = 1
	E_DOMAIN             H3Error = 2
	E_LATLNG_DOMAIN      H3Error = 3
	E_RES_DOMAIN         H3Error = 4
	E_CELL_INVALID       H3Error = 5
	E_DIR_EDGE_INVALID   H3Error = 6
	E_UNDIR_EDGE_INVALID H3Error = 7
	E_VERTEX_INVALID     H3Error = 8
	E_PENTAGON           H3Error = 9
	E_DUPLICATE_INPUT    H3Error = 10
	E_NOT_NEIGHBORS      H3Error = 11
	E_RES_MISMATCH       H3Error = 12
	E_MEMORY_ALLOC       H3Error = 13
	E_MEMORY_BOUNDS      H3Error = 14
	E_OPTION_INVALID     H3Error = 15
)

// Integer limits for overflow checking (from stdint.h)
const (
	INT32_MAX = 2147483647
	INT32_MIN = -2147483648
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
