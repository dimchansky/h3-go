package h3

// baseCellDataEntry represents information on a single base cell.
// Mirrors the baseCellDataEntry struct from baseCells.h.
type baseCellDataEntry struct {
	HomeFijk     faceIJK  // "home" face and normalized ijk coordinates on that face
	IsPentagon   bool     // is this base cell a pentagon?
	CwOffsetPent [2]int32 // if a pentagon, what are its two clockwise offset faces?
}

// baseCellData mirrors the baseCellDataEntry array from baseCells.c
// Resolution 0 base cell data lookup table
// Ported from H3 C: baseCells.c::baseCellData.
var baseCellData = [numBaseCells]baseCellDataEntry{
	{faceIJK{1, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 0
	{faceIJK{2, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 1
	{faceIJK{1, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 2
	{faceIJK{2, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 3
	{faceIJK{0, coordIJK{2, 0, 0}}, true, [2]int32{-1, -1}},  // base cell 4
	{faceIJK{1, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 5
	{faceIJK{1, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 6
	{faceIJK{2, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 7
	{faceIJK{0, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 8
	{faceIJK{2, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 9
	{faceIJK{1, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 10
	{faceIJK{1, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 11
	{faceIJK{3, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 12
	{faceIJK{3, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 13
	{faceIJK{11, coordIJK{2, 0, 0}}, true, [2]int32{2, 6}},   // base cell 14
	{faceIJK{4, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 15
	{faceIJK{0, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 16
	{faceIJK{6, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 17
	{faceIJK{0, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 18
	{faceIJK{2, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 19
	{faceIJK{7, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 20
	{faceIJK{2, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 21
	{faceIJK{0, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 22
	{faceIJK{6, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 23
	{faceIJK{10, coordIJK{2, 0, 0}}, true, [2]int32{1, 5}},   // base cell 24
	{faceIJK{6, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 25
	{faceIJK{3, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 26
	{faceIJK{11, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 27
	{faceIJK{4, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},   // base cell 28
	{faceIJK{3, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 29
	{faceIJK{0, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 30
	{faceIJK{4, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 31
	{faceIJK{5, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 32
	{faceIJK{0, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 33
	{faceIJK{7, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 34
	{faceIJK{11, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 35
	{faceIJK{7, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 36
	{faceIJK{10, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 37
	{faceIJK{12, coordIJK{2, 0, 0}}, true, [2]int32{3, 7}},   // base cell 38
	{faceIJK{6, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 39
	{faceIJK{7, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 40
	{faceIJK{4, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 41
	{faceIJK{3, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 42
	{faceIJK{3, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 43
	{faceIJK{4, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 44
	{faceIJK{6, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 45
	{faceIJK{11, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 46
	{faceIJK{8, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 47
	{faceIJK{5, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 48
	{faceIJK{14, coordIJK{2, 0, 0}}, true, [2]int32{0, 9}},   // base cell 49
	{faceIJK{5, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 50
	{faceIJK{12, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 51
	{faceIJK{10, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 52
	{faceIJK{4, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},   // base cell 53
	{faceIJK{12, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 54
	{faceIJK{7, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 55
	{faceIJK{11, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 56
	{faceIJK{10, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 57
	{faceIJK{13, coordIJK{2, 0, 0}}, true, [2]int32{4, 8}},   // base cell 58
	{faceIJK{10, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 59
	{faceIJK{11, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 60
	{faceIJK{9, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 61
	{faceIJK{8, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},   // base cell 62
	{faceIJK{6, coordIJK{2, 0, 0}}, true, [2]int32{11, 15}},  // base cell 63
	{faceIJK{8, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 64
	{faceIJK{9, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},   // base cell 65
	{faceIJK{14, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 66
	{faceIJK{5, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 67
	{faceIJK{16, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 68
	{faceIJK{8, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 69
	{faceIJK{5, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 70
	{faceIJK{12, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 71
	{faceIJK{7, coordIJK{2, 0, 0}}, true, [2]int32{12, 16}},  // base cell 72
	{faceIJK{12, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 73
	{faceIJK{10, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 74
	{faceIJK{9, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},   // base cell 75
	{faceIJK{13, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 76
	{faceIJK{16, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 77
	{faceIJK{15, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 78
	{faceIJK{15, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 79
	{faceIJK{16, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 80
	{faceIJK{14, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 81
	{faceIJK{13, coordIJK{1, 1, 0}}, false, [2]int32{0, 0}},  // base cell 82
	{faceIJK{5, coordIJK{2, 0, 0}}, true, [2]int32{10, 19}},  // base cell 83
	{faceIJK{8, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 84
	{faceIJK{14, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 85
	{faceIJK{9, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},   // base cell 86
	{faceIJK{14, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 87
	{faceIJK{17, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 88
	{faceIJK{12, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 89
	{faceIJK{16, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 90
	{faceIJK{17, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 91
	{faceIJK{15, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 92
	{faceIJK{16, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 93
	{faceIJK{9, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},   // base cell 94
	{faceIJK{15, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 95
	{faceIJK{13, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 96
	{faceIJK{8, coordIJK{2, 0, 0}}, true, [2]int32{13, 17}},  // base cell 97
	{faceIJK{13, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 98
	{faceIJK{17, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 99
	{faceIJK{19, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 100
	{faceIJK{14, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 101
	{faceIJK{19, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 102
	{faceIJK{17, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 103
	{faceIJK{13, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 104
	{faceIJK{17, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 105
	{faceIJK{16, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 106
	{faceIJK{9, coordIJK{2, 0, 0}}, true, [2]int32{14, 18}},  // base cell 107
	{faceIJK{15, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 108
	{faceIJK{15, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 109
	{faceIJK{18, coordIJK{0, 1, 1}}, false, [2]int32{0, 0}},  // base cell 110
	{faceIJK{18, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 111
	{faceIJK{19, coordIJK{0, 0, 1}}, false, [2]int32{0, 0}},  // base cell 112
	{faceIJK{17, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 113
	{faceIJK{19, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 114
	{faceIJK{18, coordIJK{0, 1, 0}}, false, [2]int32{0, 0}},  // base cell 115
	{faceIJK{18, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 116
	{faceIJK{19, coordIJK{2, 0, 0}}, true, [2]int32{-1, -1}}, // base cell 117
	{faceIJK{19, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 118
	{faceIJK{18, coordIJK{0, 0, 0}}, false, [2]int32{0, 0}},  // base cell 119
	{faceIJK{19, coordIJK{1, 0, 1}}, false, [2]int32{0, 0}},  // base cell 120
	{faceIJK{18, coordIJK{1, 0, 0}}, false, [2]int32{0, 0}},  // base cell 121
}
