package h3

const (
	M_PI      = 3.14159265358979323846                       // pi
	M_PI_2    = 1.5707963267948966                           // pi / 2.0
	M_2PI     = 6.28318530717958647692528676655900576839433  // 2.0 * PI
	M_PI_180  = 0.0174532925199432957692369076848861271111   // pi / 180
	M_180_PI  = 57.29577951308232087679815481410517033240547 // pi * 180
	EPSILON   = 0.0000000000000001                           // threshold epsilon
	M_SQRT3_2 = 0.8660254037844386467637231707529361834714   // sqrt(3) / 2.0
	M_SIN60   = M_SQRT3_2                                    // sin(60')

	// M_AP7_ROT_RADS is the rotation angle between Class II and Class III resolution axes
	// (asin(sqrt(3.0 / 28.0)))
	M_AP7_ROT_RADS = 0.333473172251832115336090755351601070065900389
	M_SIN_AP7_ROT  = 0.3273268353539885718950318 // sin(M_AP7_ROT_RADS)
	M_COS_AP7_ROT  = 0.9449111825230680680167902 // cos(M_AP7_ROT_RADS)

	// EARTH_RADIUS_KM is the Earth radius in kilometers using WGS84 authalic radius
	EARTH_RADIUS_KM = 6371.007180918475

	// RES0_U_GNOMONIC is the scaling factor from hex2d resolution 0 unit length
	// (or distance between adjacent cell center points
	// on the plane) to gnomonic unit length.
	RES0_U_GNOMONIC = 0.38196601125010500003

	MAX_H3_RES = 15 // max H3 resolution; H3 version 1 has 16 resolutions, numbered 0 through 15.

	NUM_ICOSA_FACES = 20  // The number of faces on an icosahedron
	NUM_BASE_CELLS  = 122 // The number of H3 base cells
	NUM_HEX_VERTS   = 6   // The number of vertices in a hexagon
	NUM_PENT_VERTS  = 5   // The number of vertices in a pentagon
	NUM_PENTAGONS   = 12  // The number of pentagons per resolution

	// H3 index modes

	H3_HEXAGON_MODE = 1
	H3_UNIEDGE_MODE = 2
)
