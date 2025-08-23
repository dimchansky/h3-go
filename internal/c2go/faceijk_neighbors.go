package c2go

// faceNeighbors provides face transition information for coordinate system transformations
// Ported from H3 C: faceijk.c::faceNeighbors
var faceNeighbors = [NUM_ICOSA_FACES][4]FaceOrientIJK{
	{
		// face 0
		{0, CoordIJK{0, 0, 0}, 0}, // central face
		{4, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{1, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{5, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 1
		{1, CoordIJK{0, 0, 0}, 0}, // central face
		{0, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{2, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{6, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 2
		{2, CoordIJK{0, 0, 0}, 0}, // central face
		{1, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{3, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{7, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 3
		{3, CoordIJK{0, 0, 0}, 0}, // central face
		{2, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{4, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{8, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 4
		{4, CoordIJK{0, 0, 0}, 0}, // central face
		{3, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{0, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{9, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 5
		{5, CoordIJK{0, 0, 0}, 0},  // central face
		{10, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{14, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{0, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 6
		{6, CoordIJK{0, 0, 0}, 0},  // central face
		{11, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{10, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{1, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 7
		{7, CoordIJK{0, 0, 0}, 0},  // central face
		{12, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{11, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{2, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 8
		{8, CoordIJK{0, 0, 0}, 0},  // central face
		{13, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{12, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{3, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 9
		{9, CoordIJK{0, 0, 0}, 0},  // central face
		{14, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{13, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{4, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 10
		{10, CoordIJK{0, 0, 0}, 0}, // central face
		{5, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{6, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{15, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 11
		{11, CoordIJK{0, 0, 0}, 0}, // central face
		{6, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{7, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{16, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 12
		{12, CoordIJK{0, 0, 0}, 0}, // central face
		{7, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{8, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{17, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 13
		{13, CoordIJK{0, 0, 0}, 0}, // central face
		{8, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{9, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{18, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 14
		{14, CoordIJK{0, 0, 0}, 0}, // central face
		{9, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{5, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{19, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 15
		{15, CoordIJK{0, 0, 0}, 0}, // central face
		{16, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{19, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{10, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 16
		{16, CoordIJK{0, 0, 0}, 0}, // central face
		{17, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{15, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{11, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 17
		{17, CoordIJK{0, 0, 0}, 0}, // central face
		{18, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{16, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{12, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 18
		{18, CoordIJK{0, 0, 0}, 0}, // central face
		{19, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{17, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{13, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 19
		{19, CoordIJK{0, 0, 0}, 0}, // central face
		{15, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{18, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{14, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
}
