package tables

// BaseCells contains metadata for all base cells.
// Data extracted from H3 C v4.3.0 baseCells.c.
var BaseCells = [NumBaseCells]BaseCellData{
	{Face: 1, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 0
	{Face: 2, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 1
	{Face: 1, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 2
	{Face: 2, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 3
	{Face: 0, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{-1, -1}},  // base cell 4 (pentagon)
	{Face: 1, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 5
	{Face: 1, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 6
	{Face: 2, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 7
	{Face: 0, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 8
	{Face: 2, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 9
	{Face: 1, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 10
	{Face: 1, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 11
	{Face: 3, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 12
	{Face: 3, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 13
	{Face: 11, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{2, 6}},   // base cell 14 (pentagon)
	{Face: 4, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 15
	{Face: 0, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 16
	{Face: 6, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 17
	{Face: 0, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 18
	{Face: 2, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 19
	{Face: 7, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 20
	{Face: 2, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 21
	{Face: 0, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 22
	{Face: 6, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 23
	{Face: 10, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{1, 5}},   // base cell 24 (pentagon)
	{Face: 6, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 25
	{Face: 3, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 26
	{Face: 11, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 27
	{Face: 4, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 28
	{Face: 3, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 29
	{Face: 0, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 30
	{Face: 4, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 31
	{Face: 5, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 32
	{Face: 0, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 33
	{Face: 7, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 34
	{Face: 11, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 35
	{Face: 7, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 36
	{Face: 10, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 37
	{Face: 12, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{3, 7}},   // base cell 38 (pentagon)
	{Face: 6, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 39
	{Face: 7, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 40
	{Face: 4, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 41
	{Face: 3, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 42
	{Face: 3, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 43
	{Face: 4, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 44
	{Face: 6, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 45
	{Face: 11, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 46
	{Face: 8, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 47
	{Face: 5, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 48
	{Face: 14, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{0, 9}},   // base cell 49 (pentagon)
	{Face: 5, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 50
	{Face: 12, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 51
	{Face: 10, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 52
	{Face: 4, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 53
	{Face: 12, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 54
	{Face: 7, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 55
	{Face: 11, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 56
	{Face: 10, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 57
	{Face: 13, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{4, 8}},   // base cell 58 (pentagon)
	{Face: 10, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 59
	{Face: 11, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 60
	{Face: 9, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 61
	{Face: 8, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 62
	{Face: 6, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{11, 15}},  // base cell 63 (pentagon)
	{Face: 8, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 64
	{Face: 9, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 65
	{Face: 14, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 66
	{Face: 5, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 67
	{Face: 16, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 68
	{Face: 8, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 69
	{Face: 5, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 70
	{Face: 12, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 71
	{Face: 7, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{12, 16}},  // base cell 72 (pentagon)
	{Face: 12, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 73
	{Face: 10, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 74
	{Face: 9, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 75
	{Face: 13, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 76
	{Face: 16, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 77
	{Face: 15, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 78
	{Face: 15, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 79
	{Face: 16, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 80
	{Face: 14, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 81
	{Face: 13, IJK0: [3]int{1, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 82
	{Face: 5, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{10, 19}},  // base cell 83 (pentagon)
	{Face: 8, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 84
	{Face: 14, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 85
	{Face: 9, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 86
	{Face: 14, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 87
	{Face: 17, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 88
	{Face: 12, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 89
	{Face: 16, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 90
	{Face: 17, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 91
	{Face: 15, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 92
	{Face: 16, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 93
	{Face: 9, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},    // base cell 94
	{Face: 15, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 95
	{Face: 13, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 96
	{Face: 8, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{13, 17}},  // base cell 97 (pentagon)
	{Face: 13, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 98
	{Face: 17, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 99
	{Face: 19, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 100
	{Face: 14, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 101
	{Face: 19, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 102
	{Face: 17, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 103
	{Face: 13, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 104
	{Face: 17, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 105
	{Face: 16, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 106
	{Face: 9, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{14, 18}},  // base cell 107 (pentagon)
	{Face: 15, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 108
	{Face: 15, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 109
	{Face: 18, IJK0: [3]int{0, 1, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 110
	{Face: 18, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 111
	{Face: 19, IJK0: [3]int{0, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 112
	{Face: 17, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 113
	{Face: 19, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 114
	{Face: 18, IJK0: [3]int{0, 1, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 115
	{Face: 18, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 116
	{Face: 19, IJK0: [3]int{2, 0, 0}, IsPentagon: 1, CWOffsetPent: [2]int{-1, -1}}, // base cell 117 (pentagon)
	{Face: 19, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 118
	{Face: 18, IJK0: [3]int{0, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 119
	{Face: 19, IJK0: [3]int{1, 0, 1}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 120
	{Face: 18, IJK0: [3]int{1, 0, 0}, IsPentagon: 0, CWOffsetPent: [2]int{0, 0}},   // base cell 121
}
