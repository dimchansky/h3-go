package c2go

// Local H3 constants mirrored from testref h3lib/include/constants.h
// Keep values in sync with the referenced H3 version used in tests.

const (
    // Mathematical constants
    M_PI        = 3.14159265358979323846
    M_PI_2      = 1.5707963267948966
    M_2PI       = 6.28318530717958647692528676655900576839433
    M_PI_180    = 0.0174532925199432957692369076848861271111
    M_180_PI    = 57.29577951308232087679815481410517033240547
    EPSILON     = 0.0000000000000001
    M_SQRT3_2   = 0.8660254037844386467637231707529361834714
    M_SIN60     = M_SQRT3_2
    M_RSIN60    = 1.1547005383792515290182975610039149112953
    M_ONETHIRD  = 0.333333333333333333333333333333333333333
    M_ONESEVENTH = 0.14285714285714285714285714285714285
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
    H3_CELL_MODE        = 1
    H3_DIRECTEDEDGE_MODE = 2
    H3_EDGE_MODE         = 3
    H3_VERTEX_MODE       = 4
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

func isBaseCellPentagon(base int) bool {
    if base < 0 || base >= len(isBaseCellPentagonArr) {
        return false
    }
    return isBaseCellPentagonArr[base]
}
