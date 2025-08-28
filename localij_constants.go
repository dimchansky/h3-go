package h3

// Pentagon rotation lookup tables for local IJ coordinate transformations.
// These tables handle pentagon distortion when converting between H3 cells and local IJK coordinates.

// PENTAGON_ROTATIONS provides rotation amounts (60 degrees clockwise steps) based on
// the origin cell's leading digit and the target direction.
// Origin leading digit -> index leading digit -> rotations 60 cw
// Either being 1 (K axis) is invalid.
// No good default at 0.
// Ported from H3 C: localij.c::PENTAGON_ROTATIONS
var PENTAGON_ROTATIONS = [7][7]int32{
	{0, -1, 0, 0, 0, 0, 0},       // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, -1, 0, 0, 0, 1, 0},       // 2
	{0, -1, 0, 0, 1, 1, 0},       // 3
	{0, -1, 0, 5, 0, 0, 0},       // 4
	{0, -1, 5, 5, 0, 0, 0},       // 5
	{0, -1, 0, 0, 0, 0, 0},       // 6
}

// PENTAGON_ROTATIONS_REVERSE provides counter-clockwise rotation amounts for
// reversing the rotation introduced in PENTAGON_ROTATIONS when the origin is on a pentagon
// (regardless of the base cell of the index).
// Reverse base cell direction -> leading index digit -> rotations 60 ccw
// Ported from H3 C: localij.c::PENTAGON_ROTATIONS_REVERSE
var PENTAGON_ROTATIONS_REVERSE = [7][7]int32{
	{0, 0, 0, 0, 0, 0, 0},        // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, 1, 0, 0, 0, 0, 0},        // 2
	{0, 1, 0, 0, 0, 1, 0},        // 3
	{0, 5, 0, 0, 0, 0, 0},        // 4
	{0, 5, 0, 5, 0, 0, 0},        // 5
	{0, 0, 0, 0, 0, 0, 0},        // 6
}

// PENTAGON_ROTATIONS_REVERSE_NONPOLAR provides counter-clockwise rotation amounts for
// reversing the rotation introduced in PENTAGON_ROTATIONS when the index is on a
// non-polar pentagon and the origin is not.
// Reverse base cell direction -> leading index digit -> rotations 60 ccw
// Ported from H3 C: localij.c::PENTAGON_ROTATIONS_REVERSE_NONPOLAR
var PENTAGON_ROTATIONS_REVERSE_NONPOLAR = [7][7]int32{
	{0, 0, 0, 0, 0, 0, 0},        // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, 1, 0, 0, 0, 0, 0},        // 2
	{0, 1, 0, 0, 0, 1, 0},        // 3
	{0, 5, 0, 0, 0, 0, 0},        // 4
	{0, 1, 0, 5, 1, 1, 0},        // 5
	{0, 0, 0, 0, 0, 0, 0},        // 6
}

// PENTAGON_ROTATIONS_REVERSE_POLAR provides counter-clockwise rotation amounts for
// reversing the rotation introduced in PENTAGON_ROTATIONS when the index is on a
// polar pentagon and the origin is not.
// Reverse base cell direction -> leading index digit -> rotations 60 ccw
// Ported from H3 C: localij.c::PENTAGON_ROTATIONS_REVERSE_POLAR
var PENTAGON_ROTATIONS_REVERSE_POLAR = [7][7]int32{
	{0, 0, 0, 0, 0, 0, 0},        // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, 1, 1, 1, 1, 1, 1},        // 2
	{0, 1, 0, 0, 0, 1, 0},        // 3
	{0, 1, 0, 0, 1, 1, 1},        // 4
	{0, 1, 0, 5, 1, 1, 0},        // 5
	{0, 1, 1, 0, 1, 1, 1},        // 6
}

// FAILED_DIRECTIONS determines prohibited directions when unfolding a pentagon.
// Indexes by two directions, both relative to the pentagon base cell. The first
// is the direction of the origin index and the second is the direction of the
// index to unfold. Direction refers to the direction from base cell to base
// cell if the indexes are on different base cells, or the leading digit if
// within the pentagon base cell.
//
// This logic prevents unfolding across more than one icosahedron face.
// Ported from H3 C: localij.c::FAILED_DIRECTIONS
var FAILED_DIRECTIONS = [7][7]bool{
	{false, false, false, false, false, false, false}, // 0
	{false, false, false, false, false, false, false}, // 1
	{false, false, false, false, true, true, false},   // 2
	{false, false, false, false, true, false, true},   // 3
	{false, false, true, true, false, false, false},   // 4
	{false, false, true, false, false, false, true},   // 5
	{false, false, false, true, false, true, false},   // 6
}
