package h3

// Overage represents overage status for coordinate adjustments
type Overage int

const (
	NO_OVERAGE Overage = 0 // No overage
	FACE_EDGE  Overage = 1 // On face edge (only occurs on substrate grids)
	NEW_FACE   Overage = 2 // Overage on new face interior
)

// FaceOrientIJK contains information to transform into an adjacent face IJK system
type FaceOrientIJK struct {
	Face      int32    // face number
	Translate CoordIJK // res 0 translation relative to primary face
	CcwRot60  int32    // number of 60 degree ccw rotations relative to primary face
}

// Quadrant direction constants for faceNeighbors table
const (
	IJ = 1 // IJ quadrant faceNeighbors table direction
	KI = 2 // KI quadrant faceNeighbors table direction
	JK = 3 // JK quadrant faceNeighbors table direction
)

// maxDimByCIIres provides maximum dimension value for each Class II resolution
// Ported from H3 C: faceijk.c::maxDimByCIIres
var maxDimByCIIres = [17]int32{
	2,        // res  0
	-1,       // res  1
	14,       // res  2
	-1,       // res  3
	98,       // res  4
	-1,       // res  5
	686,      // res  6
	-1,       // res  7
	4802,     // res  8
	-1,       // res  9
	33614,    // res 10
	-1,       // res 11
	235298,   // res 12
	-1,       // res 13
	1647086,  // res 14
	-1,       // res 15
	11529602, // res 16
}

// unitScaleByCIIres provides unit scale for each Class II resolution
// Ported from H3 C: faceijk.c::unitScaleByCIIres
var unitScaleByCIIres = [17]int32{
	1,       // res  0
	-1,      // res  1
	7,       // res  2
	-1,      // res  3
	49,      // res  4
	-1,      // res  5
	343,     // res  6
	-1,      // res  7
	2401,    // res  8
	-1,      // res  9
	16807,   // res 10
	-1,      // res 11
	117649,  // res 12
	-1,      // res 13
	823543,  // res 14
	-1,      // res 15
	5764801, // res 16
}

// adjacentFaceDir maps from face to face to the directions of the adjacent face
// Ported from H3 C: faceijk.c::adjacentFaceDir
var adjacentFaceDir = [20][20]int32{
	{0, KI, -1, -1, IJ, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 0
	{IJ, 0, KI, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 1
	{-1, IJ, 0, KI, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 2
	{-1, -1, IJ, 0, KI, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 3
	{KI, -1, -1, IJ, 0, -1, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 4
	{JK, -1, -1, -1, -1, 0, -1, -1, -1, -1, IJ, -1, -1, -1, KI, -1, -1, -1, -1, -1}, // face 5
	{-1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1, -1, -1, -1}, // face 6
	{-1, -1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1, -1, -1}, // face 7
	{-1, -1, -1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1, -1}, // face 8
	{-1, -1, -1, -1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1}, // face 9
	{-1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1, -1, -1, -1}, // face 10
	{-1, -1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1, -1, -1}, // face 11
	{-1, -1, -1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1, -1}, // face 12
	{-1, -1, -1, -1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1}, // face 13
	{-1, -1, -1, -1, -1, KI, -1, -1, -1, IJ, -1, -1, -1, -1, 0, -1, -1, -1, -1, JK}, // face 14
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, -1, 0, IJ, -1, -1, KI}, // face 15
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, KI, 0, IJ, -1, -1}, // face 16
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, KI, 0, IJ, -1}, // face 17
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, KI, 0, IJ}, // face 18
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, IJ, -1, -1, KI, 0}, // face 19
}
