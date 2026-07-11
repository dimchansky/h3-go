package h3

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
// Ported from H3 C: algos.c::DIRECTIONS.
var algosDirections = [6]direction{
	jAxesDigit,  // 0
	jkAxesDigit, // 1
	kAxesDigit,  // 2
	ikAxesDigit, // 3
	iAxesDigit,  // 4
	ijAxesDigit, // 5
}

// direction used for traversing to the next outward hexagonal ring.
// Ported from H3 C: algos.c::NEXT_RING_DIRECTION.
const nextRingDirection = iAxesDigit

// k value which will encompass all cells at resolution 15.
// This is the largest possible k in the H3 grid system.
// Ported from H3 C: algos.c::K_ALL_CELLS_AT_RES_15.
const kAllCellsAtRes15 = 13780510

// Maximum number of cells in a ring (including center cell).
// This is 7 cells: the center cell plus up to 6 neighbors.
// Ported from H3 C: algos.c::MAX_ONE_RING_SIZE.
const maxOneRingSize = 7

// Buffer size for polygon to cells estimation.
// When the polygon is very small, near an icosahedron edge and is an odd
// resolution, the line tracing needs an extra buffer than the estimator
// function provides.
// Ported from H3 C: algos.c::POLYGON_TO_CELLS_BUFFER.
const polygonToCellsBuffer = 12
