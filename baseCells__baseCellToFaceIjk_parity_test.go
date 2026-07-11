//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_baseCellToFaceIjk_parity(t *testing.T) {
	tests := []struct {
		name     string
		baseCell int32
	}{
		// Test all base cells 0-121
		{"base_cell_0", 0},
		{"base_cell_1", 1},
		{"base_cell_2", 2},
		{"base_cell_3", 3},
		{"base_cell_4", 4}, // pentagon
		{"base_cell_5", 5},
		{"base_cell_6", 6},
		{"base_cell_7", 7},
		{"base_cell_8", 8},
		{"base_cell_9", 9},
		{"base_cell_10", 10},
		{"base_cell_11", 11},
		{"base_cell_12", 12},
		{"base_cell_13", 13},
		{"base_cell_14", 14}, // pentagon
		{"base_cell_15", 15},
		{"base_cell_16", 16},
		{"base_cell_17", 17},
		{"base_cell_18", 18},
		{"base_cell_19", 19},
		{"base_cell_20", 20},
		{"base_cell_21", 21},
		{"base_cell_22", 22},
		{"base_cell_23", 23},
		{"base_cell_24", 24}, // pentagon
		{"base_cell_25", 25},
		{"base_cell_26", 26},
		{"base_cell_27", 27},
		{"base_cell_28", 28},
		{"base_cell_29", 29},
		{"base_cell_30", 30},
		{"base_cell_31", 31},
		{"base_cell_32", 32},
		{"base_cell_33", 33},
		{"base_cell_34", 34},
		{"base_cell_35", 35},
		{"base_cell_36", 36},
		{"base_cell_37", 37},
		{"base_cell_38", 38}, // pentagon
		{"base_cell_39", 39},
		{"base_cell_40", 40},
		{"base_cell_41", 41},
		{"base_cell_42", 42},
		{"base_cell_43", 43},
		{"base_cell_44", 44},
		{"base_cell_45", 45},
		{"base_cell_46", 46},
		{"base_cell_47", 47},
		{"base_cell_48", 48},
		{"base_cell_49", 49}, // pentagon
		{"base_cell_50", 50},
		{"base_cell_51", 51},
		{"base_cell_52", 52},
		{"base_cell_53", 53},
		{"base_cell_54", 54},
		{"base_cell_55", 55},
		{"base_cell_56", 56},
		{"base_cell_57", 57},
		{"base_cell_58", 58}, // pentagon
		{"base_cell_59", 59},
		{"base_cell_60", 60},
		{"base_cell_61", 61},
		{"base_cell_62", 62},
		{"base_cell_63", 63}, // pentagon
		{"base_cell_64", 64},
		{"base_cell_65", 65},
		{"base_cell_66", 66},
		{"base_cell_67", 67},
		{"base_cell_68", 68},
		{"base_cell_69", 69},
		{"base_cell_70", 70},
		{"base_cell_71", 71},
		{"base_cell_72", 72}, // pentagon
		{"base_cell_73", 73},
		{"base_cell_74", 74},
		{"base_cell_75", 75},
		{"base_cell_76", 76},
		{"base_cell_77", 77},
		{"base_cell_78", 78},
		{"base_cell_79", 79},
		{"base_cell_80", 80},
		{"base_cell_81", 81},
		{"base_cell_82", 82}, // pentagon
		{"base_cell_83", 83}, // pentagon
		{"base_cell_84", 84},
		{"base_cell_85", 85},
		{"base_cell_86", 86},
		{"base_cell_87", 87},
		{"base_cell_88", 88},
		{"base_cell_89", 89},
		{"base_cell_90", 90},
		{"base_cell_91", 91},
		{"base_cell_92", 92},
		{"base_cell_93", 93},
		{"base_cell_94", 94},
		{"base_cell_95", 95},
		{"base_cell_96", 96},
		{"base_cell_97", 97}, // pentagon
		{"base_cell_98", 98},
		{"base_cell_99", 99},
		{"base_cell_100", 100},
		{"base_cell_101", 101},
		{"base_cell_102", 102},
		{"base_cell_103", 103},
		{"base_cell_104", 104},
		{"base_cell_105", 105},
		{"base_cell_106", 106},
		{"base_cell_107", 107}, // pentagon
		{"base_cell_108", 108},
		{"base_cell_109", 109},
		{"base_cell_110", 110},
		{"base_cell_111", 111},
		{"base_cell_112", 112},
		{"base_cell_113", 113},
		{"base_cell_114", 114},
		{"base_cell_115", 115},
		{"base_cell_116", 116},
		{"base_cell_117", 117}, // pentagon
		{"base_cell_118", 118},
		{"base_cell_119", 119},
		{"base_cell_120", 120},
		{"base_cell_121", 121},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goResult FaceIJK
			_baseCellToFaceIjk(tt.baseCell, &goResult)

			cResult := _baseCellToFaceIjkC(tt.baseCell)

			if goResult.Face != cResult.Face ||
				goResult.Coord.I != cResult.Coord.I ||
				goResult.Coord.J != cResult.Coord.J ||
				goResult.Coord.K != cResult.Coord.K {
				t.Errorf("_baseCellToFaceIjk(%d): Go=%+v, C=%+v",
					tt.baseCell, goResult, cResult)
			}
		})
	}
}
