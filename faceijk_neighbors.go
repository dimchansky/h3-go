package h3

// faceNeighbors provides face transition information for coordinate system transformations
// Ported from H3 C: faceijk.c::faceNeighbors.
var faceNeighbors = [numIcosaFaces][4]faceOrientIJK{
	{
		// face 0
		{0, coordIJK{0, 0, 0}, 0}, // central face
		{4, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{1, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{5, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 1
		{1, coordIJK{0, 0, 0}, 0}, // central face
		{0, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{2, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{6, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 2
		{2, coordIJK{0, 0, 0}, 0}, // central face
		{1, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{3, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{7, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 3
		{3, coordIJK{0, 0, 0}, 0}, // central face
		{2, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{4, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{8, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 4
		{4, coordIJK{0, 0, 0}, 0}, // central face
		{3, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{0, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{9, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 5
		{5, coordIJK{0, 0, 0}, 0},  // central face
		{10, coordIJK{2, 2, 0}, 3}, // ij quadrant
		{14, coordIJK{2, 0, 2}, 3}, // ki quadrant
		{0, coordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 6
		{6, coordIJK{0, 0, 0}, 0},  // central face
		{11, coordIJK{2, 2, 0}, 3}, // ij quadrant
		{10, coordIJK{2, 0, 2}, 3}, // ki quadrant
		{1, coordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 7
		{7, coordIJK{0, 0, 0}, 0},  // central face
		{12, coordIJK{2, 2, 0}, 3}, // ij quadrant
		{11, coordIJK{2, 0, 2}, 3}, // ki quadrant
		{2, coordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 8
		{8, coordIJK{0, 0, 0}, 0},  // central face
		{13, coordIJK{2, 2, 0}, 3}, // ij quadrant
		{12, coordIJK{2, 0, 2}, 3}, // ki quadrant
		{3, coordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 9
		{9, coordIJK{0, 0, 0}, 0},  // central face
		{14, coordIJK{2, 2, 0}, 3}, // ij quadrant
		{13, coordIJK{2, 0, 2}, 3}, // ki quadrant
		{4, coordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 10
		{10, coordIJK{0, 0, 0}, 0}, // central face
		{5, coordIJK{2, 2, 0}, 3},  // ij quadrant
		{6, coordIJK{2, 0, 2}, 3},  // ki quadrant
		{15, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 11
		{11, coordIJK{0, 0, 0}, 0}, // central face
		{6, coordIJK{2, 2, 0}, 3},  // ij quadrant
		{7, coordIJK{2, 0, 2}, 3},  // ki quadrant
		{16, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 12
		{12, coordIJK{0, 0, 0}, 0}, // central face
		{7, coordIJK{2, 2, 0}, 3},  // ij quadrant
		{8, coordIJK{2, 0, 2}, 3},  // ki quadrant
		{17, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 13
		{13, coordIJK{0, 0, 0}, 0}, // central face
		{8, coordIJK{2, 2, 0}, 3},  // ij quadrant
		{9, coordIJK{2, 0, 2}, 3},  // ki quadrant
		{18, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 14
		{14, coordIJK{0, 0, 0}, 0}, // central face
		{9, coordIJK{2, 2, 0}, 3},  // ij quadrant
		{5, coordIJK{2, 0, 2}, 3},  // ki quadrant
		{19, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 15
		{15, coordIJK{0, 0, 0}, 0}, // central face
		{16, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{19, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{10, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 16
		{16, coordIJK{0, 0, 0}, 0}, // central face
		{17, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{15, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{11, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 17
		{17, coordIJK{0, 0, 0}, 0}, // central face
		{18, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{16, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{12, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 18
		{18, coordIJK{0, 0, 0}, 0}, // central face
		{19, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{17, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{13, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 19
		{19, coordIJK{0, 0, 0}, 0}, // central face
		{15, coordIJK{2, 0, 2}, 1}, // ij quadrant
		{18, coordIJK{2, 2, 0}, 5}, // ki quadrant
		{14, coordIJK{0, 2, 2}, 3}, // jk quadrant
	},
}
