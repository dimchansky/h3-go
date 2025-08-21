package faceijk

import (
	"testing"

	"github.com/dimchansky/h3-go/internal/coordijk"
)

func TestFaceIJKToBaseCell(t *testing.T) {
	tests := []struct {
		name        string
		fijk        FaceIJK
		expectValid bool
	}{
		{"face0_000", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 0}}, true},
		{"face0_100", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 1, J: 0, K: 0}}, true},
		{"face0_200", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 2, J: 0, K: 0}}, true},
		{"face0_010", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 0, J: 1, K: 0}}, true},
		{"face0_020", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 0, J: 2, K: 0}}, true},
		{"face1_000", FaceIJK{Face: 1, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 0}}, true},
		{"face1_100", FaceIJK{Face: 1, Coord: coordijk.CoordIJK{I: 1, J: 0, K: 0}}, true},
		{"face5_110", FaceIJK{Face: 5, Coord: coordijk.CoordIJK{I: 1, J: 1, K: 0}}, true},
		{"face10_001", FaceIJK{Face: 10, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 1}}, true},
		{"face19_220", FaceIJK{Face: 19, Coord: coordijk.CoordIJK{I: 2, J: 2, K: 0}}, true},
		{"invalid_face", FaceIJK{Face: -1, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 0}}, false},
		{"invalid_i_coord", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 3, J: 0, K: 0}}, false},
		{"invalid_j_coord", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 0, J: 3, K: 0}}, false},
		{"invalid_negative", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: -1, J: 0, K: 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := faceIJKToBaseCell(tt.fijk)
			if tt.expectValid && got < 0 {
				t.Errorf("faceIJKToBaseCell(%v) = %d, expected >=0", tt.fijk, got)
			}
			if !tt.expectValid && got >= 0 {
				t.Errorf("faceIJKToBaseCell(%v) = %d, expected <0", tt.fijk, got)
			}
		})
	}
}

func TestFaceIJKToBaseCellCCWrot60(t *testing.T) {
	tests := []struct {
		name        string
		fijk        FaceIJK
		expectValid bool
	}{
		{"face0_000", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 0}}, true},
		{"face0_100", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 1, J: 0, K: 0}}, true},
		{"face5_110", FaceIJK{Face: 5, Coord: coordijk.CoordIJK{I: 1, J: 1, K: 0}}, true},
		{"invalid_face", FaceIJK{Face: -1, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 0}}, false},
		{"invalid_coord", FaceIJK{Face: 0, Coord: coordijk.CoordIJK{I: 3, J: 0, K: 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := faceIJKToBaseCellCCWrot60(tt.fijk)
			if tt.expectValid && got < 0 {
				t.Errorf("faceIJKToBaseCellCCWrot60(%v) = %d, expected >=0", tt.fijk, got)
			}
			if !tt.expectValid && got >= 0 {
				t.Errorf("faceIJKToBaseCellCCWrot60(%v) = %d, expected <0", tt.fijk, got)
			}
		})
	}
}
