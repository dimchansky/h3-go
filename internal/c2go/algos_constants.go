package c2go

// Directions used for traversing a hexagonal ring counterclockwise around {1, 0, 0}
//
//	   _
//	 _/ \_
//	/ \5/ \
//	\0/ \4/
//	/ \_/ \
//	\1/ \3/
//	  \2/
//
// Ported from H3 C: algos.c::DIRECTIONS
var DIRECTIONS = [6]Direction{
	J_AXES_DIGIT,  // 0
	JK_AXES_DIGIT, // 1
	K_AXES_DIGIT,  // 2
	IK_AXES_DIGIT, // 3
	I_AXES_DIGIT,  // 4
	IJ_AXES_DIGIT, // 5
}

// Direction used for traversing to the next outward hexagonal ring.
// Ported from H3 C: algos.c::NEXT_RING_DIRECTION
const NEXT_RING_DIRECTION = I_AXES_DIGIT

// k value which will encompass all cells at resolution 15.
// This is the largest possible k in the H3 grid system.
// Ported from H3 C: algos.c::K_ALL_CELLS_AT_RES_15
const K_ALL_CELLS_AT_RES_15 = 13780510

// Buffer size for polygon to cells estimation.
// When the polygon is very small, near an icosahedron edge and is an odd
// resolution, the line tracing needs an extra buffer than the estimator
// function provides.
// Ported from H3 C: algos.c::POLYGON_TO_CELLS_BUFFER
const POLYGON_TO_CELLS_BUFFER = 12
