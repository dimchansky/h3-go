package h3

// Local H3 constants mirrored from testref h3lib/include/constants.h
// Keep values in sync with the referenced H3 version used in tests.

const (
	// Mathematical constants.
	mPi              = 3.14159265358979323846
	mPi2             = 1.5707963267948966
	m2pi             = 6.28318530717958647692528676655900576839433
	mPi180           = 0.0174532925199432957692369076848861271111
	m180Pi           = 57.29577951308232087679815481410517033240547
	epsilon          = 0.0000000000000001
	mSqrt32          = 0.8660254037844386467637231707529361834714
	mSin60           = mSqrt32
	mRsin60          = 1.1547005383792515290182975610039149112953
	mOnethird        = 0.333333333333333333333333333333333333333
	mOneseventh      = 0.14285714285714285714285714285714285
	mAp7RotRads      = 0.333473172251832115336090755351601070065900389
	mSinAp7Rot       = 0.3273268353539885718950318
	mCosAp7Rot       = 0.9449111825230680680167902
	mSqrt7           = 2.6457513110645905905016157536392604257102
	mRsqrt7          = 0.37796447300922722721451653623418006081576
	invRes0UGnomonic = 2.61803398874989588842
	res0UGnomonic    = 0.38196601125010500003

	// Earth radius (km).
	earthRadiusKm = 6371.007180918475

	// Epsilon constants from latLng.h.
	epsilonDeg = .000000001          // epsilon of ~0.1mm in degrees
	epsilonRad = epsilonDeg * mPi180 // epsilon of ~0.1mm in radians

	// Resolution and topology.
	maxH3Res      = 15
	numIcosaFaces = 20
	numBaseCells  = 122
	numHexVerts   = 6
	numPentVerts  = 5
	numPentagons  = 12

	// Base cell constants.
	invalidBaseCell  = 127
	invalidRotations = -1
	maxFaceCoord     = 2

	// Face constants.
	invalidFace = -1

	// Modes.
	h3CellMode         = 1
	h3DirectededgeMode = 2
	h3EdgeMode         = 3
	h3VertexMode       = 4
)

// Integer limits for overflow checking (from stdint.h).
const (
	int32Max  = 2147483647
	int32Min  = -2147483648
	int32Max3 = int32Max / 3 // Used in aperture 7 overflow checking
)
