//go:build cgo

package c2go

import (
	"testing"
)

func Test_faceIjkToBaseCellCCWrot60_parity(t *testing.T) {
	tests := []struct {
		name string
		h    FaceIJK
	}{
		// Test all valid combinations for each face
		// Face 0
		{"face_0_000", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 0, K: 0}}},
		{"face_0_001", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 0, K: 1}}},
		{"face_0_002", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 0, K: 2}}},
		{"face_0_010", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 1, K: 0}}},
		{"face_0_011", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 1, K: 1}}},
		{"face_0_012", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 1, K: 2}}},
		{"face_0_020", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 2, K: 0}}},
		{"face_0_021", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 2, K: 1}}},
		{"face_0_022", FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 2, K: 2}}},
		{"face_0_100", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 0, K: 0}}},
		{"face_0_101", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 0, K: 1}}},
		{"face_0_102", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 0, K: 2}}},
		{"face_0_110", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 1, K: 0}}},
		{"face_0_111", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 1, K: 1}}},
		{"face_0_112", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 1, K: 2}}},
		{"face_0_120", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 2, K: 0}}},
		{"face_0_121", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 2, K: 1}}},
		{"face_0_122", FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 2, K: 2}}},
		{"face_0_200", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 0, K: 0}}},
		{"face_0_201", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 0, K: 1}}},
		{"face_0_202", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 0, K: 2}}},
		{"face_0_210", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 1, K: 0}}},
		{"face_0_211", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 1, K: 1}}},
		{"face_0_212", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 1, K: 2}}},
		{"face_0_220", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 2, K: 0}}},
		{"face_0_221", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 2, K: 1}}},
		{"face_0_222", FaceIJK{Face: 0, Coord: CoordIJK{I: 2, J: 2, K: 2}}},

		// Test a few from other faces
		{"face_5_111", FaceIJK{Face: 5, Coord: CoordIJK{I: 1, J: 1, K: 1}}},
		{"face_10_111", FaceIJK{Face: 10, Coord: CoordIJK{I: 1, J: 1, K: 1}}},
		{"face_15_111", FaceIJK{Face: 15, Coord: CoordIJK{I: 1, J: 1, K: 1}}},
		{"face_19_111", FaceIJK{Face: 19, Coord: CoordIJK{I: 1, J: 1, K: 1}}},
		{"face_19_222", FaceIJK{Face: 19, Coord: CoordIJK{I: 2, J: 2, K: 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := _faceIjkToBaseCellCCWrot60(&tt.h)
			cResult := _faceIjkToBaseCellCCWrot60C(&tt.h)

			if goResult != cResult {
				t.Errorf("_faceIjkToBaseCellCCWrot60(%+v): Go=%d, C=%d",
					tt.h, goResult, cResult)
			}
		})
	}
}

func Test_faceIjkToBaseCellCCWrot60_all_faces(t *testing.T) {
	// Test all valid combinations systematically
	for face := 0; face < 20; face++ {
		for i := 0; i <= 2; i++ {
			for j := 0; j <= 2; j++ {
				for k := 0; k <= 2; k++ {
					h := FaceIJK{
						Face:  face,
						Coord: CoordIJK{I: int32(i), J: int32(j), K: int32(k)},
					}

					goResult := _faceIjkToBaseCellCCWrot60(&h)
					cResult := _faceIjkToBaseCellCCWrot60C(&h)

					if goResult != cResult {
						t.Errorf("_faceIjkToBaseCellCCWrot60(face=%d, i=%d, j=%d, k=%d): Go=%d, C=%d",
							face, i, j, k, goResult, cResult)
					}
				}
			}
		}
	}
}
