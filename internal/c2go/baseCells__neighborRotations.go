package c2go

// baseCellNeighbor60CCWRots table contains neighboring base cell rotations in each IJK direction.
// For each base cell, for each direction, the number of 60 degree CCW rotations
// to the coordinate system of the neighbor is given.
// -1 indicates there is no neighbor in that direction.
var baseCellNeighbor60CCWRots = [NUM_BASE_CELLS][7]int{
	{0, 5, 0, 0, 1, 5, 1},   // base cell 0
	{0, 0, 1, 0, 1, 0, 1},   // base cell 1
	{0, 0, 0, 0, 0, 5, 0},   // base cell 2
	{0, 5, 0, 0, 2, 5, 1},   // base cell 3
	{0, -1, 1, 0, 3, 4, 2},  // base cell 4 (pentagon)
	{0, 0, 1, 0, 1, 0, 1},   // base cell 5
	{0, 0, 0, 3, 5, 5, 0},   // base cell 6
	{0, 0, 0, 0, 0, 5, 0},   // base cell 7
	{0, 5, 0, 0, 0, 5, 1},   // base cell 8
	{0, 0, 1, 5, 0, 0, 1},   // base cell 9
	{0, 0, 1, 0, 1, 0, 1},   // base cell 10
	{0, 3, 0, 0, 0, 0, 0},   // base cell 11
	{0, 5, 0, 0, 3, 5, 1},   // base cell 12
	{0, 0, 1, 0, 1, 0, 1},   // base cell 13
	{0, -1, 3, 0, 5, 2, 0},  // base cell 14 (pentagon)
	{0, 5, 0, 0, 4, 5, 1},   // base cell 15
	{0, 0, 0, 0, 0, 5, 0},   // base cell 16
	{0, 3, 3, 3, 0, 0, 0},   // base cell 17
	{0, 0, 0, 0, 0, 5, 0},   // base cell 18
	{0, 3, 3, 3, 0, 0, 0},   // base cell 19
	{0, 3, 3, 3, 0, 0, 3},   // base cell 20
	{0, 0, 0, 0, 0, 5, 0},   // base cell 21
	{0, 0, 1, 0, 1, 0, 1},   // base cell 22
	{0, 3, 0, 0, 0, 0, 0},   // base cell 23
	{0, -1, 3, 0, 5, 2, 0},  // base cell 24 (pentagon)
	{0, 0, 0, 0, 0, 0, 0},   // base cell 25
	{0, 0, 0, 0, 1, 0, 1},   // base cell 26
	{0, 3, 0, 0, 2, 3, 3},   // base cell 27
	{0, 0, 1, 0, 1, 0, 1},   // base cell 28
	{0, 0, 1, 3, 0, 0, 1},   // base cell 29
	{0, 3, 3, 3, 0, 0, 0},   // base cell 30
	{0, 0, 0, 0, 0, 5, 0},   // base cell 31
	{0, 3, 3, 3, 3, 0, 3},   // base cell 32
	{0, 3, 3, 3, 3, 0, 3},   // base cell 33
	{0, 3, 3, 3, 3, 3, 0},   // base cell 34
	{0, 0, 3, 0, 0, 0, 0},   // base cell 35
	{0, 0, 0, 0, 0, 0, 0},   // base cell 36
	{0, 3, 0, 0, 0, 0, 0},   // base cell 37
	{0, -1, 3, 0, 5, 2, 0},  // base cell 38 (pentagon)
	{0, 3, 0, 0, 3, 3, 0},   // base cell 39
	{0, 3, 0, 0, 3, 3, 0},   // base cell 40
	{0, 0, 0, 3, 5, 5, 0},   // base cell 41
	{0, 0, 0, 0, 0, 5, 0},   // base cell 42
	{0, 1, 3, 1, 5, 5, 0},   // base cell 43
	{0, 0, 0, 0, 0, 0, 0},   // base cell 44
	{0, 0, 0, 0, 0, 3, 0},   // base cell 45
	{0, 0, 0, 0, 3, 0, 0},   // base cell 46
	{0, 3, 3, 3, 0, 0, 0},   // base cell 47
	{0, 3, 3, 3, 0, 0, 0},   // base cell 48
	{0, -1, 3, 0, 5, 2, 0},  // base cell 49 (pentagon)
	{0, 0, 0, 3, 0, 0, 0},   // base cell 50
	{0, 3, 0, 0, 0, 0, 3},   // base cell 51
	{0, 0, 3, 0, 0, 0, 0},   // base cell 52
	{0, 0, 0, 0, 0, 5, 0},   // base cell 53
	{0, 0, 3, 0, 0, 3, 0},   // base cell 54
	{0, 0, 0, 0, 0, 0, 0},   // base cell 55
	{0, 0, 0, 0, 0, 0, 3},   // base cell 56
	{0, 3, 3, 3, 0, 0, 0},   // base cell 57
	{0, -1, 3, 0, 5, 2, 0},  // base cell 58 (pentagon)
	{0, 3, 3, 3, 3, 3, 0},   // base cell 59
	{0, 0, 0, 0, 0, 0, 0},   // base cell 60
	{0, 0, 0, 3, 0, 0, 0},   // base cell 61
	{0, 3, 3, 3, 0, 0, 0},   // base cell 62
	{0, -1, 3, 0, 5, 2, 0},  // base cell 63 (pentagon)
	{0, 0, 0, 0, 0, 0, 0},   // base cell 64
	{0, 0, 0, 0, 0, 5, 0},   // base cell 65
	{0, 3, 3, 3, 0, 3, 0},   // base cell 66
	{0, 0, 0, 0, 0, 0, 0},   // base cell 67
	{0, 3, 0, 0, 0, 0, 3},   // base cell 68
	{0, 3, 0, 0, 0, 3, 0},   // base cell 69
	{0, 3, 0, 0, 3, 3, 0},   // base cell 70
	{0, 0, 0, 0, 0, 0, 3},   // base cell 71
	{0, -1, 3, 0, 5, 2, 0},  // base cell 72 (pentagon)
	{0, 3, 3, 3, 0, 0, 3},   // base cell 73
	{0, 3, 0, 0, 0, 3, 0},   // base cell 74
	{0, 0, 0, 3, 0, 0, 0},   // base cell 75
	{0, 0, 0, 0, 0, 0, 0},   // base cell 76
	{0, 0, 3, 0, 0, 0, 0},   // base cell 77
	{0, 0, 0, 0, 0, 3, 0},   // base cell 78
	{0, 3, 3, 3, 3, 0, 3},   // base cell 79
	{0, 0, 0, 0, 0, 0, 3},   // base cell 80
	{0, 3, 3, 3, 0, 0, 0},   // base cell 81
	{0, 0, 3, 0, 0, 0, 0},   // base cell 82
	{0, -1, 3, 0, 5, 2, 0},  // base cell 83 (pentagon)
	{0, 0, 3, 0, 0, 0, 3},   // base cell 84
	{0, 3, 3, 3, 0, 0, 0},   // base cell 85
	{0, 0, 0, 3, 0, 0, 0},   // base cell 86
	{0, 3, 0, 0, 3, 0, 0},   // base cell 87
	{0, 0, 3, 0, 3, 0, 3},   // base cell 88
	{0, 0, 3, 0, 3, 0, 0},   // base cell 89
	{0, 0, 0, 3, 0, 3, 0},   // base cell 90
	{0, 0, 0, 0, 0, 0, 0},   // base cell 91
	{0, 3, 3, 3, 0, 0, 0},   // base cell 92
	{0, 0, 3, 0, 0, 3, 3},   // base cell 93
	{0, 0, 0, 0, 0, 0, 0},   // base cell 94
	{0, 0, 0, 0, 3, 0, 0},   // base cell 95
	{0, 0, 0, 0, 0, 3, 0},   // base cell 96
	{0, -1, 3, 0, 5, 2, 0},  // base cell 97 (pentagon)
	{0, 3, 3, 3, 0, 0, 0},   // base cell 98
	{0, 0, 0, 0, 3, 0, 3},   // base cell 99
	{0, 3, 0, 0, 3, 3, 0},   // base cell 100
	{0, 3, 0, 0, 0, 0, 0},   // base cell 101
	{0, 0, 3, 0, 3, 3, 0},   // base cell 102
	{0, 0, 0, 0, 3, 0, 0},   // base cell 103
	{0, 3, 3, 3, 0, 0, 0},   // base cell 104
	{0, 0, 0, 0, 0, 0, 3},   // base cell 105
	{0, 0, 3, 0, 0, 0, 3},   // base cell 106
	{0, -1, 3, 0, 5, 2, 0},  // base cell 107 (pentagon)
	{0, 3, 0, 0, 0, 3, 0},   // base cell 108
	{0, 3, 3, 3, 3, 0, 0},   // base cell 109
	{0, 0, 0, 0, 0, 0, 0},   // base cell 110
	{0, 0, 3, 0, 0, 0, 0},   // base cell 111
	{0, 0, 3, 0, 0, 3, 3},   // base cell 112
	{0, 0, 0, 0, 3, 0, 3},   // base cell 113
	{0, 0, 0, 3, 3, 0, 0},   // base cell 114
	{0, 0, 0, 0, 3, 0, 0},   // base cell 115
	{0, 0, 0, 0, 0, 3, 0},   // base cell 116
	{0, -1, 0, 0, 3, 5, 3},  // base cell 117 (pentagon)
	{0, 0, 3, 3, 3, 0, 0},   // base cell 118
	{0, 0, 0, 0, 3, 3, 0},   // base cell 119
	{0, 0, 3, 3, 3, 3, 0},   // base cell 120
	{0, 0, 0, 3, 3, 3, 0},   // base cell 121
}