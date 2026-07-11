//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_faceIjkToBaseCell_parity(t *testing.T) {
	tests := []struct {
		name string
		h    faceIJK
	}{
		// Test all valid combinations for each face
		// Face 0
		{"face_0_000", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}}},
		{"face_0_001", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 1}}},
		{"face_0_002", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 2}}},
		{"face_0_010", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 1, K: 0}}},
		{"face_0_011", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 1, K: 1}}},
		{"face_0_012", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 1, K: 2}}},
		{"face_0_020", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 2, K: 0}}},
		{"face_0_021", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 2, K: 1}}},
		{"face_0_022", faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 2, K: 2}}},
		{"face_0_100", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 0, K: 0}}},
		{"face_0_101", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 0, K: 1}}},
		{"face_0_102", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 0, K: 2}}},
		{"face_0_110", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 1, K: 0}}},
		{"face_0_111", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 1, K: 1}}},
		{"face_0_112", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 1, K: 2}}},
		{"face_0_120", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 2, K: 0}}},
		{"face_0_121", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 2, K: 1}}},
		{"face_0_122", faceIJK{Face: 0, Coord: coordIJK{I: 1, J: 2, K: 2}}},
		{"face_0_200", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 0, K: 0}}},
		{"face_0_201", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 0, K: 1}}},
		{"face_0_202", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 0, K: 2}}},
		{"face_0_210", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 1, K: 0}}},
		{"face_0_211", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 1, K: 1}}},
		{"face_0_212", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 1, K: 2}}},
		{"face_0_220", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 2, K: 0}}},
		{"face_0_221", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 2, K: 1}}},
		{"face_0_222", faceIJK{Face: 0, Coord: coordIJK{I: 2, J: 2, K: 2}}},

		// Test a few from other faces
		{"face_5_111", faceIJK{Face: 5, Coord: coordIJK{I: 1, J: 1, K: 1}}},
		{"face_10_111", faceIJK{Face: 10, Coord: coordIJK{I: 1, J: 1, K: 1}}},
		{"face_15_111", faceIJK{Face: 15, Coord: coordIJK{I: 1, J: 1, K: 1}}},
		{"face_19_111", faceIJK{Face: 19, Coord: coordIJK{I: 1, J: 1, K: 1}}},
		{"face_19_222", faceIJK{Face: 19, Coord: coordIJK{I: 2, J: 2, K: 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := _faceIjkToBaseCell(&tt.h)
			cResult := _faceIjkToBaseCellC(&tt.h)

			if goResult != cResult {
				t.Errorf("_faceIjkToBaseCell(%+v): Go=%d, C=%d",
					tt.h, goResult, cResult)
			}
		})
	}
}

func Test_faceIjkToBaseCell_all_faces(t *testing.T) {
	// Test all valid combinations systematically
	for face := 0; face < 20; face++ {
		for i := 0; i <= 2; i++ {
			for j := 0; j <= 2; j++ {
				for k := 0; k <= 2; k++ {
					h := faceIJK{
						Face:  int32(face),
						Coord: coordIJK{I: int32(i), J: int32(j), K: int32(k)},
					}

					goResult := _faceIjkToBaseCell(&h)
					cResult := _faceIjkToBaseCellC(&h)

					if goResult != cResult {
						t.Errorf("_faceIjkToBaseCell(face=%d, i=%d, j=%d, k=%d): Go=%d, C=%d",
							face, i, j, k, goResult, cResult)
					}
				}
			}
		}
	}
}
