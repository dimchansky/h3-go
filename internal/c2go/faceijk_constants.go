package c2go

// Overage represents overage status for coordinate adjustments
type Overage int

const (
	NO_OVERAGE Overage = 0 // No overage
	FACE_EDGE  Overage = 1 // On face edge (only occurs on substrate grids)
	NEW_FACE   Overage = 2 // Overage on new face interior
)

// FaceOrientIJK contains information to transform into an adjacent face IJK system
type FaceOrientIJK struct {
	Face      int      // face number
	Translate CoordIJK // res 0 translation relative to primary face
	CcwRot60  int      // number of 60 degree ccw rotations relative to primary face
}

// Quadrant direction constants for faceNeighbors table
const (
	IJ = 1 // IJ quadrant faceNeighbors table direction
	KI = 2 // KI quadrant faceNeighbors table direction
	JK = 3 // JK quadrant faceNeighbors table direction
)

// maxDimByCIIres provides maximum dimension value for each Class II resolution
// Ported from H3 C: faceijk.c::maxDimByCIIres
var maxDimByCIIres = [17]int{
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
var unitScaleByCIIres = [17]int{
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
