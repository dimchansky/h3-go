package h3

// BaseCellData represents information on a single base cell.
// Mirrors the BaseCellData struct from baseCells.h
type BaseCellData struct {
	HomeFijk     FaceIJK  // "home" face and normalized ijk coordinates on that face
	IsPentagon   bool     // is this base cell a pentagon?
	CwOffsetPent [2]int32 // if a pentagon, what are its two clockwise offset faces?
}

// baseCellData mirrors the BaseCellData array from baseCells.c
// Resolution 0 base cell data lookup table
// Ported from H3 C: baseCells.c::baseCellData
var baseCellData = [NUM_BASE_CELLS]BaseCellData{
	{FaceIJK{1, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 0
	{FaceIJK{2, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 1
	{FaceIJK{1, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 2
	{FaceIJK{2, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 3
	{FaceIJK{0, CoordIJK{2, 0, 0}}, true, [2]int32{-1, -1}},  // base cell 4
	{FaceIJK{1, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 5
	{FaceIJK{1, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 6
	{FaceIJK{2, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 7
	{FaceIJK{0, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 8
	{FaceIJK{2, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 9
	{FaceIJK{1, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 10
	{FaceIJK{1, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 11
	{FaceIJK{3, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 12
	{FaceIJK{3, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 13
	{FaceIJK{11, CoordIJK{2, 0, 0}}, true, [2]int32{2, 6}},   // base cell 14
	{FaceIJK{4, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 15
	{FaceIJK{0, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 16
	{FaceIJK{6, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 17
	{FaceIJK{0, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 18
	{FaceIJK{2, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 19
	{FaceIJK{7, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 20
	{FaceIJK{2, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 21
	{FaceIJK{0, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 22
	{FaceIJK{6, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 23
	{FaceIJK{10, CoordIJK{2, 0, 0}}, true, [2]int32{1, 5}},   // base cell 24
	{FaceIJK{6, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 25
	{FaceIJK{3, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 26
	{FaceIJK{11, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 27
	{FaceIJK{4, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 28
	{FaceIJK{3, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 29
	{FaceIJK{0, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 30
	{FaceIJK{4, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 31
	{FaceIJK{5, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 32
	{FaceIJK{0, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 33
	{FaceIJK{7, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 34
	{FaceIJK{11, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 35
	{FaceIJK{7, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 36
	{FaceIJK{10, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 37
	{FaceIJK{12, CoordIJK{2, 0, 0}}, true, [2]int32{3, 7}},   // base cell 38
	{FaceIJK{6, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 39
	{FaceIJK{7, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 40
	{FaceIJK{4, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 41
	{FaceIJK{3, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 42
	{FaceIJK{3, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 43
	{FaceIJK{4, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 44
	{FaceIJK{6, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 45
	{FaceIJK{11, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 46
	{FaceIJK{8, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 47
	{FaceIJK{5, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 48
	{FaceIJK{14, CoordIJK{2, 0, 0}}, true, [2]int32{0, 9}},   // base cell 49
	{FaceIJK{5, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 50
	{FaceIJK{12, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 51
	{FaceIJK{10, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 52
	{FaceIJK{4, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 53
	{FaceIJK{12, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 54
	{FaceIJK{7, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 55
	{FaceIJK{11, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 56
	{FaceIJK{10, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 57
	{FaceIJK{13, CoordIJK{2, 0, 0}}, true, [2]int32{4, 8}},   // base cell 58
	{FaceIJK{10, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 59
	{FaceIJK{11, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 60
	{FaceIJK{9, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 61
	{FaceIJK{8, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 62
	{FaceIJK{6, CoordIJK{2, 0, 0}}, true, [2]int32{11, 15}},  // base cell 63
	{FaceIJK{8, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 64
	{FaceIJK{9, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 65
	{FaceIJK{14, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 66
	{FaceIJK{5, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 67
	{FaceIJK{16, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 68
	{FaceIJK{8, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 69
	{FaceIJK{5, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 70
	{FaceIJK{12, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 71
	{FaceIJK{7, CoordIJK{2, 0, 0}}, true, [2]int32{12, 16}},  // base cell 72
	{FaceIJK{12, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 73
	{FaceIJK{10, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 74
	{FaceIJK{9, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 75
	{FaceIJK{13, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 76
	{FaceIJK{16, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 77
	{FaceIJK{15, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 78
	{FaceIJK{15, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 79
	{FaceIJK{16, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 80
	{FaceIJK{14, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 81
	{FaceIJK{13, CoordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 82
	{FaceIJK{5, CoordIJK{2, 0, 0}}, true, [2]int32{10, 19}},  // base cell 83
	{FaceIJK{8, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 84
	{FaceIJK{14, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 85
	{FaceIJK{9, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 86
	{FaceIJK{14, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 87
	{FaceIJK{17, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 88
	{FaceIJK{12, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 89
	{FaceIJK{16, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 90
	{FaceIJK{17, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 91
	{FaceIJK{15, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 92
	{FaceIJK{16, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 93
	{FaceIJK{9, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 94
	{FaceIJK{15, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 95
	{FaceIJK{13, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 96
	{FaceIJK{8, CoordIJK{2, 0, 0}}, true, [2]int32{13, 17}},  // base cell 97
	{FaceIJK{13, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 98
	{FaceIJK{17, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 99
	{FaceIJK{19, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 100
	{FaceIJK{14, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 101
	{FaceIJK{19, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 102
	{FaceIJK{17, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 103
	{FaceIJK{13, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 104
	{FaceIJK{17, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 105
	{FaceIJK{16, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 106
	{FaceIJK{9, CoordIJK{2, 0, 0}}, true, [2]int32{14, 18}},  // base cell 107
	{FaceIJK{15, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 108
	{FaceIJK{15, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 109
	{FaceIJK{18, CoordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 110
	{FaceIJK{18, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 111
	{FaceIJK{19, CoordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 112
	{FaceIJK{17, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 113
	{FaceIJK{19, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 114
	{FaceIJK{18, CoordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 115
	{FaceIJK{18, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 116
	{FaceIJK{19, CoordIJK{2, 0, 0}}, true, [2]int32{-1, -1}}, // base cell 117
	{FaceIJK{19, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 118
	{FaceIJK{18, CoordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 119
	{FaceIJK{19, CoordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 120
	{FaceIJK{18, CoordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 121
}
